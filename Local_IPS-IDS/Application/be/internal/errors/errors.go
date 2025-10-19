package errors

import (
	"fmt"
	"net/http"
)

// AppError 應用錯誤結構
type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"-"`
}

// Error 實現 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// 錯誤碼常量
const (
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeTimeout        = "TIMEOUT"
)

// New 創建新的應用錯誤
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewWithDetails 創建帶詳情的應用錯誤
func NewWithDetails(code, message string, statusCode int, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: statusCode,
	}
}

// 預定義錯誤
var (
	ErrInternal = New(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrNotFound = New(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrUnauthorized = New(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden = New(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrBadRequest = New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrServiceUnavailable = New(ErrCodeServiceUnavailable, "Service unavailable", http.StatusServiceUnavailable)
	ErrTimeout = New(ErrCodeTimeout, "Request timeout", http.StatusRequestTimeout)
)

// Wrap 包裝錯誤
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsAppError 檢查是否為應用錯誤
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

