package main

import (
	"blog-backend/api/batch"
	"blog-backend/common/config"
	"blog-backend/common/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 DB
	db := config.InitDB()

	// 初始化 Batch Service 與 API
	service := batch.NewBatchService(db.DB)
	api := batch.NewBatchAPI(service)

	// 建立 Gin Engine
	r := gin.Default()

	// 註冊統一錯誤處理 middleware
	middleware.RegisterExceptionHandler(r)

	// Debug Header（可選）
	r.Use(func(c *gin.Context) {
		fmt.Println("Scheduler Trigger:", c.GetHeader("User-Agent"))
		c.Next()
	})

	api.RegisterRoutes(r)

	// 讀取 PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8380"
	}

	// 啟動伺服器
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
