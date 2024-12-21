package mdw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dhlanshan/go-saillibs/bee"
	"github.com/dhlanshan/go-saillibs/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"time"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type LogMWCmd struct {
	NotReqBodyRoute  []string // 不记录请求内容的路由列表
	NotRespBodyRoute []string // 不记录响应内容的路由列表
}

// LogMiddleware 日志
func LogMiddleware(cmd *LogMWCmd) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UnixMilli()
		// 重写request.Body方法使其支持重读
		reqBodyBytes, _ := c.GetRawData()
		if len(reqBodyBytes) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}
		// 重写response使其支持储存
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		// 记录API请求日志 格式："[Api] | 唯一ID | GET | url | header | body | END"
		msgId := fmt.Sprintf("A%s", uuid.New().String())
		// 将消息ID加入到上下文中
		c.Set("trace_id", msgId)

		header, _ := json.Marshal(c.Request.Header)
		msgFormat := "[Api] | %s | %s | %s | Header:%s | Body:%s | END"
		if cmd.NotReqBodyRoute != nil && tools.InSlice[string](cmd.NotReqBodyRoute, c.Request.RequestURI) {
			bee.Logger.Info(fmt.Sprintf(msgFormat, msgId, c.Request.Method, c.Request.RequestURI, header, "当前接口不记录请求内容"))
		} else {
			bee.Logger.Info(fmt.Sprintf(msgFormat, msgId, c.Request.Method, c.Request.RequestURI, header, string(reqBodyBytes)))
		}

		// 执行请求处理程序和其他中间件
		c.Next()
		// 记录API响应日志 格式："[Api] | 唯一ID | 状态 | 内容 | 耗时 | END"
		endTime := time.Now().UnixMilli()
		eTime := fmt.Sprintf("%.5fs", (float64(endTime-startTime))*0.001)
		//记录json响应
		msgFormat = "[Api] | %s | 响应状态: %d | RespBody: %s | 耗时:%s | END"
		if cmd.NotRespBodyRoute != nil && tools.InSlice[string](cmd.NotRespBodyRoute, c.Request.RequestURI) {
			bee.Logger.Info(fmt.Sprintf(msgFormat, msgId, c.Writer.Status(), "当前接口不记录响应内容", eTime))
		} else {
			bee.Logger.Info(fmt.Sprintf(msgFormat, msgId, c.Writer.Status(), blw.body.String(), eTime))
		}
	}
}
