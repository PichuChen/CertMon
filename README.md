# CertMon Monorepo

本專案為 SSL 憑證到期自動監控系統，採用 Golang + Vue + Postgres，前後端同倉管理。

## 專案結構

```
CertMon/
├── backend/        # Golang server (API, 憑證檢查、排程)
├── frontend/       # Vue 3 + Tailwind 前端專案
│   └── src/
│       └── components/
│           └── CertMonHome.vue
├── docs/           # 專案文件、會議紀錄
│   └── meeting/
│       └── 20240608_01.md
├── LICENSE         # MIT License
├── README.md       # 專案說明
└── .gitignore      # Node, Go, 其他忽略規則
```

## 快速開始

### 前端 (Vue 3 + Tailwind)
```bash
cd frontend
npm install
npm run dev
```

### 後端 (Golang)
```bash
cd backend
# 初始化與啟動指令，依後端設計補充
```

---

- 詳細 API、DB schema、設計稿請見 docs/meeting/20240608_01.md
- Figma 設計稿：[CertMon Figma](https://www.figma.com/design/7D4JnVFhEInYE8WLyqnd3g/CertMon?node-id=3604-1626&t=iTMPph78L3sM1lER-1)
