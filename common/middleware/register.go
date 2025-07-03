package middleware

import "github.com/gin-gonic/gin"

// 將錯誤處理中間件統一註冊進 router
func RegisterExceptionHandler(r *gin.Engine) {
	r.Use(ExceptionMiddleware())
}
