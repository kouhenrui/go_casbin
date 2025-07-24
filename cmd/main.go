package main

import (
	"go_casbin/api"
	"go_casbin/internal/config"
	"go_casbin/internal/logger"
	"go_casbin/pkg/casbin"
	"go_casbin/pkg/database"
	"go_casbin/pkg/etcd"
	"go_casbin/pkg/jwt"
	"go_casbin/pkg/redis"
)

func init() {
	logger.Init(nil) // 初始化logger

	// 使用相对于项目根目录的配置文件路径
	config.InitConfig("configs/config.dev.yaml") // 初始化配置
	// 初始化pg数据库连接
	database.InitDB(database.DBOption{
		DBType: config.ViperConfig.Database.Driver,
		DBName: config.ViperConfig.Database.DBName,
		Username: config.ViperConfig.Database.Username,
		Password: config.ViperConfig.Database.Password,
		Host: config.ViperConfig.Database.Host,
		Port: config.ViperConfig.Database.Port,
		Charset: config.ViperConfig.Database.Charset,
		ParseTime: &config.ViperConfig.Database.ParseTime,
		Loc: &config.ViperConfig.Database.Loc,
	})
	// 初始化redis连接
	redis.InitRedis(redis.RedisOptions{
		Addr: config.ViperConfig.Redis.Addr,
		Username: config.ViperConfig.Redis.Username,
		Password: config.ViperConfig.Redis.Password,
	})
	// 初始化jwt配置
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
	// 初始化etcd连接
	etcd.InitEtcd(etcd.EtcdOptions{
		Endpoints: config.ViperConfig.Etcd.Endpoints,
		DialTimeout: config.ViperConfig.Etcd.DialTimeout,
		Username: config.ViperConfig.Etcd.Username,
		Password: config.ViperConfig.Etcd.Password,
	})

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
