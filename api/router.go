package api

import (
	"go_casbin/internal/config"
	"go_casbin/internal/logger"
	"go_casbin/internal/middleware"
	errorhandler "go_casbin/internal/middleware/error"
	"go_casbin/internal/middleware/response"
	"time"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由
func InitRouter() *gin.Engine {
	r := gin.New()
	// 设置Gin配置
	r.HandleMethodNotAllowed = true
	r.MaxMultipartMemory = 20 << 20

	// 根据环境设置中间件
	if config.ViperConfig.Service.Env == "production" {
		if config.ViperConfig.Service.Mode == "release" {
			gin.SetMode(gin.ReleaseMode) //生产环境 控制台隐藏路由信息
		}
		middleware.SetupProductionMiddlewares(r)
	} else {
		gin.SetMode(gin.DebugMode) //开发环境 控制台显示路由信息
		middleware.SetupDevelopmentMiddlewares(r)
		logger.Info("开发环境中间件设置完成")
	}

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"service": config.ViperConfig.Service.Name,
			"version": config.ViperConfig.Service.Version,
			"timestamp": gin.H{
				"start_time": time.Now().Format(time.RFC3339),
				"current":    time.Now().Format(time.RFC3339),
			},
		})
	})
	RegisterRoutes(r) //挂载API

	// 注册404和405错误处理（必须在所有路由注册完成后）
	r.NoRoute(errorhandler.NoRoute())
	r.NoMethod(errorhandler.NoMethod())

	logger.Info("路由初始化完成")
	return r
}
