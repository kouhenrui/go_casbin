package config

type Config struct {
	Service  Service  `yaml:"service" json:"service" mapstructure:"service"`
	Casbin   Casbin   `yaml:"casbin" json:"casbin" mapstructure:"casbin"`
	Log      Log      `yaml:"log" json:"log" mapstructure:"log"`
	Database Database `yaml:"database" json:"database" mapstructure:"database"`
	Middleware MiddlewareConfig `yaml:"middleware" json:"middleware" mapstructure:"middleware"`
	JWT      JWT      `yaml:"jwt" json:"jwt" mapstructure:"jwt"`
}

type Service struct {
	Port string `yaml:"port" json:"port" mapstructure:"port"`
	Name string `yaml:"name" json:"name" mapstructure:"name"`
	Version string `yaml:"version" json:"version" mapstructure:"version"`
	Mode string `yaml:"mode" json:"mode" mapstructure:"mode"`
	Env string `yaml:"env" json:"env" mapstructure:"env"`
}

type Casbin struct {
	ModelPath  string `yaml:"modelPath" json:"modelPath" mapstructure:"modelPath"`
	Driver     string `yaml:"driver" json:"driver" mapstructure:"driver"`
	DataSource string `yaml:"dataSource" json:"dataSource" mapstructure:"dataSource"`
}

type JWT struct {
	SecretKey     string `yaml:"secretKey" json:"secretKey" mapstructure:"secretKey"`
	ExpireTime    int    `yaml:"expireTime" json:"expireTime" mapstructure:"expireTime"`       // 过期时间（小时）
	RefreshTime   int    `yaml:"refreshTime" json:"refreshTime" mapstructure:"refreshTime"`   // 刷新时间（天）
	Issuer        string `yaml:"issuer" json:"issuer" mapstructure:"issuer"`
	Audience      string `yaml:"audience" json:"audience" mapstructure:"audience"`
	TokenPrefix   string `yaml:"tokenPrefix" json:"tokenPrefix" mapstructure:"tokenPrefix"`
	RefreshPrefix string `yaml:"refreshPrefix" json:"refreshPrefix" mapstructure:"refreshPrefix"`
	WhiteList     []string `yaml:"whiteList" json:"whiteList" mapstructure:"whiteList"`
}

type Log struct {
	Level      string `yaml:"level" json:"level" mapstructure:"level"`
	Format     string `yaml:"format" json:"format" mapstructure:"format"`
	Output     string `yaml:"output" json:"output" mapstructure:"output"`
	LogDir     string `yaml:"logDir" json:"logDir" mapstructure:"logDir"`
	TimeFormat string `yaml:"timeFormat" json:"timeFormat" mapstructure:"timeFormat"`
	MaxSize    int    `yaml:"maxSize" json:"maxSize" mapstructure:"maxSize"`
	MaxBackups int    `yaml:"maxBackups" json:"maxBackups" mapstructure:"maxBackups"`
	MaxAge     int    `yaml:"maxAge" json:"maxAge" mapstructure:"maxAge"`
	Compress   bool   `yaml:"compress" json:"compress" mapstructure:"compress"`
	Caller     bool   `yaml:"caller" json:"caller" mapstructure:"caller"`
}

type Database struct {
	Driver          string `yaml:"driver" json:"driver" mapstructure:"driver"`
	Host            string `yaml:"host" json:"host" mapstructure:"host"`
	Port            int    `yaml:"port" json:"port" mapstructure:"port"`
	Username        string `yaml:"username" json:"username" mapstructure:"username"`
	Password        string `yaml:"password" json:"password" mapstructure:"password"`
	DBName          string `yaml:"dbname" json:"dbname" mapstructure:"dbname"`
	Charset         string `yaml:"charset" json:"charset" mapstructure:"charset"`
	ParseTime       bool   `yaml:"parseTime" json:"parseTime" mapstructure:"parseTime"`
	Loc             string `yaml:"loc" json:"loc" mapstructure:"loc"`
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns" mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime" mapstructure:"connMaxLifetime"`
}


// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 异常捕获
	UseRecovery        bool `yaml:"useRecovery" json:"useRecovery" mapstructure:"useRecovery"`
	DetailedRecovery   bool `yaml:"detailedRecovery" json:"detailedRecovery" mapstructure:"detailedRecovery"`
	
	// 跨域
	UseCORS           bool `yaml:"useCORS" json:"useCORS" mapstructure:"useCORS"`
	AllowedOrigins    []string `yaml:"allowedOrigins" json:"allowedOrigins" mapstructure:"allowedOrigins"`
	AllowedMethods    []string `yaml:"allowedMethods" json:"allowedMethods" mapstructure:"allowedMethods"`
	AllowedHeaders    []string `yaml:"allowedHeaders" json:"allowedHeaders" mapstructure:"allowedHeaders"`
	
	// 日志
	UseLogger         bool `yaml:"useLogger" json:"useLogger" mapstructure:"useLogger"`
	DetailedLogging   bool `yaml:"detailedLogging" json:"detailedLogging" mapstructure:"detailedLogging"`
	
	// 响应格式
	UseResponseFormat bool `yaml:"useResponseFormat" json:"useResponseFormat" mapstructure:"useResponseFormat"`
	
	// 错误处理
	UseErrorHandler   bool `yaml:"useErrorHandler" json:"useErrorHandler" mapstructure:"useErrorHandler"`
	
	// 响应日志
	UseResponseLogger bool `yaml:"useResponseLogger" json:"useResponseLogger" mapstructure:"useResponseLogger"`
	
	// 验证错误处理
	UseValidationErrorHandler bool `yaml:"useValidationErrorHandler" json:"useValidationErrorHandler" mapstructure:"useValidationErrorHandler"`
	
	// 限流错误处理
	UseRateLimitErrorHandler bool `yaml:"useRateLimitErrorHandler" json:"useRateLimitErrorHandler" mapstructure:"useRateLimitErrorHandler"`
}



// DevelopmentMiddlewareConfig 开发环境中间件配置
func DevelopmentMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		UseRecovery:              true,
		DetailedRecovery:         true,
		UseCORS:                  true,
		AllowedOrigins:           []string{},
		AllowedMethods:           []string{},
		AllowedHeaders:           []string{},
		UseLogger:                true,
		DetailedLogging:          true,
		UseResponseFormat:        true,
		UseErrorHandler:          true,
		UseResponseLogger:        true,
		UseValidationErrorHandler: true,
		UseRateLimitErrorHandler: false,
	}
}

// ProductionMiddlewareConfig 生产环境中间件配置
func ProductionMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		UseRecovery:              true,
		DetailedRecovery:         false,
		UseCORS:                  true,
		AllowedOrigins:           []string{"https://yourdomain.com","http://127.0.0.1:9000"},
		AllowedMethods:           []string{"GET", "POST", "PUT", "DELETE","PATCH"},
		AllowedHeaders:           []string{"Content-Type", "Authorization", "X-API-Key","X-Client-Version"},
		UseLogger:                true,
		DetailedLogging:          false,
		UseResponseFormat:        true,
		UseErrorHandler:          true,
		UseResponseLogger:        true,
		UseValidationErrorHandler: false,
		UseRateLimitErrorHandler: true,
	}
} 