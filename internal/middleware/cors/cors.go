package cors

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "*")

		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")

		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, X-API-Key, X-Client-Version, X-Platform")

		// 设置允许的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Total-Count, X-API-Version, X-Trace-ID")

		// 设置预检请求缓存时间
		c.Header("Access-Control-Max-Age", "86400")

		// 设置允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// CORSWithConfig 自定义配置的跨域中间件
func CORSWithConfig(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源

		c.Header("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))

		c.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

		c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

		// 设置允许的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Total-Count, X-API-Version, X-Trace-ID, X-Client-Version, X-Platform")

		// 设置预检请求缓存时间
		c.Header("Access-Control-Max-Age", "86400")

		// 设置允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
