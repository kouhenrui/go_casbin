package recovery

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go_casbin/internal/logger"

	"github.com/gin-gonic/gin"
)

// Recovery 全局异常捕获中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic恢复",
				logger.String("error", err),
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.String("user_agent", c.Request.UserAgent()),
				logger.String("stack", string(debug.Stack())),
			)
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "服务器内部错误",
				"error":   err,
			})
		} else {
			logger.Error("Panic恢复",
				logger.String("error", fmt.Sprintf("%v", recovered)),
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.String("user_agent", c.Request.UserAgent()),
				logger.String("stack", string(debug.Stack())),
			)
			
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "服务器内部错误",
				"error":   fmt.Sprintf("%v", recovered),
			})
		}
	})
}

// RecoveryWithLogger 带详细日志的异常捕获中间件
func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录详细的错误信息
				logger.Error("请求处理异常",
					logger.String("error", fmt.Sprintf("%v", err)),
					logger.String("method", c.Request.Method),
					logger.String("path", c.Request.URL.Path),
					logger.String("query", c.Request.URL.RawQuery),
					logger.String("client_ip", c.ClientIP()),
					logger.String("user_agent", c.Request.UserAgent()),
					logger.String("referer", c.Request.Referer()),
					logger.String("stack", string(debug.Stack())),
				)
				
				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   fmt.Sprintf("%v", err),
					"path":    c.Request.URL.Path,
					"method":  c.Request.Method,
				})
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
}

// RecoveryWithCustomHandler 自定义异常处理中间件
func RecoveryWithCustomHandler(handler func(*gin.Context, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				logger.Error("请求处理异常",
					logger.String("error", fmt.Sprintf("%v", err)),
					logger.String("method", c.Request.Method),
					logger.String("path", c.Request.URL.Path),
					logger.String("stack", string(debug.Stack())),
				)
				
				// 调用自定义处理函数
				handler(c, err)
				
				c.Abort()
			}
		}()
		
		c.Next()
	}
} 