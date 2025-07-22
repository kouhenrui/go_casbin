package jwt

import (
	"go_casbin/internal/config"
	"go_casbin/internal/logger"
	"go_casbin/internal/middleware/response"
	"go_casbin/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在排除路径中
		if isExcludedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("JWT认证失败 - 缺少Authorization头",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			response.Unauthorized(c, "缺少认证Token")
			c.Abort()
			return
		}

		// 检查Token前缀
		if !strings.HasPrefix(authHeader, config.ViperConfig.JWT.TokenPrefix) && !strings.HasPrefix(authHeader, config.ViperConfig.JWT.RefreshPrefix) {
			logger.Warn("JWT认证失败 - Token格式错误",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		jwtService := jwt.GetJWTInstance()
		claims, err := jwtService.ParseToken(authHeader)
		if err != nil {
			logger.Warn("JWT认证失败 - Token解析错误",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.String("error", err.Error()),
			)
			response.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}
		// 将用户信息存储到上下文中
		c.Set("account", claims)
		c.Next()
	}
}

// isExcludedPath 检查路径是否在白名单中
func isExcludedPath(path string) bool {
	for _, pattern := range config.ViperConfig.JWT.WhiteList {
		// 精确匹配或前缀匹配
		if path == pattern || (strings.HasSuffix(pattern, "*") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "*"))) {
			return true
		}
	}
	return false
}
