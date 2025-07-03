// X-Timestamp 驗證：請求必須在 5 分鐘內有效
// X-Signature 驗證：確保簽章是根據 secret + timestamp 計算出的
// 支援雙組 Secret 輪替：可同時驗證主、備用金鑰

package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	headerTimestamp = "X-Timestamp"
	headerSignature = "X-Signature"
	timeWindowSec   = 300 // 允許誤差時間（秒）5 分鐘
)

func VerifySignedHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestampStr := c.GetHeader(headerTimestamp)
		signature := c.GetHeader(headerSignature)

		if timestampStr == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少簽章 header"})
			return
		}

		ts, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "時間戳格式錯誤"})
			return
		}

		now := time.Now().Unix()
		if abs(now-ts) > timeWindowSec {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "簽章過期"})
			return
		}

		// 支援雙組 token 驗證（輪替用）
		// 將來要更新worker與apigw之間的token時，因為不能直接把原本的token拿掉會發生問題，所以這邊使用輪替的方式
		// 目前先不給SIGNING_SECRET_SECONDARY的環境變數，現在不給也不會發生問題
		secrets := []string{
			os.Getenv("SIGNING_SECRET"),
			os.Getenv("SIGNING_SECRET_SECONDARY"),
		}

		valid := false
		for _, secret := range secrets {
			if secret == "" {
				continue
			}
			if generateHMAC(timestampStr, secret) == signature {
				valid = true
				break
			}
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "簽章不正確"})
			return
		}

		c.Next()
	}
}

func generateHMAC(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
