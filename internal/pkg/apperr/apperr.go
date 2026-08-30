package apperr

import "fmt"

type Error struct {
	Code    int
	Status  int
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}
func (e *Error) Unwrap() error { return e.Cause }
func New(code, status int, message string) *Error {
	return &Error{Code: code, Status: status, Message: message}
}
func Wrap(code, status int, message string, cause error) *Error {
	return &Error{Code: code, Status: status, Message: message, Cause: cause}
}

var (
	BadRequest    = New(10001, 400, "参数校验失败")
	Unauthorized  = New(10002, 401, "未认证")
	Expired       = New(10003, 401, "访问令牌已过期")
	Forbidden     = New(10004, 403, "无权限")
	NotFound      = New(10005, 404, "资源不存在")
	Conflict      = New(10006, 409, "资源冲突")
	Unprocessable = New(10008, 422, "业务校验失败")
)
