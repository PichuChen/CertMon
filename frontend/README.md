# CertMon Frontend

此目錄為 Vue 3 + Tailwind CSS 前端專案。

- 主要元件請放於 src/components/
- 首頁元件：CertMonHome.vue
- 入口檔案：src/App.vue, src/main.js

啟動方式：
```bash
npm install
npm run dev
```

---

## Playwright E2E 測試

啟動本地伺服器後，另開一個終端機執行：
```bash
npx playwright test
```

錄製互動腳本（可選）：
```bash
npx playwright codegen http://localhost:5173
```

如需安裝 Playwright 相關依賴：
```bash
npm install -D playwright @playwright/test
npx playwright install
```

CI 會自動執行 E2E 測試，設定檔於專案根目錄 `.github/workflows/playwright.yml`。
