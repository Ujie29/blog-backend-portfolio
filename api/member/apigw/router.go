package apigw

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/idtoken"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	env := os.Getenv("ENV")
	r.Any("/api/:service/*proxyPath", func(c *gin.Context) {
		service := strings.ToUpper(c.Param("service")) // POST
		envKey := fmt.Sprintf("%s_member_SERVICE", service)
		target := os.Getenv(envKey)

		if target == "" {
			log.Printf("❌ 環境變數 %s 未設定，請確認 cloudbuild.yaml 是否有傳入，否則轉發 /api/%s 將會失敗", envKey, strings.ToLower(service))
			target = fmt.Sprintf("MISSING_ENV[%s]", envKey)
		}

		// ✳️ 建構要轉發的完整 URL，例如 https://post/api/post/xxx
		proxyPath := c.Param("proxyPath")
		fullURL := fmt.Sprintf("%s/api/%s%s", target, strings.ToLower(service), proxyPath)

		// ✳️ 讀取原始 body 並重建（否則 POST body 會失效）
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read body"})
			return
		}

		// 👉 建立轉發請求
		req, err := http.NewRequest(c.Request.Method, fullURL, bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header = c.Request.Header.Clone()
		req.URL.RawQuery = c.Request.URL.RawQuery

		// 🔐 如果不是本地，就幫這支 request 加上 ID Token
		if env != "local" {
			tokenSource, err := idtoken.NewTokenSource(context.Background(), target)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token source", "detail": err.Error()})
				return
			}
			token, err := tokenSource.Token()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ID token", "detail": err.Error()})
				return
			}
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		}

		// ✳️ 傳送 request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "proxy request failed", "detail": err.Error()})
			return
		}
		defer resp.Body.Close()

		// ✳️ 傳回原始結果
		c.Status(resp.StatusCode)
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		io.Copy(c.Writer, resp.Body)
	})
}
