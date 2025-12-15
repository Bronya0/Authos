package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey  string
	ExpireTime time.Duration
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	AppID    uint   `json:"appId"`   // 应用数字ID，用于数据库操作
	AppUUID  string `json:"appUuid"` // 应用UUID，用于对外传输和引用
	jwt.RegisteredClaims
}

// SystemJWTClaims 系统管理员JWT声明
type SystemJWTClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

// AppJWTClaims 应用JWT声明
type AppJWTClaims struct {
	AppID   uint   `json:"appId"`
	AppUUID string `json:"appUuid"`
	AppCode string `json:"appCode"`
	jwt.RegisteredClaims
}

// NewJWTConfig 创建JWT配置实例
func NewJWTConfig(secretKey string, expireTime time.Duration) *JWTConfig {
	return &JWTConfig{
		SecretKey:  secretKey,
		ExpireTime: expireTime,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWTConfig) GenerateToken(userID uint, username string, appID uint, appUUID string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		AppID:    appID,
		AppUUID:  appUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// GenerateSystemToken 生成系统管理员JWT令牌
func (j *JWTConfig) GenerateSystemToken(username string) (string, error) {
	claims := SystemJWTClaims{
		Username: username,
		IsAdmin:  true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// GenerateAppToken 生成应用JWT令牌
func (j *JWTConfig) GenerateAppToken(appID uint, appUUID, appCode string) (string, error) {
	claims := AppJWTClaims{
		AppID:   appID,
		AppUUID: appUUID,
		AppCode: appCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// ParseToken 解析JWT令牌
func (j *JWTConfig) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseSystemToken 解析系统管理员JWT令牌
func (j *JWTConfig) ParseSystemToken(tokenString string) (*SystemJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SystemJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SystemJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid system token")
}

// ParseAppToken 解析应用JWT令牌
func (j *JWTConfig) ParseAppToken(tokenString string) (*AppJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AppJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid app token")
}
