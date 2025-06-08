package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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

// GetDomainHandler 取得單一 Domain 詳細資訊（假資料）
func GetDomainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
	// 解析 /api/domains/{idOrDomain}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid domain id or name"}`))
		return
	}
	key := parts[3]
	// 支援數字 id 或 domain 名稱
	if key == "1" || key == "example.com" {
		resp := model.Domain{
			ID:       1,
			Domain:   "example.com",
			ValidTo:  "2025-07-30T23:59:59Z",
			DaysLeft: 52,
			Status:   "valid",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	if key == "2" || key == "test.com" {
		resp := model.Domain{
			ID:       2,
			Domain:   "test.com",
			ValidTo:  "2025-06-30T23:59:59Z",
			DaysLeft: 21,
			Status:   "expiring",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	// 其他情境可依需求擴充
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "domain not found"}`))
}
