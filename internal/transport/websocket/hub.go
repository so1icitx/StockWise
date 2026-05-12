package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/so1icitx/StockWise/internal/application"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	sendBufferSize = 32
)

// Hub manages WebSocket notification clients and publishes application events to them.
type Hub struct {
	upgrader   gorillawebsocket.Upgrader
	clients    map[*client]struct{}
	register   chan *client
	unregister chan *client
	broadcast  chan application.NotificationEvent
}

// NewHub creates a WebSocket notification hub.
func NewHub() *Hub {
	return &Hub{
		upgrader: gorillawebsocket.Upgrader{
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
		clients:    make(map[*client]struct{}),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan application.NotificationEvent, 128),
	}
}

// Run processes client registration and broadcast events until the context is cancelled.
func (hub *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for client := range hub.clients {
				hub.closeClient(client)
			}
			return
		case client := <-hub.register:
			hub.clients[client] = struct{}{}
		case client := <-hub.unregister:
			hub.closeClient(client)
		case event := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- event:
				default:
					hub.closeClient(client)
				}
			}
		}
	}
}

// Publish broadcasts a notification event to connected WebSocket clients.
func (hub *Hub) Publish(ctx context.Context, event application.NotificationEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}

	select {
	case hub.broadcast <- event:
	case <-ctx.Done():
	default:
		log.Printf("dropped websocket notification %s because the broadcast buffer is full", event.Event)
	}
}

// HandleNotifications upgrades an HTTP request to the StockWise notification WebSocket.
func (hub *Hub) HandleNotifications(ctx *gin.Context) {
	conn, err := hub.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("upgrade websocket notification connection: %v", err)
		return
	}

	client := &client{
		hub:  hub,
		conn: conn,
		send: make(chan application.NotificationEvent, sendBufferSize),
	}

	hub.register <- client
	go client.writePump()
	go client.readPump()
}

func (hub *Hub) closeClient(client *client) {
	if _, ok := hub.clients[client]; !ok {
		return
	}
	delete(hub.clients, client)
	close(client.send)
}
