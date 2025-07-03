package main

import (
	"blog-backend/api/admin/apigw"
	"fmt"
	"log"
	"os"

	"blog-backend/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ✅ 載入本地的 .env 檔案（只會影響本地）
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 沒有載入 .env 檔案，可能是雲端環境")
	}
	// 建立 Gin Engine
	r := gin.Default()
	// ✅ 套用共用的 CORS middleware
	r.Use(middleware.Handler())
	// ✅ 加入簽章驗證 middleware（防止非 Worker 請求與重放攻擊）
	r.Use(middleware.VerifySignedHeaders())
	// ✅ 保持尾端加斜線統一化
	r.RedirectTrailingSlash = true
	// 註冊路由
	apigw.RegisterRoutes(r)
	// 從環境變數讀取 PORT，預設為 8180
	port := os.Getenv("PORT")
	if port == "" {
		port = "8180"
	}
	// 啟動伺服器
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
