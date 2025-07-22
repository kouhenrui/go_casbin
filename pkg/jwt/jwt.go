package jwt

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var (
	once sync.Once
	jwtConfig *JWTConfig
)
// JWTClaims JWT声明结构
type JWTClaims struct {
	Account Account `json:"account" mapstructure:"account"`
	jwt.RegisteredClaims
}
type Account struct {
	ID       string `json:"id"`//用户id
	Username string `json:"username"`//用户名
	Role     []string `json:"role"`//角色
	Permissions []string `json:"permissions,omitempty"`//权限
	OpenId     string `json:"open_id,omitempty"`//openid
	Platform   string `json:"platform,omitempty"`//平台
	SystemId   string `json:"system_id,omitempty"`//系统id
	TenantId   string `json:"tenant_id,omitempty"`//租户id
	AppId      string `json:"app_id,omitempty"`//应用id
	Status    int8 `json:"status,omitempty"`//状态
	IsVerified bool `json:"is_verified,omitempty"`//是否验证
	IsLocked   bool `json:"is_locked,omitempty"`//是否锁定
}
// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string        `json:"secret_key"`     // 密钥
	ExpireTime    time.Duration `json:"expire_time"`    // 过期时间
	RefreshTime   time.Duration `json:"refresh_time"`   // 刷新时间
	Issuer        string        `json:"issuer"`         // 签发者
	Audience      string        `json:"audience,omitempty"`       // 受众
	TokenPrefix   string        `json:"token_prefix,omitempty"`   // Token前缀
	RefreshPrefix string        `json:"refresh_prefix,omitempty"` // 刷新Token前缀
}

// DefaultJWTConfig 默认JWT配置
func defaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:     "your-secret-key-change-in-production",
		ExpireTime:    24 * time.Hour,
		RefreshTime:   7 * 24 * time.Hour,
		Issuer:        "go_casbin",
		Audience:      "go_casbin-api",
		TokenPrefix:   "Bearer ",
		RefreshPrefix: "Refresh ",
	}
}
func InitJWTConfig(option *JWTConfig) {
	once.Do(func() {
		jwtConfig =defaultJWTConfig()
		if option != nil {
			if option.SecretKey != "" {
				jwtConfig.SecretKey = option.SecretKey
			}
			if option.ExpireTime != 0 {
				jwtConfig.ExpireTime = option.ExpireTime
			}
			if option.RefreshTime != 0 {
				jwtConfig.RefreshTime = option.RefreshTime
			}
			if option.Issuer != "" {
				jwtConfig.Issuer = option.Issuer
			}
			if option.Audience != "" {
				jwtConfig.Audience = option.Audience
			}
			if option.TokenPrefix != "" {
				jwtConfig.TokenPrefix = option.TokenPrefix
			}
			if option.RefreshPrefix != "" {
				jwtConfig.RefreshPrefix = option.RefreshPrefix
			}
	}
	})
}

func GetJWTInstance() *JWTConfig {
	return jwtConfig
}

// GenerateJWTToken 生成JWT Token
func(j *JWTConfig) GenerateJWTToken(payload Account) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		Account: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.Issuer,
			Audience:  []string{j.Audience},
			Subject:   payload.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// GenerateRefreshToken 生成刷新Token
func(j *JWTConfig) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(j.RefreshTime)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    j.Issuer,
		Audience:  []string{j.Audience + "-refresh"},
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// ParseToken 解析JWT Token
func(j *JWTConfig) ParseToken(tokenString string) (*Account, error) {	
	// 移除前缀
	if len(tokenString) > len(j.TokenPrefix) && tokenString[:len(j.TokenPrefix)] == j.TokenPrefix {
		tokenString = tokenString[len(j.TokenPrefix):]
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &claims.Account, nil
	}

	return nil, errors.New("invalid token")
}

// ParseRefreshToken 解析刷新Token
func(j *JWTConfig) parseRefreshToken(tokenString string) (*JWTClaims, error) {
// 移除前缀
	if len(tokenString) > len(j.RefreshPrefix) && tokenString[:len(j.RefreshPrefix)] == j.RefreshPrefix {
		tokenString = tokenString[len(j.RefreshPrefix):]
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid refresh token")
}

// ParseRefreshToken 解析刷新Token
func(j *JWTConfig) ParseRefreshToken(tokenString string) (*Account, error) {
	claims,err:= j.parseRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &claims.Account, nil
}

// ValidateToken 验证Token是否有效
func(j *JWTConfig) ValidateToken(tokenString string) (bool, error) {
	_, err := j.ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetTokenExpiration 获取Token过期时间
func(j *JWTConfig) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := j.parseRefreshToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

// IsTokenExpired 检查Token是否过期
func(j *JWTConfig) IsTokenExpired(tokenString string) (bool, error) {
	expiration, err := j.GetTokenExpiration(tokenString)
	if err != nil {
		return true, err
	}
	return time.Now().After(expiration), nil
}

// RefreshTokenPair 刷新Token对
func(j *JWTConfig) RefreshTokenPair(refreshToken string) (string, string, error) {
	// 解析刷新Token获取用户ID
	account, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 生成新的Token对
	newToken, err := j.GenerateJWTToken(*account)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := j.GenerateRefreshToken(account.ID)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
} 