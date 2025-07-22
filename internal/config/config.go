// internal/config/config.go
package config

import (
	"go_casbin/internal/logger"
	"go_casbin/pkg/path"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ViperConfig Config
var v *viper.Viper
var once sync.Once

// InitConfig 初始化配置，单例模式
func InitConfig(configPath string) {
	once.Do(func() {
		v = viper.New()

		// 获取项目根目录的绝对路径
		configAbsPath, err := path.GetAbsolutePath(configPath)
		if err != nil {
			logger.Log().Fatal("获取配置文件路径失败", logger.Field("error", err))
			panic(err)
		}
		v.SetConfigFile(configAbsPath)
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			logger.Log().Fatal("读取配置文件失败", logger.Field("error", err), logger.Field("path", configAbsPath))
			panic(err)
		}
		logger.Log().Info("配置文件加载成功")
		viperLoadConf()
		v.WatchConfig() //开启监听
		v.OnConfigChange(func(in fsnotify.Event) {
			viperLoadConf() // 加载配置的方法
		})

	})

}

func viperLoadConf() {
	if err := v.Unmarshal(&ViperConfig); err != nil {
		log.Fatalf("Unable to decode into config struct: %v", err)
	}
}
