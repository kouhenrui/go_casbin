package middleware

import (
	"go_casbin/internal/config"
	"go_casbin/internal/middleware/cors"
	errorhandler "go_casbin/internal/middleware/error"
	"go_casbin/internal/middleware/logger"
	"go_casbin/internal/middleware/recovery"
	"go_casbin/internal/middleware/response"

	"github.com/gin-gonic/gin"
)

// SetupMiddlewares 设置所有中间件
func SetupMiddlewares(r *gin.Engine) {
	// 1. 全局异常捕获（必须在最前面）
	r.Use(recovery.Recovery())
	
	// 2. 跨域中间件
	r.Use(cors.CORS())	
	
	// 3. 日志中间件
	r.Use(logger.Logger())
	
	// 4. 统一响应格式中间件
	r.Use(response.ResponseMiddleware())


	// 5. 错误处理中间件
	r.Use(errorhandler.ErrorHandler())
	
	// 6. 响应日志中间件
	r.Use(response.ResponseLogger())

}

// SetupDevelopmentMiddlewares 开发环境中间件
func SetupDevelopmentMiddlewares(r *gin.Engine) {
	// 开发环境使用详细日志
	r.Use(recovery.RecoveryWithLogger())// 全局异常捕获（必须在最前面）
	r.Use(cors.CORS())// 跨域中间件
	r.Use(logger.RequestLogger()) // 详细请求日志
	r.Use(response.ResponseMiddleware())// 统一响应格式中间件
	r.Use(errorhandler.ErrorHandler())// 错误处理中间件
	r.Use(response.ResponseLogger())// 响应日志中间件
	r.Use(errorhandler.ValidationErrorHandler())// 验证错误处理
	

}


// SetupProductionMiddlewares 生产环境中间件
func SetupProductionMiddlewares(r *gin.Engine) {
	// 生产环境使用精简日志
	r.Use(recovery.Recovery())// 全局异常捕获（必须在最前面）
	r.Use(cors.CORSWithConfig(
		[]string{"https://yourdomain.com","http://127.0.0.1:9000"}, // 只允许特定域名
		[]string{"GET", "POST", "PUT", "DELETE","PATCH"},
		[]string{"Content-Type", "Authorization", "X-API-Key","X-Client-Version","X-Platform"},
	))// 跨域中间件
	r.Use(logger.Logger())// 日志中间件
	r.Use(response.ResponseMiddleware())// 统一响应格式中间件
	r.Use(errorhandler.DefaultErrorHandler())// 错误处理中间件
	r.Use(errorhandler.RateLimitErrorHandler())// 限流错误处理
}

// SetupCustomMiddlewares 自定义中间件配置
func SetupCustomMiddlewares(r *gin.Engine, config config.MiddlewareConfig) {
	// 异常捕获
	if config.UseRecovery {
		if config.DetailedRecovery {
			r.Use(recovery.RecoveryWithLogger())
		} else {
			r.Use(recovery.Recovery())
		}
	}
	
	// 跨域
	if config.UseCORS {
		if len(config.AllowedOrigins) > 0 {
			r.Use(cors.CORS())
		} else {
			r.Use(cors.CORS())
		}
	}
	
	// 日志
	if config.UseLogger {
		if config.DetailedLogging {
			r.Use(logger.RequestLogger())
		} else {
			r.Use(logger.Logger())
		}
	}
	
	// 响应格式
	if config.UseResponseFormat {
		r.Use(response.ResponseMiddleware())
	}
	
	// 错误处理
	if config.UseErrorHandler {
		r.Use(errorhandler.ErrorHandler())
	}
	
	// 响应日志
	if config.UseResponseLogger {
		r.Use(response.ResponseLogger())
	}
	
	// 验证错误处理
	if config.UseValidationErrorHandler {
		r.Use(errorhandler.ValidationErrorHandler())
	}
	
	// 限流错误处理
	if config.UseRateLimitErrorHandler {
		r.Use(errorhandler.RateLimitErrorHandler())
	}
}
