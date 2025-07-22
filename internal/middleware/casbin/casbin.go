package casbin

import (
	"go_casbin/internal/middleware/response"
	casbinService "go_casbin/pkg/casbin"
	"go_casbin/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		casbinService := casbinService.GetCasbinInstance()
		accountVal, _ := c.Get("account")
		account, _ := accountVal.(jwt.Account)
		if account.Role == nil {
			response.Forbidden(c, "无权限")
			c.Abort()
			return
		}
		ok,err:=casbinService.Enforce(account.Role[0], c.Request.Method, c.Request.URL.Path)
		if err != nil {
			response.InternalServerError(c, err.Error())
			c.Abort()
			return
		}
		if !ok {
			response.Forbidden(c, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}