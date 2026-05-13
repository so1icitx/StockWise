package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/config"
	"github.com/so1icitx/StockWise/internal/domain"
	"github.com/so1icitx/StockWise/internal/infrastructure/postgres"
	"github.com/so1icitx/StockWise/internal/transport/httpapi"
	websocketapi "github.com/so1icitx/StockWise/internal/transport/websocket"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	integrationDB         *sql.DB
	integrationGormStore  *postgres.Store
	postgresContainer     *tcpostgres.PostgresContainer
	integrationSkipReason string
)

type baselineIDs struct {
	adminID            domain.ID
	managerID          domain.ID
	operatorID         domain.ID
	mainWarehouseID    domain.ID
	targetWarehouseID  domain.ID
	emptyWarehouseID   domain.ID
	categoryID         domain.ID
	productID          domain.ID
	secondaryProductID domain.ID
}

type testApp struct {
	server *httptest.Server
	ids    baselineIDs
	cancel context.CancelFunc
}

type warehouseResponse struct {
	ID       domain.ID `json:"id"`
	Code     string    `json:"code"`
	IsActive bool      `json:"is_active"`
}

type supplyResponse struct {
	ID          domain.ID           `json:"id"`
	WarehouseID domain.ID           `json:"warehouse_id"`
	Status      domain.SupplyStatus `json:"status"`
}

type outboundRequestResponse struct {
	ID          domain.ID                    `json:"id"`
	WarehouseID domain.ID                    `json:"warehouse_id"`
	Status      domain.OutboundRequestStatus `json:"status"`
}

