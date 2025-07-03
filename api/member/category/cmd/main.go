package main

import (
	"blog-backend/api/member/category"
	"blog-backend/common/config"
	"blog-backend/common/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()

	service := category.NewCategoryService(db.DB)
	api := category.NewCategoryAPI(service)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		fmt.Println("Auth header:", c.GetHeader("Authorization"))
		c.Next()
	})

	middleware.RegisterExceptionHandler(r)
	api.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
