package handler

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pichu/CertMon/backend/model"
)

var db *sql.DB

// InitDB 初始化資料庫連線，支援 .env 設定。
func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
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
		dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + dbname + " sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

// ListDomainsHandler 回傳資料庫內容
func ListDomainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	rows, err := db.Query(`SELECT id, domain, created_at, updated_at FROM domains ORDER BY id ASC`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "db error"}`))
		return
	}
	defer rows.Close()
	var domains []model.Domain
	for rows.Next() {
		var d model.Domain
		var createdAt, updatedAt string
		if err := rows.Scan(&d.ID, &d.Domain, &createdAt, &updatedAt); err != nil {
			continue
		}
		// 這裡不查憑證狀態，僅回傳 domain 基本資料
		d.Status = "unknown"
		d.ValidTo = ""
		d.DaysLeft = 0
		domains = append(domains, d)
	}
	resp := struct {
		Offset  int            `json:"offset"`
		Limit   int            `json:"limit"`
		Total   int            `json:"total"`
		Domains []model.Domain `json:"domains"`
	}{
		Offset:  0,
		Limit:   20,
		Total:   len(domains),
		Domains: domains,
	}
	json.NewEncoder(w).Encode(resp)
}

// GetDomainHandler 取得單一 Domain 詳細資訊（實際抓取 SSL 憑證）
func GetDomainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid domain id or name"}`))
		return
	}
	domain := parts[3]
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "empty domain"}`))
		return
	}
	// 確保 domain 沒有 http(s):// 前綴
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	// 預設 port 443
	host := domain
	if !strings.Contains(host, ":") {
		host += ":443"
	}
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"domain": domain,
			"status": "disconnected",
			"error":  err.Error(),
		})
		return
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"domain": domain,
			"status": "error",
			"error":  "no certificate found",
		})
		return
	}
	cert := certs[0]
	validTo := cert.NotAfter.UTC()
	daysLeft := int(time.Until(validTo).Hours() / 24)
	status := "valid"
	if daysLeft < 0 {
		status = "expired"
	} else if daysLeft < 30 {
		status = "expiring"
	}
	resp := map[string]interface{}{
		"domain":     domain,
		"issuer":     cert.Issuer.CommonName,
		"valid_from": cert.NotBefore.Format(time.RFC3339),
		"valid_to":   cert.NotAfter.Format(time.RFC3339),
		"days_left":  daysLeft,
		"status":     status,
		"serial":     cert.SerialNumber.String(),
		"type":       "",
		"san":        strings.Join(cert.DNSNames, ","),
	}
	json.NewEncoder(w).Encode(resp)
}

// AddDomainHandler 新增監控 Domain
func AddDomainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	var req struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
		return
	}
	var id int64
	err := db.QueryRow(
		`INSERT INTO domains (domain) VALUES ($1) ON CONFLICT (domain) DO NOTHING RETURNING id`,
		req.Domain,
	).Scan(&id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "domain already exists"}`))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "db error"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"domain": req.Domain,
	})
}
