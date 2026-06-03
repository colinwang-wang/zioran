package errcode

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

var (
	ErrParam        = &Error{Code: 40001, Message: "参数错误"}
	ErrUnauthorized = &Error{Code: 40101, Message: "未认证"}
	ErrForbidden    = &Error{Code: 40301, Message: "无权限"}
	ErrNotFound     = &Error{Code: 40401, Message: "资源不存在"}
	ErrInternal     = &Error{Code: 50001, Message: "服务器错误"}
)

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}
