package response

import (
	"go_casbin/internal/logger"
	"go_casbin/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`      // 状态码
	Message   string      `json:"message"`   // 消息
	Data      interface{} `json:"data"`      // 数据
	Timestamp int64       `json:"timestamp"` // 时间戳
	Path      string      `json:"path"`      // 请求路径
	Method    string      `json:"method"`    // 请求方法
	TraceID   string      `json:"trace_id"`  // 追踪ID
}

// ResponseMiddleware 统一响应格式中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置响应头
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("X-API-Version", "1.0.0")
		c.Header("X-Server-Time", time.Now().Format(time.RFC3339))
		
		// 生成追踪ID
		traceID := util.GenerateTraceID()
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		
		// 重写响应方法
		c.Next()
		// 如果响应已经被写入，则不处理
		if c.Writer.Written() {
			return
		}

		// 获取响应数据
		data, exists := c.Get("response_data")
		if exists {
			// 构建统一响应格式
			response := Response{
				Code:      http.StatusOK,
				Message:   "success",
				Data:      data,
				Timestamp: time.Now().Unix(),
				Path:      c.Request.URL.Path,
				Method:    c.Request.Method,
				TraceID:   traceID,
			}
			
			c.JSON(http.StatusOK, response)
		}
	}
}



// getTraceID 获取追踪ID
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return util.GenerateTraceID()
}

// ResponseLogger 响应日志中间件
func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		
		// 记录响应日志
		logger.Info("HTTP响应",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.RequestURI),
			logger.Int("status_code", statusCode),
			logger.Duration("latency", latency),
			logger.String("client_ip", c.ClientIP()),
			logger.String("trace_id", GetTraceID(c)),
		)
	}
} 
// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	response := Response{
		Code:      http.StatusOK,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		TraceID:   GetTraceID(c),
	}
	c.JSON(http.StatusOK, response)
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	response := Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		TraceID:   GetTraceID(c),
	}
	
	c.JSON(code, response)
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}


// InternalServerError 500错误响应
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

//程序内部逻辑错误
func LogicError(c *gin.Context, message string) {
	Error(c, -1, message)
}

// ValidationError 验证错误响应
func ValidationError(c *gin.Context, errors interface{}) {
	response := Response{
		Code:      http.StatusBadRequest,
		Message:   "validation_error",
		Data:      errors,
		Timestamp: time.Now().Unix(),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		TraceID:   GetTraceID(c),
	}
	
	c.JSON(http.StatusBadRequest, response)
}

// PaginatedResponse 分页响应
func PaginatedResponse(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	paginatedData := gin.H{
		"list":      data,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	}
	Success(c, paginatedData)
}