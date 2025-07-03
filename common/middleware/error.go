package middleware

import "fmt"

type AppError struct {
	Code    string
	Message string
}

// 定義錯誤
var (
	// ✅ 一般成功
	ErrOK = New("OK", "成功")

	// ❌ 請求錯誤（輸入錯、格式錯、驗證錯）
	ErrBadRequest = New("ErrBadRequest", "請求格式錯誤或參數無效")
	ErrValidation = New("ErrValidation", "輸入驗證失敗，請確認欄位格式與內容")

	// 🔐 權限相關
	ErrUnauthorized = New("ErrUnauthorized", "未經授權的存取，請先登入或提供有效憑證")
	ErrForbidden    = New("ErrForbidden", "沒有權限存取此資源，請聯絡管理員")

	// 📦 資源查無（文章、使用者、檔案不存在）
	ErrNotFound = New("ErrNotFound", "找不到請求的資源")

	// 🧱 資料層錯誤（DB 失敗、資料有問題）
	ErrDB        = New("ErrDB", "資料庫操作失敗")
	ErrDataError = New("ErrDataError", "資料不正確或不一致")
	ErrContentEmpty  = New("ErrContentEmpty", "內容不能為空")
	ErrTransaction   = New("ErrTransaction", "資料儲存失敗，請稍後再試")

	// 🌐 外部服務錯誤（例如 API call、redis、queue）
	ErrExternalService = New("ErrExternalService", "呼叫外部服務時發生錯誤")

	// 💥 系統錯誤（panic、預期外錯誤）
	ErrInternal   = New("ErrInternal", "伺服器內部錯誤，請聯絡開發團隊處理")
	ErrUnexpected = New("ErrUnexpected", "發生非預期錯誤，請稍後再試")
)

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%s, message=%s", e.Code, e.Message)
}

func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Newf(code, format string, a ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
