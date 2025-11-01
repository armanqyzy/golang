package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

var db *pgxpool.Pool

func main() {
	var err error

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/products_db?sslmode=disable"
	}

	db, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		log.Fatal("Cannot ping database:", err)
	}
	log.Println("Successfully connected to database")

	http.HandleFunc("/products", getProductsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	params := r.URL.Query()

	var whereConditions []string
	var args []interface{}
	argIndex := 1

	if category := params.Get("category"); category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.name = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	if minPrice := params.Get("min_price"); minPrice != "" {
		if price, err := strconv.Atoi(minPrice); err == nil {
			whereConditions = append(whereConditions, fmt.Sprintf("p.price >= $%d", argIndex))
			args = append(args, price)
			argIndex++
		}
	}

	if maxPrice := params.Get("max_price"); maxPrice != "" {
		if price, err := strconv.Atoi(maxPrice); err == nil {
			whereConditions = append(whereConditions, fmt.Sprintf("p.price <= $%d", argIndex))
			args = append(args, price)
			argIndex++
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	orderClause := ""
	if sort := params.Get("sort"); sort != "" {
		switch sort {
		case "price_asc":
			orderClause = "ORDER BY p.price ASC"
		case "price_desc":
			orderClause = "ORDER BY p.price DESC"
		default:
			orderClause = ""
		}
	}

	limit := 10
	if limitParam := params.Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetParam := params.Get("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	query := fmt.Sprintf(`
        SELECT p.id, p.name, c.name as category, p.price
        FROM products p
        JOIN categories c ON p.category_id = c.id
        %s
        %s
        LIMIT $%d OFFSET $%d
    `, whereClause, orderClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(w, "Failed to parse results", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		http.Error(w, "Error reading results", http.StatusInternalServerError)
		return
	}

	queryTime := time.Since(start).Milliseconds()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Query-Time", fmt.Sprintf("%dms", queryTime))

	if products == nil {
		products = []Product{}
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Query completed in %dms, returned %d products", queryTime, len(products))
}
