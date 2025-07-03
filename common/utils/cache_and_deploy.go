package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

// PurgeWorkerCacheAndDeployVercel 會先清除快取，再觸發前台部署
func PurgeWorkerCacheAndDeployVercel() error {
	// Step 1: 呼叫 Worker 清除快取
	workerURL := os.Getenv("WORKER_CACHE_PURGE_URL")
	workerToken := os.Getenv("SIGNING_SECRET")

	if workerURL == "" || workerToken == "" {
		return fmt.Errorf("缺少 WORKER_CACHE_PURGE_URL 或 SIGNING_SECRET 環境變數")
	}

	req, err := http.NewRequest("POST", workerURL, nil)
	if err != nil {
		return fmt.Errorf("建立清除快取請求失敗：%v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", workerToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("呼叫清除快取 API 失敗：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("清除快取 API 回傳非預期狀態碼：%d", resp.StatusCode)
	}

	fmt.Println("🧹 成功清除快取")

	// Step 2: 呼叫 Vercel Deploy Hook
	vercelHook := os.Getenv("VERCEL_DEPLOY_HOOK_URL")
	if vercelHook == "" {
		return fmt.Errorf("缺少 VERCEL_DEPLOY_HOOK_URL 環境變數")
	}

	vercelResp, err := http.Post(vercelHook, "application/json", bytes.NewBuffer(nil))
	if err != nil {
		return fmt.Errorf("呼叫 Vercel Deploy Hook 失敗：%v", err)
	}
	defer vercelResp.Body.Close()

	if vercelResp.StatusCode >= 300 {
		return fmt.Errorf("vercel Deploy Hook 回傳非預期狀態碼：%d", vercelResp.StatusCode)
	}

	fmt.Println("🚀 成功觸發 Vercel 部署")

	return nil
}
