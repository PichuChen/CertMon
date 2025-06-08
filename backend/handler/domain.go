package handler

import (
	"encoding/json"
	"net/http"

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
