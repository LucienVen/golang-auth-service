package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/pkg/jwt"
)

// JWTService JWT服务接口（避免依赖service包）
type JWTService interface {
	ValidateToken(tokenString string) (*jwt.Claims, error)
}

// TokenValidator 令牌验证器，统一管理验证逻辑
type TokenValidator struct {
	jwtService JWTService
}

// NewTokenValidator 创建令牌验证器
func NewTokenValidator(jwtService JWTService) *TokenValidator {
	return &TokenValidator{
		jwtService: jwtService,
	}
}

// ValidateTokenString 验证令牌字符串
func (v *TokenValidator) ValidateTokenString(ctx context.Context, tokenString string) (*response.UserResponse, error) {
	if tokenString == "" {
		return nil, errors.NewInvalidTokenError("令牌不能为空")
	}

	// 验证JWT令牌
	claims, err := v.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.NewInvalidTokenError("无效的令牌: " + err.Error())
	}

	// 构建用户响应
	userResponse := &response.UserResponse{
		ID:       claims.UserID,
		Username: claims.Username,
		Status:   1, // 活跃状态
	}

	return userResponse, nil
}

// ValidateTokenFromRequest 从HTTP请求中验证令牌
func (v *TokenValidator) ValidateTokenFromRequest(ctx context.Context, req *http.Request) (*response.UserResponse, error) {
	// 从Authorization头中提取令牌
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.NewInvalidTokenError("缺少认证令牌")
	}

	// 检查Bearer前缀
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, errors.NewInvalidTokenError("无效的令牌格式")
	}

	// 提取令牌
	tokenString := authHeader[len(bearerPrefix):]
	if tokenString == "" {
		return nil, errors.NewInvalidTokenError("令牌不能为空")
	}

	// 验证令牌
	return v.ValidateTokenString(ctx, tokenString)
}