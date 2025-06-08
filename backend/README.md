# CertMon Backend

請在此目錄下開發 Golang API 與排程相關程式。

- 主程式建議放於 main.go
- 可依功能拆分 handler、service、model 等子目錄

---

## 開發環境建議

本專案建議使用 docker-compose 啟動 PostgreSQL，確保所有開發者環境一致。

### 快速啟動資料庫

1. 於專案根目錄執行：

```bash
docker-compose up -d postgres
```

2. 複製 .env.example 作為本地設定：

```bash
cp backend/.env.example backend/.env
```

3. 依照 .env 設定連線資訊，預設如下：
- DB_HOST=localhost
- DB_PORT=5432
- DB_USER=certmon
- DB_PASS=certmon
- DB_NAME=certmon

4. 啟動 backend server 前，請確認資料庫已啟動。

---

### 其他說明
- 若需重設資料庫，可刪除 docker volume：
  ```bash
  docker-compose down -v
  ```
- CI/CD 亦會自動啟動 PostgreSQL，確保測試環境一致。
- 若有權限或 port 衝突，請調整 .env 或 docker-compose.yaml。
