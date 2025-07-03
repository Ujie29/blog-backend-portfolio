package middleware

import (
	"blog-backend/common/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 處理所有成功與錯誤的統一回傳格式
func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 執行後續 handler

		// 錯誤處理
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*AppError); ok {
				c.JSON(http.StatusBadRequest, model.APIResponseAny{
					Code:    appErr.Code,
					Message: appErr.Message,
				})
			} else {
				c.JSON(http.StatusInternalServerError, model.APIResponseAny{
					Code:    ErrInternal.Code,
					Message: err.Error(),
				})
			}
			return
		}

		// 成功處理：只要有設 data 就包裝格式回傳
		if data, exists := c.Get("data"); exists {
			c.JSON(http.StatusOK, model.APIResponseAny{
				Code:    ErrOK.Code,
				Message: ErrOK.Message,
				Data:    data,
			})
		}
	}
}

func WrapDBErr(action string, err error) *AppError {
	return Newf(ErrDB.Code, "%s：%v", action, err)
}