type transferResponse struct {
	ID                domain.ID             `json:"id"`
	SourceWarehouseID domain.ID             `json:"source_warehouse_id"`
	TargetWarehouseID domain.ID             `json:"target_warehouse_id"`
	Status            domain.TransferStatus `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type graphqlResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := runPostgresContainer(ctx)
	if err != nil {
		integrationSkipReason = fmt.Sprintf("integration tests require Docker/Testcontainers PostgreSQL: %v", err)
		os.Exit(m.Run())
	}
	postgresContainer = container

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get postgres connection string: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open test database: %v\n", err)
		os.Exit(1)
	}
	integrationDB = db

	if err := db.PingContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ping test database: %v\n", err)
		os.Exit(1)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Fprintf(os.Stderr, "set goose dialect: %v\n", err)
		os.Exit(1)
	}

	migrationsDir, err := migrationsPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve migrations path: %v\n", err)
		os.Exit(1)
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		fmt.Fprintf(os.Stderr, "run goose migrations: %v\n", err)
		os.Exit(1)
	}

	gormDB, err := postgres.Open(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open gorm test database: %v\n", err)
		os.Exit(1)
	}
	integrationGormStore = postgres.NewStore(gormDB)

	code := m.Run()

	sqlDB, err := gormDB.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
	_ = db.Close()
	if postgresContainer != nil {
		_ = postgresContainer.Terminate(context.Background())
	}

	os.Exit(code)
}

func runPostgresContainer(ctx context.Context) (container *tcpostgres.PostgresContainer, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("testcontainers startup panic: %v", recovered)
		}
	}()

	return tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("stockwise_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.BasicWaitStrategies(),
	)
}

func TestCreateWarehouseSucceeds(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	var warehouse warehouseResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/warehouses", map[string]any{
		"name":     "Integration Warehouse",
		"code":     "WH-INT",
		"location": "Ruse",
	}, nil, &warehouse)

	assertStatus(t, resp, http.StatusCreated)
	if warehouse.ID == 0 {
		t.Fatal("expected created warehouse id")
	}
	if warehouse.Code != "WH-INT" {
		t.Fatalf("expected warehouse code WH-INT, got %s", warehouse.Code)
	}
	if !warehouse.IsActive {
		t.Fatal("expected created warehouse to be active")
	}
}

func TestInvalidProductReturnsValidationError(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	var body errorResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/products", map[string]any{
		"sku":                 "BAD-001",
		"category_id":         app.ids.categoryID,
		"unit_of_measure":     "pcs",
		"min_stock_threshold": 0,
	}, nil, &body)

	assertStatus(t, resp, http.StatusBadRequest)
	if body.Error == "" {
		t.Fatal("expected validation error message")
	}
}

func TestDuplicateSKUReturnsConflict(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	var body errorResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/products", map[string]any{
		"name":                "Duplicate Scanner",
		"sku":                 "ELEC-001",
		"category_id":         app.ids.categoryID,
		"unit_of_measure":     "pcs",
		"min_stock_threshold": 5,
	}, nil, &body)

	assertStatus(t, resp, http.StatusConflict)
	if !strings.Contains(body.Error, "sku") {
		t.Fatalf("expected SKU conflict message, got %q", body.Error)
	}
}

func TestConfirmSupplyUpdatesDBStock(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	supply := app.createSupply(t, app.ids.mainWarehouseID)
	app.addSupplyItem(t, supply.ID, app.ids.productID, 7)

	var confirmed supplyResponse
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/supplies/%d/confirm", supply.ID), nil, headerUser(app.ids.managerID), &confirmed)

	assertStatus(t, resp, http.StatusOK)
	if confirmed.Status != domain.SupplyStatusCompleted {
		t.Fatalf("expected completed supply, got %s", confirmed.Status)
	}
	assertDBStock(t, app.ids.mainWarehouseID, app.ids.productID, 17)
}

func TestExecuteOutboundDecreasesStock(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	request := app.createApprovedOutbound(t, 4)

	var completed outboundRequestResponse
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/outbound-requests/%d/execute", request.ID), nil, headerUser(app.ids.operatorID), &completed)

	assertStatus(t, resp, http.StatusOK)
	if completed.Status != domain.OutboundRequestStatusCompleted {
		t.Fatalf("expected completed outbound request, got %s", completed.Status)
	}
	assertDBStock(t, app.ids.mainWarehouseID, app.ids.productID, 6)
}

func TestInsufficientOutboundExecutionFails(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	request := app.createApprovedOutbound(t, 99)

	var body errorResponse
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/outbound-requests/%d/execute", request.ID), nil, headerUser(app.ids.operatorID), &body)

	assertStatus(t, resp, http.StatusConflict)
	if !strings.Contains(body.Error, "stock") {
		t.Fatalf("expected stock error, got %q", body.Error)
	}
	assertDBStock(t, app.ids.mainWarehouseID, app.ids.productID, 10)
}

func TestConfirmTransferMovesStock(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	transfer := app.createTransfer(t)
	app.addTransferItem(t, transfer.ID, app.ids.productID, 3)

	var completed transferResponse
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/transfers/%d/confirm", transfer.ID), nil, headerUser(app.ids.managerID), &completed)

	assertStatus(t, resp, http.StatusOK)
	if completed.Status != domain.TransferStatusCompleted {
		t.Fatalf("expected completed transfer, got %s", completed.Status)
	}
	assertDBStock(t, app.ids.mainWarehouseID, app.ids.productID, 7)
	assertDBStock(t, app.ids.targetWarehouseID, app.ids.productID, 3)
}

func TestDeleteWarehouseWithStockFails(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	var body errorResponse
	resp := app.doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v1/warehouses/%d", app.ids.mainWarehouseID), nil, nil, &body)

	assertStatus(t, resp, http.StatusConflict)
	if !strings.Contains(body.Error, "stock") {
		t.Fatalf("expected stock delete restriction, got %q", body.Error)
	}
}

func TestGraphQLQueryWorks(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	query := `query { products(sku: "ELEC-001") { id sku name } }`
	response := app.graphql(t, query)
	if len(response.Errors) > 0 {
		t.Fatalf("expected no GraphQL errors, got %+v", response.Errors)
	}

	var products []struct {
		ID   string `json:"id"`
		SKU  string `json:"sku"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(response.Data["products"], &products); err != nil {
		t.Fatalf("decode products: %v", err)
	}
	if len(products) != 1 || products[0].SKU != "ELEC-001" {
		t.Fatalf("expected ELEC-001 product, got %+v", products)
	}
}

func TestGraphQLMutationWorks(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	mutation := `mutation {
		createWarehouse(input: { name: "GraphQL Warehouse", code: "WH-GQL", location: "Burgas", isActive: true }) {
			id
			code
			isActive
		}
	}`
	response := app.graphql(t, mutation)
	if len(response.Errors) > 0 {
		t.Fatalf("expected no GraphQL errors, got %+v", response.Errors)
	}

	var warehouse struct {
		ID       string `json:"id"`
		Code     string `json:"code"`
		IsActive bool   `json:"isActive"`
	}
	if err := json.Unmarshal(response.Data["createWarehouse"], &warehouse); err != nil {
		t.Fatalf("decode createWarehouse: %v", err)
	}
	if warehouse.Code != "WH-GQL" || !warehouse.IsActive {
		t.Fatalf("unexpected GraphQL warehouse response: %+v", warehouse)
	}
}

