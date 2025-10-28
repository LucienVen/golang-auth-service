package service

import (
	"fmt"
	"time"

	"github.com/LucienVen/golang-auth-service/pkg/jwt"
)

// JWTService JWT令牌管理服务接口
type JWTService interface {
	// GenerateTokenPair 生成访问令牌和刷新令牌对
	GenerateTokenPair(userID string, username string) (*TokenPair, error)

	// ValidateToken 验证访问令牌
	ValidateToken(tokenString string) (*jwt.Claims, error)

	// RefreshToken 使用刷新令牌生成新的令牌对
	RefreshToken(refreshTokenString string) (*TokenPair, error)

	// RevokeToken 撤销令牌（加入黑名单）
	RevokeToken(jti string) error

	// IsTokenRevoked 检查令牌是否已被撤销
	IsTokenRevoked(jti string) (bool, error)

	// ExtractJTI 从令牌中提取JTI
	ExtractJTI(tokenString string) (string, error)
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// JWTServiceImpl JWT服务实现
type JWTServiceImpl struct {
	accessExpiresIn  time.Duration // 访问令牌过期时间
	refreshExpiresIn time.Duration // 刷新令牌过期时间
	issuer          string        // 发行者
}

// NewJWTService 创建JWT服务实例
func NewJWTService() JWTService {
	return &JWTServiceImpl{
		accessExpiresIn:  time.Hour * 2,      // 访问令牌2小时过期
		refreshExpiresIn: time.Hour * 24 * 7, // 刷新令牌7天过期
		issuer:          "golang-auth-service",
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (s *JWTServiceImpl) GenerateTokenPair(userID string, username string) (*TokenPair, error) {
	now := time.Now()

	// 生成访问令牌
	accessToken, accessClaims, err := jwt.GenerateToken(userID, username, s.accessExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken, refreshClaims, err := jwt.GenerateToken(userID, username, s.refreshExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 创建令牌对
	tokenPair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessExpiresIn.Seconds()),
		ExpiresAt:    now.Add(s.accessExpiresIn),
	}

	// 这里可以添加令牌存储逻辑（如果需要的话）
	// 比如存储JTI到Redis用于黑名单检查
	_ = accessClaims.ID
	_ = refreshClaims.ID

	return tokenPair, nil
}

// ValidateToken 验证访问令牌
func (s *JWTServiceImpl) ValidateToken(tokenString string) (*jwt.Claims, error) {
	// 解析和验证令牌
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("令牌解析失败: %w", err)
	}

	// 检查令牌是否在黑名单中
	if revoked, err := s.IsTokenRevoked(claims.ID); err != nil {
		return nil, fmt.Errorf("检查令牌状态失败: %w", err)
	} else if revoked {
		return nil, fmt.Errorf("令牌已被撤销")
	}

	return claims, nil
}

// RefreshToken 使用刷新令牌生成新的令牌对
func (s *JWTServiceImpl) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := jwt.ParseToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("刷新令牌无效: %w", err)
	}

	// 检查令牌是否在黑名单中
	if revoked, err := s.IsTokenRevoked(claims.ID); err != nil {
		return nil, fmt.Errorf("检查令牌状态失败: %w", err)
	} else if revoked {
		return nil, fmt.Errorf("刷新令牌已被撤销")
	}

	// 撤销旧的刷新令牌
	if err := s.RevokeToken(claims.ID); err != nil {
		return nil, fmt.Errorf("撤销旧令牌失败: %w", err)
	}

	// 生成新的令牌对
	return s.GenerateTokenPair(claims.UserID, claims.Username)
}

// RevokeToken 撤销令牌（加入黑名单）
func (s *JWTServiceImpl) RevokeToken(jti string) error {
	// TODO: 实现令牌黑名单逻辑
	// 这里应该将JTI添加到Redis黑名单中
	// Redis键格式: auth:blacklist:{jti} = 1
	// 过期时间应该与令牌过期时间一致

	// 目前返回nil，表示成功
	// 在后续实现会话服务时，会在这里集成Redis黑名单功能
	return nil
}

// IsTokenRevoked 检查令牌是否已被撤销
func (s *JWTServiceImpl) IsTokenRevoked(jti string) (bool, error) {
	// TODO: 实现令牌黑名单检查
	// 这里应该检查Redis中是否存在该JTI
	// Redis键: auth:blacklist:{jti}

	// 目前返回false，表示未被撤销
	// 在后续实现会话服务时，会在这里集成Redis黑名单检查
	return false, nil
}

// ExtractJTI 从令牌中提取JTI
func (s *JWTServiceImpl) ExtractJTI(tokenString string) (string, error) {
	return jwt.GetJTI(tokenString)
}

// JWTConfig JWT配置选项
type JWTConfig struct {
	AccessExpiresIn  time.Duration // 访问令牌过期时间
	RefreshExpiresIn time.Duration // 刷新令牌过期时间
	Issuer          string        // 发行者
}

// NewJWTServiceWithConfig 使用配置创建JWT服务实例
func NewJWTServiceWithConfig(config JWTConfig) JWTService {
	// 设置默认值
	if config.AccessExpiresIn == 0 {
		config.AccessExpiresIn = time.Hour * 2
	}
	if config.RefreshExpiresIn == 0 {
		config.RefreshExpiresIn = time.Hour * 24 * 7
	}
	if config.Issuer == "" {
		config.Issuer = "golang-auth-service"
	}

	return &JWTServiceImpl{
		accessExpiresIn:  config.AccessExpiresIn,
		refreshExpiresIn: config.RefreshExpiresIn,
		issuer:           config.Issuer,
	}
}