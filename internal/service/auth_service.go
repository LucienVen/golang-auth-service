package service

import (
	"context"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/appcontext"
	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/LucienVen/golang-auth-service/internal/request"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/pkg/auth"
)

// AuthService 认证服务接口（简化版本）
type AuthService interface {
	// Login 用户登录
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)

	// Register 用户注册
	Register(ctx context.Context, req *request.RegisterRequest) (*response.RegisterResponse, error)

	// Logout 用户登出
	Logout(ctx context.Context, tokenString string) error

	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, refreshToken string) (*response.RefreshTokenResponse, error)

	// ValidateToken 验证令牌
	ValidateToken(ctx context.Context, tokenString string) (*response.UserResponse, error)

	// GetProfile 获取用户信息
	GetProfile(ctx context.Context, userID string) (*response.UserResponse, error)

	// GetSessions 获取用户会话列表
	GetSessions(ctx context.Context, userID string) ([]*response.SessionResponse, error)

	// RevokeSession 撤销指定会话
	RevokeSession(ctx context.Context, userID, sessionID string) error

	// RevokeAllSessions 撤销用户所有会话
	RevokeAllSessions(ctx context.Context, userID string) error
}

// AuthServiceImpl 认证服务实现（简化版本）
type AuthServiceImpl struct {
	userService    UserService
	jwtService     JWTService
	sessionService SessionService
	validator      *auth.TokenValidator // 统一验证器
}

// NewAuthService 创建认证服务实例
func NewAuthService(appCtx *appcontext.AppContext) AuthService {
	jwtService := NewJWTService()
	validator := auth.NewTokenValidator(jwtService)

	return &AuthServiceImpl{
		userService:    &UserServiceImpl{}, // 简化实现，不依赖外部依赖
		jwtService:     jwtService,
		sessionService: NewSessionService(appCtx),
		validator:      validator,
	}
}

// NewAuthServiceWithValidator 使用外部验证器创建认证服务实例
func NewAuthServiceWithValidator(validator *auth.TokenValidator, appCtx *appcontext.AppContext) AuthService {
	return &AuthServiceImpl{
		userService:    &UserServiceImpl{}, // 简化实现，不依赖外部依赖
		jwtService:     NewJWTService(),
		sessionService: NewSessionService(appCtx),
		validator:      validator,
	}
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	// 验证登录请求
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("account", err.Error(), "required")
	}

	// TODO: 实现用户查询和验证逻辑
	// 临时实现：生成测试令牌
	tokenPair, err := s.jwtService.GenerateTokenPair("test-user", "test-username")
	if err != nil {
		return nil, errors.NewTokenGenerationError("生成令牌失败: " + err.Error())
	}

	// 创建测试用户响应
	userResponse := &response.UserResponse{
		ID:       "test-user",
		Username: "test-username",
		Status:   1, // 活跃状态
	}

	// 构建响应
	loginResponse := &response.LoginResponse{
		User:         userResponse,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	return loginResponse, nil
}

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req *request.RegisterRequest) (*response.RegisterResponse, error) {
	// 验证注册请求
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("account", err.Error(), "required")
	}

	// TODO: 实现用户创建逻辑
	// 临时实现：生成测试令牌
	tokenPair, err := s.jwtService.GenerateTokenPair("new-user", "new-username")
	if err != nil {
		return nil, errors.NewTokenGenerationError("生成令牌失败: " + err.Error())
	}

	// 创建测试用户响应
	userResponse := &response.UserResponse{
		ID:       "new-user",
		Username: "new-username",
		Status:   0, // 未激活状态
	}

	// 构建响应
	registerResponse := &response.RegisterResponse{
		User:         userResponse,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	return registerResponse, nil
}

// Logout 用户登出
func (s *AuthServiceImpl) Logout(ctx context.Context, tokenString string) error {
	// 验证令牌
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return errors.NewInvalidTokenError("无效的令牌: " + err.Error())
	}

	// 将令牌加入黑名单
	if err := s.jwtService.RevokeToken(claims.ID); err != nil {
		return errors.NewTokenRevocationError("撤销令牌失败: " + err.Error())
	}

	// 使用会话服务将令牌加入黑名单
	if err := s.sessionService.AddTokenToBlacklist(ctx, claims.ID, time.Now().Add(time.Hour*24)); err != nil {
		return errors.NewSessionError("会话登出失败: " + err.Error())
	}

	return nil
}

// RefreshToken 刷新令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*response.RefreshTokenResponse, error) {
	// 使用JWT服务刷新令牌
	tokenPair, err := s.jwtService.RefreshToken(refreshToken)
	if err != nil {
		return nil, errors.NewTokenRefreshError("刷新令牌失败: " + err.Error())
	}

	// 构建响应
	refreshResponse := &response.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	return refreshResponse, nil
}

// ValidateToken 验证令牌
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*response.UserResponse, error) {
	// 使用统一验证器验证令牌
	return s.validator.ValidateTokenString(ctx, tokenString)
}

// GetProfile 获取用户信息
func (s *AuthServiceImpl) GetProfile(ctx context.Context, userID string) (*response.UserResponse, error) {
	// TODO: 实现用户信息获取
	// 临时实现：返回测试用户
	userResponse := &response.UserResponse{
		ID:       userID,
		Username: "test-user",
		Status:   1, // 活跃状态
	}

	return userResponse, nil
}

// GetSessions 获取用户会话列表
func (s *AuthServiceImpl) GetSessions(ctx context.Context, userID string) ([]*response.SessionResponse, error) {
	sessions, err := s.sessionService.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, errors.NewSessionError("获取用户会话失败: " + err.Error())
	}

	var sessionResponses []*response.SessionResponse
	for _, session := range sessions {
		sessionResponse := &response.SessionResponse{
			SessionID:  session.SessionID,
			DeviceID:   session.DeviceID,
			UserAgent:  session.UserAgent,
			IPAddress:  session.IPAddress,
			LoginTime:  session.LoginTime,
			LastActive: session.LastActive,
			ExpiresAt:  session.ExpiresAt,
			IsActive:   session.IsActive,
		}
		sessionResponses = append(sessionResponses, sessionResponse)
	}

	return sessionResponses, nil
}

// RevokeSession 撤销指定会话
func (s *AuthServiceImpl) RevokeSession(ctx context.Context, userID, sessionID string) error {
	// 将令牌加入黑名单
	if err := s.sessionService.AddTokenToBlacklist(ctx, sessionID, time.Now().Add(time.Hour*24)); err != nil {
		return errors.NewSessionError("撤销会话失败: " + err.Error())
	}

	return nil
}

// RevokeAllSessions 撤销用户所有会话
func (s *AuthServiceImpl) RevokeAllSessions(ctx context.Context, userID string) error {
	if err := s.sessionService.RevokeUserSessions(ctx, userID); err != nil {
		return errors.NewSessionError("撤销用户所有会话失败: " + err.Error())
	}

	return nil
}