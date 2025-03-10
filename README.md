# shortURL

## 概述
`shortURL` 是一個基於 Go 架構的短網址服務，支援建立短網址並透過短網址導向原始連結。本專案使用 [Gin](https://github.com/gin-gonic/gin) 作為 HTTP 伺服器、[sqlc](https://sqlc.dev/) 來自動產生資料庫查詢程式碼，並使用 Redis 進行快取與布隆過濾器處理。

## 專案結構
- **.github/**
  包含 GitHub Actions 的 CI/CD 流程設定。

- **api/**
  API 路由與處理邏輯以及相關處理函式。

- **db/**
  資料庫相關邏輯與 migration 設定：
  - `sqlc/`：自動產生的 SQL 查詢程式碼。
  - `mock/`：使用 GoMock 模擬的資料庫與 Redis 呼叫。

- **util/**
  實用工具函式。

- 其他檔案
  - `main.go`：程式進入點，組合 API、資料庫與 Redis。
  - `makefile`：提供快速執行命令，如測試、資料庫 migration、啟動伺服器等。

## 安裝與部署

### 前置需求
- 安裝 [Go 1.18](https://golang.org/) 或更新版本。
- 確保電腦上有 Docker 以便啟動 PostgreSQL 與建立網路。
- Redis 服務：可使用本機或透過 Docker 運行 Redis。

### 環境變數
專案透過 `app.env` 設定必要的環境變數，包含資料庫連線字串及 HTTP 伺服器地址。請依需求調整此檔案。

### 資料庫設定與 Migration
使用 Makefile 指令來建立與管理資料庫：
- 建立 Docker 網路：
  ```
  make network
  ```
- 啟動 PostgreSQL：
  ```
  make postgres
  ```
- 創建資料庫：
  ```
  make createdb
  ```
- 執行 migration：
  ```
  make migrateup
  ```
- 如需回溯 migration：
  ```
  make migratedown
  ```

### 產生 SQL 查詢程式碼
使用 sqlc 產生 SQL 查詢相關函式：
```
make sqlc
```

### 啟動伺服器
啟動 API 伺服器有兩種方式：
- 使用 Makefile 指令：
  ```
  make server
  ```
- 或直接執行：
  ```
  go run main.go
  ```

## 測試
整個專案包含單元測試（涵蓋 API 路由、資料庫讀寫以及 util 工具），可使用以下指令執行測試：
```
make test
```

## Mock 與 CI/CD
- **Mock 測試**
  使用 GoMock 模擬資料庫與 Redis 呼叫。

- **CI/CD**
  GitHub Actions 設定檔位於 `.github/workflows/`，使用 golangci-lint.yml 以確保程式碼風格與測試品質。

## 結論
此專案展示了如何使用 Go、Gin、sqlc 與 Redis 建立一個高效能的短網址解決方案。歡迎參考程式碼檔案與測試案例，提出建議或貢獻內容。
