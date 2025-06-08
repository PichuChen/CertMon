package main

import (
	"log"
	"net/http"

	"github.com/pichu/CertMon/backend/handler"
)

func main() {
	mux := http.NewServeMux()

	// API 路由註冊
	mux.HandleFunc("/api/domains", handler.ListDomainsHandler)

	log.Println("CertMon backend server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
