package main

import (
	"go_casbin/api"
	"go_casbin/internal/config"
	"go_casbin/internal/logger"
	"go_casbin/pkg/casbin"
	"go_casbin/pkg/jwt"
)

func init() {
	logger.Init(nil) // 初始化logger

	// 使用相对于项目根目录的配置文件路径
	config.InitConfig("configs/config.dev.yaml") // 初始化配置
	jwt.InitJWTConfig(nil)
	// 初始化casbin服务
	err := casbin.InitCasbin(casbin.CasbinOptions{
		Driver: config.ViperConfig.Casbin.Driver,
		DataSource: config.ViperConfig.Casbin.DataSource,
		ModelPath: config.ViperConfig.Casbin.ModelPath,
	})
	if err != nil {
		logger.ErrorWithErr("初始化CasbinService失败", err)
		panic(err)
	}

}

func main() {
	r := api.InitRouter()
	var port string = config.ViperConfig.Service.Port
	if err := r.Run(port); err != nil {
		logger.ErrorWithErr("服务启动失败", err, logger.String("port", port))
		panic(err)
	}
	logger.Info("服务启动成功", logger.String("port", port))
}
