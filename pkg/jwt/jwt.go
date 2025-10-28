package jwt

// jwt 工具封装

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义 JWT 的密钥（实际项目中建议通过配置文件或环境变量管理）
var jwtSecret = []byte("your_secret_key")

// Claims 结构体，嵌入 jwt.RegisteredClaims，可根据需要添加自定义字段
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// UserClaims 别名，用于向后兼容
type UserClaims = Claims

// generateJTI 生成唯一的 JTI
func generateJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID, username string, duration time.Duration) (string, *Claims, error) {
	jti := generateJTI()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, err
	}
	return tokenString, &claims, nil
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetJTI 获取令牌的 JTI
func GetJTI(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.ID, nil
}

// RefreshToken 刷新 Token（生成新的 Token，延长有效期）
func RefreshToken(tokenString string, duration time.Duration) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	tokenString, _, err = GenerateToken(claims.UserID, claims.Username, duration)
	return tokenString, err
}
