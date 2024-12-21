package mdw

import (
	"github.com/dhlanshan/go-saillibs/bee"
	"github.com/gin-gonic/gin"
)

// ExceptionMiddleware 异常捕获中间件
func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(bee.Error); ok {
					bee.ErrorJsonResponse(c, e.Code, e.Msg)
				} else {
					bee.ErrorJsonResponse(c, bee.SystemErr, "系统错误。")
				}
				bee.Logger.Error(err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
