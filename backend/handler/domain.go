package handler

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pichu/CertMon/backend/model"
)

// ListDomainsHandler 回傳假資料，供前端測試
func ListDomainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	resp := struct {
		Offset  int            `json:"offset"`
		Limit   int            `json:"limit"`
		Total   int            `json:"total"`
		Domains []model.Domain `json:"domains"`
	}{
		Offset: 0,
		Limit:  20,
		Total:  2,
		Domains: []model.Domain{
			{ID: 1, Domain: "example.com", ValidTo: "2025-07-30T23:59:59Z", DaysLeft: 52, Status: "valid"},
			{ID: 2, Domain: "test.com", ValidTo: "2025-06-30T23:59:59Z", DaysLeft: 21, Status: "expiring"},
		},
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
