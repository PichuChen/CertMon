package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/pichu/CertMon/backend/handler"
)

func main() {
	// 讀取 DB 連線資訊（可用 os.Getenv 或直接寫死測試）
	dsn := "host=localhost port=5432 user=certmon password=certmon dbname=certmon sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	err = goose.Up(db, "./model")
	if err != nil {
		log.Fatalf("goose migration failed: %v", err)
	}

	// 初始化 DB
	handler.InitDB()

	mux := http.NewServeMux()
	// API 路由註冊
	mux.HandleFunc("/api/domains", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handler.AddDomainHandler(w, r)
		} else {
			handler.ListDomainsHandler(w, r)
		}
	})
	mux.HandleFunc("/api/domains/", handler.GetDomainHandler) // 支援 /api/domains/{id}

	log.Println("CertMon backend server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
