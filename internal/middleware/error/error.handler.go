package errorhandler

import (
	"fmt"
	"go_casbin/internal/logger"
	"go_casbin/internal/middleware/response"
	"go_casbin/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoRoute 404错误处理
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {		
		// 直接返回404响应，不调用其他响应函数
		response := response.Response{
			Code:      404,
			Message:   "请求的资源不存在",
			Data:      nil,
			Timestamp: time.Now().Unix(),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
			TraceID:   util.GenerateTraceID(),
		}
		
		c.JSON(404, response)
	}
}

// NoMethod 405错误处理
func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {		
		// 直接返回405响应，不调用其他响应函数
		response := response.Response{
			Code:      405,
			Message:   "请求方法不允许",
			Data:      nil,
			Timestamp: time.Now().Unix(),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
			TraceID:   util.GenerateTraceID(),
		}
		
		c.JSON(405, response)
	}
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			logger.Error("请求处理错误",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("error", err.Error()),
				logger.String("error_type", fmt.Sprint(err.Type)),
				logger.String("client_ip", c.ClientIP()),
			)
			
			// 根据错误类型返回不同的响应
			switch err.Type {
			case gin.ErrorTypeBind:
				response.BadRequest(c, "请求参数绑定失败: "+err.Error())
			case gin.ErrorTypePublic:
				response.BadRequest(c, err.Error())
			case gin.ErrorTypePrivate:
				response.InternalServerError(c, "服务器内部错误")
			default:
				response.InternalServerError(c, "未知错误")
			}
		}
	}
}

// CustomErrorHandler 自定义错误处理中间件
func CustomErrorHandler(errorMap map[int]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		statusCode := c.Writer.Status()
		
		// 如果是错误状态码且响应未写入
		if statusCode >= 400 && !c.Writer.Written() {
			message, exists := errorMap[statusCode]
			if !exists {
				message = "请求处理失败"
			}
			
			logger.Error("HTTP错误",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.Int("status_code", statusCode),
				logger.String("message", message),
				logger.String("client_ip", c.ClientIP()),
			)
			
			response.Error(c, statusCode, message)
		}
	}
}

// DefaultErrorMap 默认错误映射
var DefaultErrorMap = map[int]string{
	http.StatusBadRequest:          "请求参数错误",
	http.StatusUnauthorized:        "未授权访问",
	http.StatusForbidden:           "禁止访问",
	http.StatusNotFound:            "资源不存在",
	http.StatusMethodNotAllowed:    "请求方法不允许",
	http.StatusRequestTimeout:      "请求超时",
	http.StatusConflict:            "资源冲突",
	http.StatusUnprocessableEntity: "请求无法处理",
	http.StatusTooManyRequests:     "请求过于频繁",
	http.StatusInternalServerError: "服务器内部错误",
	http.StatusBadGateway:          "网关错误",
	http.StatusServiceUnavailable:  "服务不可用",
}

// DefaultErrorHandler 默认错误处理中间件
func DefaultErrorHandler() gin.HandlerFunc {
	return CustomErrorHandler(DefaultErrorMap)
}

// RateLimitErrorHandler 限流错误处理
func RateLimitErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 检查是否被限流
		if c.Writer.Status() == http.StatusTooManyRequests {
			logger.Warn("请求被限流",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.String("user_agent", c.Request.UserAgent()),
			)
			
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后重试")
		}
	}
}

// ValidationErrorHandler 验证错误处理
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 检查验证错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypeBind {
					logger.Warn("参数验证失败",
						logger.String("method", c.Request.Method),
						logger.String("path", c.Request.URL.Path),
						logger.String("error", err.Error()),
						logger.String("client_ip", c.ClientIP()),
					)
					
					response.ValidationError(c, gin.H{
						"error":   "参数验证失败",
						"details": err.Error(),
					})
					return
				}
			}
		}
	}
} 