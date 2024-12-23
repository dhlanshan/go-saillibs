package bee

import "github.com/dhlanshan/go-saillibs/internal/tools"

const (
	OK        = 200  // 成功
	SystemErr = 1000 // 系统错误
	AuthErr   = 1001 // 认证失败
	ArgErr    = 1002 // 参数错误
	ApiErr    = 1003 // 接口错误
	NormalErr = 1004 // 正常业务错误
)

// GetCodeMsg 获取状态消息
func GetCodeMsg(code int) string {
	codeMap := map[int]string{
		OK:        "success",
		SystemErr: "系统错误",
		AuthErr:   "认证失败",
		ArgErr:    "参数错误",
		ApiErr:    "接口错误",
	}

	return tools.GetMapDefault(codeMap, code, "未知错误类型")
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	if e.Msg == "" {
		e.Msg = GetCodeMsg(e.Code)
	}
	return e.Msg
}
