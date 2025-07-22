package logger

import (
	"bytes"
	"fmt"
	"go_casbin/internal/logger"
	"go_casbin/internal/middleware/response"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装gin.ResponseWriter以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 基础日志中间件（生产环境使用）
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 获取追踪ID
		traceID := ""
		if param.Keys != nil {
			if id, exists := param.Keys["trace_id"]; exists {
				traceID = id.(string)
			}
		}
		
		// 使用结构化日志记录请求信息
		logger.Info("HTTP请求",
			logger.String("method", param.Method),
			logger.String("path", param.Path),
			logger.String("client_ip", param.ClientIP),
			logger.String("user_agent", param.Request.UserAgent()),
			logger.Int("status_code", param.StatusCode),
			logger.Duration("latency", param.Latency),
			logger.Int("body_size", param.BodySize),
			logger.String("error", param.ErrorMessage),
			logger.String("trace_id", traceID),
		)
		return ""
	})
}

// RequestLogger 详细请求日志中间件（开发环境使用）
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)
		
		// 获取追踪ID
		traceID := response.GetTraceID(c)

		// 记录请求日志
		logger.Info("HTTP请求详情",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("query", c.Request.URL.RawQuery),
			logger.String("client_ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.String("referer", c.Request.Referer()),
			logger.Int("status_code", c.Writer.Status()),
			logger.Duration("latency", latency),
			logger.Int("request_size", len(requestBody)),
			logger.Int("response_size", blw.body.Len()),
			logger.String("request_body", string(requestBody)),
			logger.String("response_body", blw.body.String()),
			logger.String("trace_id", traceID),
		)

		// 记录错误日志
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("请求处理错误",
					logger.String("method", c.Request.Method),
					logger.String("path", c.Request.URL.Path),
					logger.String("error", err.Error()),
					logger.String("error_type", fmt.Sprint(err.Type)),
					logger.String("trace_id", traceID),
				)
			}
		}
	}
}

// ErrorLogger 错误日志中间件
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只记录错误状态码
		if c.Writer.Status() >= 400 {
			logger.Error("HTTP错误响应",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.Int("status_code", c.Writer.Status()),
				logger.String("user_agent", c.Request.UserAgent()),
				logger.String("trace_id", response.GetTraceID(c)),
			)
		}
	}
}