func TestWebSocketReceivesEventFromRealBusinessAction(t *testing.T) {
	app := newTestApp(t)
	defer app.close()

	wsURL := "ws" + strings.TrimPrefix(app.server.URL, "http") + "/ws/notifications"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("connect websocket: %v", err)
	}
	defer conn.Close()

	supply := app.createSupply(t, app.ids.mainWarehouseID)
	app.addSupplyItem(t, supply.ID, app.ids.productID, 2)

	var confirmed supplyResponse
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/supplies/%d/confirm", supply.ID), nil, headerUser(app.ids.managerID), &confirmed)
	assertStatus(t, resp, http.StatusOK)

	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var event application.NotificationEvent
	if err := conn.ReadJSON(&event); err != nil {
		t.Fatalf("read websocket event: %v", err)
	}
	if event.Event != application.NotificationSupplyConfirmed {
		t.Fatalf("expected supply.confirmed event, got %s", event.Event)
	}
	if got := event.Data["supply_id"]; got != float64(supply.ID) {
		t.Fatalf("expected supply_id %d in websocket payload, got %v", supply.ID, got)
	}
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()
	if integrationSkipReason != "" {
		t.Skip(integrationSkipReason)
	}
	if integrationGormStore == nil || integrationDB == nil {
		t.Fatal("integration database is not initialized")
	}

	ids := resetDatabase(t)
	hub := websocketapi.NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	go hub.Run(ctx)

	services := application.NewServices(integrationGormStore, integrationGormStore, hub)
	router := httpapi.NewRouter(config.Config{
		AppName: "StockWise",
		AppEnv:  "test",
	}, services, hub)
	server := httptest.NewServer(router)

	return &testApp{
		server: server,
		ids:    ids,
		cancel: cancel,
	}
}

func (app *testApp) close() {
	app.server.Close()
	app.cancel()
}

func (app *testApp) createSupply(t *testing.T, warehouseID domain.ID) supplyResponse {
	t.Helper()
	var supply supplyResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/supplies", map[string]any{
		"warehouse_id": warehouseID,
	}, headerUser(app.ids.operatorID), &supply)
	assertStatus(t, resp, http.StatusCreated)
	return supply
}

func (app *testApp) addSupplyItem(t *testing.T, supplyID domain.ID, productID domain.ID, quantity int64) {
	t.Helper()
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/supplies/%d/items", supplyID), map[string]any{
		"product_id":       productID,
		"quantity":         quantity,
		"unit_price_cents": 1999,
	}, nil, nil)
	assertStatus(t, resp, http.StatusCreated)
}

func (app *testApp) createApprovedOutbound(t *testing.T, quantity int64) outboundRequestResponse {
	t.Helper()
	var request outboundRequestResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/outbound-requests", map[string]any{
		"warehouse_id": app.ids.mainWarehouseID,
	}, headerUser(app.ids.operatorID), &request)
	assertStatus(t, resp, http.StatusCreated)

	resp = app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/outbound-requests/%d/items", request.ID), map[string]any{
		"product_id": app.ids.productID,
		"quantity":   quantity,
	}, nil, nil)
	assertStatus(t, resp, http.StatusCreated)

	var approved outboundRequestResponse
	resp = app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/outbound-requests/%d/approve", request.ID), nil, headerUser(app.ids.managerID), &approved)
	assertStatus(t, resp, http.StatusOK)
	return approved
}

func (app *testApp) createTransfer(t *testing.T) transferResponse {
	t.Helper()
	var transfer transferResponse
	resp := app.doJSON(t, http.MethodPost, "/api/v1/transfers", map[string]any{
		"source_warehouse_id": app.ids.mainWarehouseID,
		"target_warehouse_id": app.ids.targetWarehouseID,
	}, headerUser(app.ids.operatorID), &transfer)
	assertStatus(t, resp, http.StatusCreated)
	return transfer
}

func (app *testApp) addTransferItem(t *testing.T, transferID domain.ID, productID domain.ID, quantity int64) {
	t.Helper()
	resp := app.doJSON(t, http.MethodPost, fmt.Sprintf("/api/v1/transfers/%d/items", transferID), map[string]any{
		"product_id": productID,
		"quantity":   quantity,
	}, nil, nil)
	assertStatus(t, resp, http.StatusCreated)
}

func (app *testApp) graphql(t *testing.T, query string) graphqlResponse {
	t.Helper()
	var response graphqlResponse
	resp := app.doJSON(t, http.MethodPost, "/graphql", map[string]any{
		"query": query,
	}, nil, &response)
	assertStatus(t, resp, http.StatusOK)
	return response
}

