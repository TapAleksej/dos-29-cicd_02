package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var requestCount int
var startTime = time.Now()

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/customers", customersHandler)
	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", requestCounter(http.DefaultServeMux))
}

// Middleware для подсчёта запросов
func requestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		next.ServeHTTP(w, r)
	})
}

// ----------- Handlers -----------

// Products CRUD
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			var price float64
			rows.Scan(&id, &name, &price)
			products = append(products, map[string]interface{}{
				"id":    id,
				"name":  name,
				"price": price,
			})
		}
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		var req struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("INSERT INTO products (name, price) VALUES ($1, $2)", req.Name, req.Price)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		var req struct {
			ID    int     `json:"id"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", req.Name, req.Price, req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		var req struct {
			ID int `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("DELETE FROM products WHERE id=$1", req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// Customers CRUD
func customersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, email FROM customers")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var customers []map[string]interface{}
		for rows.Next() {
			var id int
			var name, email string
			rows.Scan(&id, &name, &email)
			customers = append(customers, map[string]interface{}{
				"id":    id,
				"name":  name,
				"email": email,
			})
		}
		json.NewEncoder(w).Encode(customers)

	case http.MethodPost:
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("INSERT INTO customers (name, email) VALUES ($1, $2)", req.Name, req.Email)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		var req struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("UPDATE customers SET name=$1, email=$2 WHERE id=$3", req.Name, req.Email, req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		var req struct {
			ID int `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("DELETE FROM customers WHERE id=$1", req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// Orders CRUD
func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, customer_id, product_id, quantity FROM orders")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var orders []map[string]interface{}
		for rows.Next() {
			var id, customerID, productID, quantity int
			rows.Scan(&id, &customerID, &productID, &quantity)
			orders = append(orders, map[string]interface{}{
				"id":          id,
				"customer_id": customerID,
				"product_id":  productID,
				"quantity":    quantity,
			})
		}
		json.NewEncoder(w).Encode(orders)

	case http.MethodPost:
		var req struct {
			CustomerID int `json:"customer_id"`
			ProductID  int `json:"product_id"`
			Quantity   int `json:"quantity"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("INSERT INTO orders (customer_id, product_id, quantity) VALUES ($1, $2, $3)", req.CustomerID, req.ProductID, req.Quantity)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		var req struct {
			ID         int `json:"id"`
			Quantity   int `json:"quantity"`
			ProductID  int `json:"product_id"`
			CustomerID int `json:"customer_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("UPDATE orders SET customer_id=$1, product_id=$2, quantity=$3 WHERE id=$4", req.CustomerID, req.ProductID, req.Quantity, req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		var req struct {
			ID int `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		_, err := db.Exec("DELETE FROM orders WHERE id=$1", req.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// Метрики
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Количество товаров
	var productCount int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&productCount)
	if err != nil {
		productCount = -1 // для явно ошибочного значения
	}

	uptime := int(time.Since(startTime).Seconds())
	out := map[string]interface{}{
		"product_count": productCount,
		"request_count": requestCount,
		"uptime_sec":    uptime,
	}
	json.NewEncoder(w).Encode(out)
} 
