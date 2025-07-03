package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Middleware 回傳 gin 的跨域處理 middleware，會從 CORS_ALLOW_ORIGINS 讀取白名單
func Handler() gin.HandlerFunc {
	corsOriginsEnv := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsOriginsEnv == "" {
		log.Fatal("❌ CORS_ALLOW_ORIGINS not set! Refusing to run with insecure defaults.")
	}
	corsOrigins := strings.Split(corsOriginsEnv, ",")

	return cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
