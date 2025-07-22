package config

import (
	"os"
	"strconv"
	"strings"
)

// EnvConfig 环境变量配置
type EnvConfig struct {
	// 服务配置
	ServicePort string
	ServiceHost string
	ServiceMode string // debug, release, test
	
	// 数据库配置
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	// Redis配置
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
	
	// 日志配置
	LogLevel    string
	LogPath     string
	LogMaxSize  int
	LogMaxAge   int
	LogCompress bool
	
	// Casbin配置
	CasbinModelPath string
	CasbinDriver    string
	CasbinDataSource string
}

// LoadEnvConfig 从环境变量加载配置
func LoadEnvConfig() *EnvConfig {
	config := &EnvConfig{}
	
	// 服务配置
	config.ServicePort = getEnv("SERVICE_PORT", "8080")
	config.ServiceHost = getEnv("SERVICE_HOST", "0.0.0.0")
	config.ServiceMode = getEnv("SERVICE_MODE", "debug")
	
	// 数据库配置
	config.DBHost = getEnv("DB_HOST", "localhost")
	config.DBPort = getEnvAsInt("DB_PORT", 5432)
	config.DBUser = getEnv("DB_USER", "postgres")
	config.DBPassword = getEnv("DB_PASSWORD", "")
	config.DBName = getEnv("DB_NAME", "go_casbin")
	config.DBSSLMode = getEnv("DB_SSL_MODE", "disable")
	
	// Redis配置
	config.RedisHost = getEnv("REDIS_HOST", "localhost")
	config.RedisPort = getEnvAsInt("REDIS_PORT", 6379)
	config.RedisPassword = getEnv("REDIS_PASSWORD", "")
	config.RedisDB = getEnvAsInt("REDIS_DB", 0)
	
	// 日志配置
	config.LogLevel = getEnv("LOG_LEVEL", "info")
	config.LogPath = getEnv("LOG_PATH", "logs")
	config.LogMaxSize = getEnvAsInt("LOG_MAX_SIZE", 100)
	config.LogMaxAge = getEnvAsInt("LOG_MAX_AGE", 7)
	config.LogCompress = getEnvAsBool("LOG_COMPRESS", true)
	
	// Casbin配置
	config.CasbinModelPath = getEnv("CASBIN_MODEL_PATH", "configs/model.conf")
	config.CasbinDriver = getEnv("CASBIN_DRIVER", "file")
	config.CasbinDataSource = getEnv("CASBIN_DATA_SOURCE", "configs/policy.csv")
	
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为bool
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(strings.ToLower(value)); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsFloat 获取环境变量并转换为float64
func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// SetEnv 设置环境变量
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetEnv 删除环境变量
func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}

// GetAllEnv 获取所有环境变量
func GetAllEnv() map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envs[pair[0]] = pair[1]
		}
	}
	return envs
}

// GetEnvWithPrefix 获取指定前缀的环境变量
func GetEnvWithPrefix(prefix string) map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 && strings.HasPrefix(pair[0], prefix) {
			envs[pair[0]] = pair[1]
		}
	}
	return envs
} 