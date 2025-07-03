package main

import (
	"blog-backend/common/config"
	"blog-backend/common/middleware"
	"blog-backend/api/member/post"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()

	// repo := post.NewPostRepository(db.DB)
	service := post.NewPostService(db.DB)
	api := post.NewPostAPI(service)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		fmt.Println("Auth header:", c.GetHeader("Authorization"))
		c.Next()
	})
	middleware.RegisterExceptionHandler(r)
	api.RegisterRoutes(r)

	// 讀取 PORT 環境變數，如果沒有就用 8081（給本地用）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// 啟動伺服器
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
