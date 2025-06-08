package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/pichu/CertMon/backend/handler"
)

func main() {
	// 讀取 DB 連線資訊，優先用環境變數
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "certmon"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "certmon"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "certmon"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	// goose migration with retry/failback waiting
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		err = goose.Up(db, "./model")
		if err == nil {
			break
		}
		log.Printf("goose migration failed (try %d/%d): %v", i, maxRetry, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("goose migration failed after %d retries: %v", maxRetry, err)
	}

	// 初始化 DB
	handler.InitDB()

	// 啟動 background worker，每 30 分鐘自動檢查所有 domain
	go func() {
		for {
			log.Println("[CertMon] 執行自動 SSL 憑證狀態檢查...")
			domains, err := getAllDomains()
			if err == nil {
				for _, d := range domains {
					go handler.CheckAndLogDomain(d)
				}
			}
			time.Sleep(30 * time.Minute)
		}
	}()

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

// 取得所有 domain 名稱
func getAllDomains() ([]string, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=certmon password=certmon dbname=certmon sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT domain FROM domains")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var domains []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err == nil {
			domains = append(domains, d)
		}
	}
	return domains, nil
}
