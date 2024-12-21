package bee

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// responseType 响应类型枚举
type responseTypeEnum string

const (
	StrEnum       responseTypeEnum = "str"       //  String响应结构
	JsonEnum      responseTypeEnum = "json"      // Json响应结构
	AsciiJsonEnum responseTypeEnum = "asciiJson" // AsciiJson响应结构
	XmlEnum       responseTypeEnum = "xml"       // XML响应结构
	RedirectEnum  responseTypeEnum = "redirect"  // 重定向
)

// result 响应体
type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// response 返回响应
func response(c *gin.Context, respType responseTypeEnum, code int, msg string, data any) {
	switch respType {
	case StrEnum:
		c.String(http.StatusOK, data.(string))
	case JsonEnum:
		res := result{code, msg, data}
		c.JSON(code, res)
	case AsciiJsonEnum:
		res := result{code, msg, data}
		c.AsciiJSON(code, res)
	case XmlEnum:
		res := result{code, msg, data}
		c.XML(code, res)
	case RedirectEnum:
		c.Redirect(http.StatusFound, data.(string))
	}
}

// OkStrResponse 成功响应体 - string
func OkStrResponse(c *gin.Context, data string) {
	msg := GetCodeMsg(OK)
	response(c, StrEnum, OK, msg, data)
}

// OkJsonResponse 成功响应体-json
func OkJsonResponse(c *gin.Context, data any) {
	if data == nil {
		data = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, JsonEnum, OK, msg, data)
}

// ErrorJsonResponse 失败响应体-json
func ErrorJsonResponse(c *gin.Context, code int, msg string) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, JsonEnum, code, msg, nil)
}

// OkAsciiJsonResponse 成功响应体-ascii json
func OkAsciiJsonResponse(c *gin.Context, data any) {
	if data == nil {
		data = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, AsciiJsonEnum, OK, msg, data)
}

// ErrorAsciiJsonResponse 失败响应体-ascii json
func ErrorAsciiJsonResponse(c *gin.Context, code int, msg string) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, AsciiJsonEnum, code, msg, nil)
}

// OkXMLResponse 成功响应体 - XML
func OkXMLResponse(c *gin.Context, data any) {
	if data == nil {
		data = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, XmlEnum, OK, msg, data)
}

// ErrorXMLResponse 失败响应体-XML
func ErrorXMLResponse(c *gin.Context, code int, msg string) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, XmlEnum, code, msg, nil)
}
