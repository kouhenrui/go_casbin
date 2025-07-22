package api

import (
	"go_casbin/internal/logger"
	"go_casbin/internal/middleware/response"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	
	// API版本路由组
	v1 := r.Group("/api/v1")
	{
		// 测试接口
		v1.GET("/test", func(c *gin.Context) {
			logger.Info("测试接口被调用",
				logger.String("client_ip", c.ClientIP()),
				logger.String("user_agent", c.Request.UserAgent()),
			)
	
			response.Success(c,gin.H{
					"service": "go_casbin",
					"version": "1.0.0",
					"time":    "2024-01-01T00:00:00Z",
			})
		})
		// authController := controller.NewAuthController()
		// v1.POST("/login", authController.Login)
		
		// 错误测试接口
		v1.GET("/error-test", func(c *gin.Context) {
			// 测试不同类型的错误
			errorType := c.Query("type")
			
			switch errorType {
			case "400":
				response.BadRequest(c, "这是一个400错误测试")
			case "401":
				response.Unauthorized(c, "这是一个401错误测试")
			case "403":
				response.Forbidden(c, "这是一个403错误测试")
			case "500":
				response.InternalServerError(c, "这是一个500错误测试")
			case "panic":
				panic("这是一个panic测试")
			default:
				response.Success(c, gin.H{
					"message": "错误测试接口",
					"usage":   "使用 ?type=400|401|403|404|405|500|panic 来测试不同错误",
				})
			}
		})
		
		// 分页测试接口
		v1.GET("/pagination-test", func(c *gin.Context) {
			page := 1
			pageSize := 10
			total := int64(100)
			
			// 模拟数据
			data := []gin.H{
				{"id": 1, "name": "Item 1"},
				{"id": 2, "name": "Item 2"},
				{"id": 3, "name": "Item 3"},
			}
			
			response.PaginatedResponse(c, data, total, page, pageSize)
		})
		//模拟接口错误
		v1.GET("/error", func(c *gin.Context) {
			// middleware.Error(c,666,"模拟错误发出")
			panic("这是一个panic测试")
		})	
	}
}