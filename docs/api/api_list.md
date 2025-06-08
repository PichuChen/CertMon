# CertMon API 列表（2025/06/08 初版）

本文件列出第一階段後端 API，供前後端協作與後續細節討論。

---

## 1. Domain（監控對象）管理

- `GET /api/domains`
  - 取得所有監控中的 Domain 列表，支援分頁（預設 20 筆）
  - 輸入參數：`offset`（預設 0）、`limit`（預設 20）
  - 回傳：
    ```json
    {
      "offset": 0,
      "limit": 20,
      "total": 100,
      "domains": [
        {
          "id": 1,
          "domain": "example.com",
          "valid_to": "2025-07-30T23:59:59Z",
          "days_left": 52,
          "status": "valid" // valid | expiring | expired | error | disconnected
        }
        // ...更多 domain
      ]
    }
    ```

- `POST /api/domains`
  - 新增監控 Domain
  - 輸入（Content-Type: application/json）：
    ```json
    {
      "domain": "example.com"
    }
    ```
  - 回傳：
    ```json
    {
      "id": 1,
      "domain": "example.com",
      "created_at": "2025-06-08T12:00:00Z"
    }
    ```

- `GET /api/domains/{id}`
  - 取得單一 Domain 詳細資訊
  - 回傳：
    ```json
    {
      "id": 1,
      "domain": "example.com",
      "valid_to": "2025-07-30T23:59:59Z",
      "days_left": 52,
      "status": "valid",
      "created_at": "2025-06-08T12:00:00Z",
      "updated_at": "2025-06-08T12:00:00Z"
    }
    ```

- `PUT /api/domains/{id}`
  - （目前無實際用途，預留未來擴充，例如允許修改 domain 名稱時使用）
  - 輸入（Content-Type: application/json）：
    ```json
    {
      "domain": "example.com"
    }
    ```
  - 回傳：同 `GET /api/domains/{id}`

- `DELETE /api/domains/{id}`
  - 刪除 Domain
  - 回傳：
    ```json
    {
      "success": true
    }
    ```

## 2. 憑證狀態查詢

- `GET /api/domains/{id}/certificate`
  - 查詢指定 Domain 的 SSL 憑證狀態（到期日、剩餘天數、狀態）
  - 回傳：
    ```json
    {
      "domain_id": 1,
      "domain": "example.com",
      "issuer": "Let's Encrypt",
      "valid_from": "2025-05-01T00:00:00Z",
      "valid_to": "2025-07-30T23:59:59Z",
      "status": "valid", // valid | expiring | expired | error | disconnected
      "days_left": 52,
      "last_checked": "2025-06-08T12:00:00Z",
      "error": null
    }
    ```

- `POST /api/domains/check`
  - 手動觸發所有 Domain 憑證檢查（非同步，僅回傳成功）
  - 輸入：無
  - 回傳：
    ```json
    {
      "success": true
    }
    ```

---

> 登入（login）與通知（notification）API 將於下一階段實作，暫不列入本階段。
> 若前端有 Dashboard 或統計需求，可再補充相關 API。
