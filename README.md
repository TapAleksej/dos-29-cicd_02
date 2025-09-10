# Ecommerce API
A simple API for managing online store products using Go, PostgreSQL, and Docker.
## Functions
- **CRUD for products:** create, read, update, delete
- **Metrics:** total number of goods, total stock balance, number of goods by category
- Containerization using Docker

## Endpoints

- `GET|POST|PUT|DELETE /products`
- `GET|POST|PUT|DELETE /customers`
- `GET|POST|PUT|DELETE /orders`
- `GET /metrics`
`
  {
     "product_count": 3,
     "request_count": 42,
     "uptime_sec": 1234
  }
`

Поля для Product: `name`, `price`.\
Поля для Customer: `name`, `email`.\
Поля для Order: `customer_id`, `product_id`, `quantity`.


## Query examples
### Product creation
`curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Laptop","price":123.45}'`

### Customer creation
`curl -X POST http://localhost:8080/customers -H "Content-Type: application/json" -d '{"name":"Vasya","email":"vasya@example.com"}'`
### Order creation
`curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d '{"customer_id":1,"product_id":2,"quantity":5}'`
### Get netrics
`curl -X GET http://localhost:8080/metrics`
### Get orders
`curl -X GET http://localhost:8080/orders`
### Get products
`curl -X GET http://localhost:8080/products`
