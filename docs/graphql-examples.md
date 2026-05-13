# StockWise GraphQL Examples

GraphQL endpoint:

```text
http://localhost:8080/graphql
```

Playground:

```text
http://localhost:8080/graphql/playground
```

The GraphQL API uses the same application services as the REST API. Workflow mutations accept user IDs directly in the mutation input or argument instead of reading `X-User-ID`.

## Users

```graphql
query Users {
  users(isActive: true) {
    id
    name
    email
    role
    isActive
  }
}
```

```graphql
mutation CreateUser {
  createUser(input: {
    name: "Nina Operator"
    email: "nina.operator@stockwise.local"
    role: "Operator"
  }) {
    id
    name
    email
    role
  }
}
```

## Warehouses

```graphql
query Warehouses {
  warehouses(isActive: true) {
    id
    name
    code
    location
  }
}
```

```graphql
mutation CreateWarehouse {
  createWarehouse(input: {
    name: "North Warehouse"
    code: "WH-NORTH"
    location: "Ruse, Bulgaria"
  }) {
    id
    name
    code
    isActive
  }
}
```

```graphql
mutation DeactivateWarehouse {
  deactivateWarehouse(id: "3")
}
```

## Categories

```graphql
query Categories {
  categories(isActive: true) {
    id
    name
    description
  }
}
```

```graphql
mutation CreateCategory {
  createCategory(input: {
    name: "Equipment"
    description: "Warehouse equipment and tools"
  }) {
    id
    name
    isActive
  }
}
```

## Products And Stock

```graphql
query Products {
  products(search: "scanner", isActive: true) {
    id
    name
    sku
    categoryID
    unitOfMeasure
    minStockThreshold
    isActive
  }
}
```

```graphql
query ProductsByCategory {
  productsByCategory(categoryID: "1") {
    id
    name
    sku
  }
}
```

```graphql
mutation CreateProduct {
  createProduct(input: {
    name: "Industrial Scale"
    sku: "EQUIP-001"
    categoryID: "1"
    unitOfMeasure: "pcs"
    minStockThreshold: 2
  }) {
    id
    name
    sku
  }
}
```

```graphql
query ProductStock {
  productStock(productID: "1") {
    warehouseID
    productID
    quantity
  }
}
```

```graphql
query ProductTotalStock {
  productTotalStock(productID: "1") {
    productID
    quantity
  }
}
```

```graphql
query LowStockProducts {
  lowStockProducts {
    state
    product {
      id
      name
      sku
      minStockThreshold
    }
    stockItem {
      warehouseID
      quantity
    }
  }
}
```

## Supplies

```graphql
mutation CreateSupply {
  createSupply(input: {
    warehouseID: "1"
    createdByUserID: "3"
  }) {
    id
    status
    warehouseID
  }
}
```

```graphql
mutation AddSupplyItem {
  addSupplyItem(supplyID: "1", input: {
    productID: "1"
    quantity: 10
    unitPriceCents: 12999
  }) {
    id
    productID
    quantity
    unitPriceCents
  }
}
```

```graphql
mutation ConfirmSupply {
  confirmSupply(id: "1", userID: "2") {
    id
    status
    confirmedByUserID
    items {
      productID
      quantity
    }
  }
}
```

## Outbound Requests

```graphql
mutation CreateOutboundRequest {
  createOutboundRequest(input: {
    warehouseID: "1"
    createdByUserID: "3"
  }) {
    id
    status
  }
}
```

```graphql
mutation AddOutboundItem {
  addOutboundRequestItem(outboundRequestID: "1", input: {
    productID: "1"
    quantity: 2
  }) {
    id
    productID
    quantity
  }
}
```

```graphql
mutation ApproveAndExecuteOutbound {
  approveOutboundRequest(id: "1", userID: "2") {
    id
    status
  }
  executeOutboundRequest(id: "1", userID: "3") {
    id
    status
    executedByUserID
  }
}
```

## Transfers

```graphql
mutation CreateTransfer {
  createTransfer(input: {
    sourceWarehouseID: "1"
    targetWarehouseID: "2"
    createdByUserID: "3"
  }) {
    id
    status
    sourceWarehouseID
    targetWarehouseID
  }
}
```

```graphql
mutation AddTransferItem {
  addTransferItem(transferID: "1", input: {
    productID: "1"
    quantity: 1
  }) {
    id
    productID
    quantity
  }
}
```

```graphql
mutation ConfirmTransfer {
  confirmTransfer(id: "1", userID: "2") {
    id
    status
    confirmedByUserID
    items {
      productID
      quantity
    }
  }
}
```

## Movements

```graphql
query ProductMovements {
  productMovements(productID: "1", limit: 20) {
    kind
    operationID
    productID
    warehouseID
    relatedWarehouseID
    quantity
    occurredAt
  }
}
```

```graphql
query WarehouseMovements {
  warehouseMovements(warehouseID: "1", limit: 20) {
    kind
    operationID
    productID
    warehouseID
    relatedWarehouseID
    quantity
    occurredAt
  }
}
```