func (app *testApp) doJSON(t *testing.T, method string, path string, payload any, headers map[string]string, target any) *http.Response {
	t.Helper()
	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		body = bytes.NewReader(raw)
	}

	request, err := http.NewRequest(method, app.server.URL+path, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := app.server.Client().Do(request)
	if err != nil {
		t.Fatalf("perform request %s %s: %v", method, path, err)
	}
	defer response.Body.Close()

	if target != nil {
		if err := json.NewDecoder(response.Body).Decode(target); err != nil {
			t.Fatalf("decode response body for %s %s: %v", method, path, err)
		}
	}

	return response
}

func resetDatabase(t *testing.T) baselineIDs {
	t.Helper()

	statements := []string{
		`TRUNCATE transfer_items, transfers, outbound_request_items, outbound_requests, supply_items, supplies, stock_items, products, categories, warehouses, users RESTART IDENTITY CASCADE`,
		`INSERT INTO users (name, email, role, is_active) VALUES
			('Alice Admin', 'alice.admin@integration.test', 'Admin', true),
			('Martin Manager', 'martin.manager@integration.test', 'Manager', true),
			('Olivia Operator', 'olivia.operator@integration.test', 'Operator', true)`,
		`INSERT INTO warehouses (name, code, location, is_active) VALUES
			('Central Warehouse', 'WH-CENTRAL', 'Sofia', true),
			('Retail Warehouse', 'WH-RETAIL', 'Plovdiv', true),
			('Empty Warehouse', 'WH-EMPTY', 'Varna', true)`,
		`INSERT INTO categories (name, description, is_active) VALUES
			('Electronics', 'Integration electronics', true)`,
		`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
			SELECT 'Wireless Barcode Scanner', 'ELEC-001', id, 'pcs', 5, true FROM categories WHERE name = 'Electronics'`,
		`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
			SELECT 'Thermal Label Printer', 'ELEC-002', id, 'pcs', 3, true FROM categories WHERE name = 'Electronics'`,
		`INSERT INTO stock_items (warehouse_id, product_id, quantity)
			SELECT w.id, p.id, 10 FROM warehouses w, products p WHERE w.code = 'WH-CENTRAL' AND p.sku = 'ELEC-001'`,
	}

	for _, statement := range statements {
		if _, err := integrationDB.Exec(statement); err != nil {
			t.Fatalf("reset database statement failed: %v\n%s", err, statement)
		}
	}

	return baselineIDs{
		adminID:            lookupID(t, "users", "email", "alice.admin@integration.test"),
		managerID:          lookupID(t, "users", "email", "martin.manager@integration.test"),
		operatorID:         lookupID(t, "users", "email", "olivia.operator@integration.test"),
		mainWarehouseID:    lookupID(t, "warehouses", "code", "WH-CENTRAL"),
		targetWarehouseID:  lookupID(t, "warehouses", "code", "WH-RETAIL"),
		emptyWarehouseID:   lookupID(t, "warehouses", "code", "WH-EMPTY"),
		categoryID:         lookupID(t, "categories", "name", "Electronics"),
		productID:          lookupID(t, "products", "sku", "ELEC-001"),
		secondaryProductID: lookupID(t, "products", "sku", "ELEC-002"),
	}
}

func lookupID(t *testing.T, table string, column string, value string) domain.ID {
	t.Helper()
	var id uint64
	query := fmt.Sprintf("SELECT id FROM %s WHERE %s = $1", table, column)
	if err := integrationDB.QueryRow(query, value).Scan(&id); err != nil {
		t.Fatalf("lookup id from %s.%s=%s: %v", table, column, value, err)
	}
	return domain.ID(id)
}

func assertDBStock(t *testing.T, warehouseID domain.ID, productID domain.ID, want int64) {
	t.Helper()
	var quantity int64
	err := integrationDB.QueryRow(`SELECT quantity FROM stock_items WHERE warehouse_id = $1 AND product_id = $2`, warehouseID, productID).Scan(&quantity)
	if err != nil {
		t.Fatalf("query stock row warehouse=%d product=%d: %v", warehouseID, productID, err)
	}
	if quantity != want {
		t.Fatalf("expected DB stock quantity %d, got %d", want, quantity)
	}
}

func headerUser(id domain.ID) map[string]string {
	return map[string]string{"X-User-ID": fmt.Sprint(id)}
}

func assertStatus(t *testing.T, response *http.Response, want int) {
	t.Helper()
	if response.StatusCode != want {
		t.Fatalf("expected HTTP %d, got %d", want, response.StatusCode)
	}
}

func migrationsPath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(workingDir, "..", "..", "migrations"))
}
