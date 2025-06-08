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

// ListDomainsHandler 回傳資料庫內容，並從 domain_logs 取得最新狀態
func ListDomainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	rows, err := db.Query(`SELECT d.id, d.domain, d.created_at, d.updated_at, l.status, l.valid_to, l.days_left FROM domains d LEFT JOIN LATERAL (
		SELECT status, valid_to, days_left FROM domain_logs WHERE domain_id = d.id ORDER BY checked_at DESC LIMIT 1
	) l ON true ORDER BY d.id ASC`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "db error"}`))
		return
	}
	defer rows.Close()
	var domains []model.Domain
	for rows.Next() {
		var d model.Domain
		var createdAt, updatedAt, validTo sql.NullString
		var status sql.NullString
		var daysLeft sql.NullInt64
		if err := rows.Scan(&d.ID, &d.Domain, &createdAt, &updatedAt, &status, &validTo, &daysLeft); err != nil {
			continue
		}
		if status.Valid {
			d.Status = status.String
		} else {
			d.Status = "unknown"
		}
		if validTo.Valid {
			d.ValidTo = validTo.String
		} else {
			d.ValidTo = ""
		}
		if daysLeft.Valid {
			d.DaysLeft = int(daysLeft.Int64)
		} else {
			d.DaysLeft = 0
		}
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

// GetDomainHandler 取得單一 Domain 詳細資訊（實際抓取 SSL 憑證，並寫入狀態）
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
	var status, validFrom, validTo, issuer, serial, san string
	var daysLeft int
	if err != nil {
		status = "disconnected"
	} else {
		defer conn.Close()
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			status = "error"
		} else {
			cert := certs[0]
			validToTime := cert.NotAfter.UTC()
			validFrom = cert.NotBefore.Format(time.RFC3339)
			validTo = cert.NotAfter.Format(time.RFC3339)
			daysLeft = int(time.Until(validToTime).Hours() / 24)
			issuer = cert.Issuer.CommonName
			serial = cert.SerialNumber.String()
			san = strings.Join(cert.DNSNames, ",")
			status = "valid"
			if daysLeft < 0 {
				status = "expired"
			} else if daysLeft < 30 {
				status = "expiring"
			}
		}
	}
	// 將狀態寫入資料庫
	_, _ = db.Exec(`UPDATE domains SET updated_at=NOW() WHERE domain=$1`, domain)
	// 新增 log 紀錄，將非索引資料寫入 extra 欄位
	var domainID int64
	db.QueryRow(`SELECT id FROM domains WHERE domain=$1`, domain).Scan(&domainID)
	extra := map[string]interface{}{
		"issuer": issuer,
		"serial": serial,
		"san":    san,
		"error": func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}
	extraJSON, _ := json.Marshal(extra)
	_, _ = db.Exec(`INSERT INTO domain_logs (domain_id, checked_at, status, valid_from, valid_to, days_left, extra) VALUES ($1, NOW(), $2, $3, $4, $5, $6)`,
		domainID, status, validFrom, validTo, daysLeft, extraJSON,
	)
	// 若有需要可擴充寫入 status, valid_to, days_left 等欄位
	resp := map[string]interface{}{
		"domain":     domain,
		"issuer":     issuer,
		"valid_from": validFrom,
		"valid_to":   validTo,
		"days_left":  daysLeft,
		"status":     status,
		"serial":     serial,
		"type":       "",
		"san":        san,
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
