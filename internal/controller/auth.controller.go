package controller

import (
	"go_casbin/internal/middleware/response"
	"go_casbin/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface{
	Login(c *gin.Context)
}

type AuthControllerImpl struct{
	authService service.AccountService
}
func NewAuthController() AuthController{
	return &AuthControllerImpl{
		authService: service.NewAccountService(),
	}
}
func(a *AuthControllerImpl) Login(c *gin.Context){

	response.Success(c, "login success")
}