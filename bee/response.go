package bee

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	Ext  any    `json:"ext,omitempty"`
}

func pickFirst(items []any) any {
	if len(items) == 0 {
		return nil
	}
	return items[0]
}

// response 返回响应
func response(c *gin.Context, respType responseTypeEnum, code int, msg string, data, ext any) {
	switch respType {
	case StrEnum:
		if s, ok := data.(string); ok {
			c.String(http.StatusOK, s)
			return
		}
		c.String(http.StatusOK, fmt.Sprint(data))
	case JsonEnum:
		res := result{code, msg, data, ext}
		c.JSON(http.StatusOK, res)
	case AsciiJsonEnum:
		res := result{code, msg, data, ext}
		c.AsciiJSON(http.StatusOK, res)
	case XmlEnum:
		res := result{code, msg, data, ext}
		c.XML(http.StatusOK, res)
	case RedirectEnum:
		if s, ok := data.(string); ok {
			c.Redirect(http.StatusFound, s)
			return
		}
		c.Redirect(http.StatusFound, fmt.Sprint(data))
	}
}

// OkStrResponse 成功响应体 - string
func OkStrResponse(c *gin.Context, data string, ext ...any) {
	msg := GetCodeMsg(OK)
	response(c, StrEnum, OK, msg, data, pickFirst(ext))
}

// OkJsonResponse 成功响应体-json
func OkJsonResponse(c *gin.Context, data any, ext ...any) {
	extVal := pickFirst(ext)
	if data == nil {
		data = map[string]any{}
	}
	if extVal == nil {
		extVal = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, JsonEnum, OK, msg, data, extVal)
}

// ErrorJsonResponse 失败响应体-json
func ErrorJsonResponse(c *gin.Context, code int, msg string, ext ...any) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, JsonEnum, code, msg, nil, pickFirst(ext))
}

// OkAsciiJsonResponse 成功响应体-ascii json
func OkAsciiJsonResponse(c *gin.Context, data any, ext ...any) {
	extVal := pickFirst(ext)
	if data == nil {
		data = map[string]any{}
	}
	if extVal == nil {
		extVal = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, AsciiJsonEnum, OK, msg, data, extVal)
}

// ErrorAsciiJsonResponse 失败响应体-ascii json
func ErrorAsciiJsonResponse(c *gin.Context, code int, msg string, ext ...any) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, AsciiJsonEnum, code, msg, nil, pickFirst(ext))
}

// OkXMLResponse 成功响应体 - XML
func OkXMLResponse(c *gin.Context, data any, ext ...any) {
	extVal := pickFirst(ext)
	if data == nil {
		data = map[string]any{}
	}
	if extVal == nil {
		extVal = map[string]any{}
	}
	msg := GetCodeMsg(OK)
	response(c, XmlEnum, OK, msg, data, extVal)
}

// ErrorXMLResponse 失败响应体-XML
func ErrorXMLResponse(c *gin.Context, code int, msg string, ext ...any) {
	if msg == "" {
		msg = GetCodeMsg(code)
	}
	response(c, XmlEnum, code, msg, nil, pickFirst(ext))
}
