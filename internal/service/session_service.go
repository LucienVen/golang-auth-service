package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/appcontext"
)

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	JTI          string    `json:"jti"`
	DeviceID     string    `json:"device_id"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastUsedAt   time.Time `json:"last_used_at"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SessionID  string    `json:"session_id"`
	UserID     string    `json:"user_id"`
	DeviceID   string    `json:"device_id"`
	UserAgent  string    `json:"user_agent"`
	IPAddress  string    `json:"ip_address"`
	LoginTime  time.Time `json:"login_time"`
	LastActive time.Time `json:"last_active"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsActive   bool      `json:"is_active"`
}

// SessionService 会话管理服务接口（简化版本）
type SessionService interface {
	// StoreToken 存储用户令牌信息
	StoreToken(ctx context.Context, userID string, tokenInfo *TokenInfo) error

	// GetToken 获取用户当前活跃令牌
	GetToken(ctx context.Context, userID string) (*TokenInfo, error)

	// RemoveToken 移除用户令牌
	RemoveToken(ctx context.Context, userID string) error

	// AddTokenToBlacklist 将令牌加入黑名单
	AddTokenToBlacklist(ctx context.Context, jti string, expiresAt time.Time) error

	// IsTokenInBlacklist 检查令牌是否在黑名单中
	IsTokenInBlacklist(ctx context.Context, jti string) (bool, error)

	// RevokeUserSessions 撤销用户所有会话
	RevokeUserSessions(ctx context.Context, userID string) error

	// GetUserSessions 获取用户所有活跃会话
	GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error)

	// GetActiveUsers 获取活跃用户数量
	GetActiveUsers(ctx context.Context) (int64, error)

	// CleanupExpiredSessions 清理过期会话
	CleanupExpiredSessions(ctx context.Context) error
}

// SessionServiceImpl 会话服务实现（简化版本）
type SessionServiceImpl struct {
	// 暂时不依赖 Redis，使用内存存储（生产环境需要 Redis）
	tokenStore     map[string]*TokenInfo
	blacklistStore map[string]bool
}

// NewSessionService 创建会话服务实例
func NewSessionService(appCtx *appcontext.AppContext) SessionService {
	return &SessionServiceImpl{
		tokenStore:     make(map[string]*TokenInfo),
		blacklistStore: make(map[string]bool),
	}
}

// StoreToken 存储用户令牌信息
func (s *SessionServiceImpl) StoreToken(ctx context.Context, userID string, tokenInfo *TokenInfo) error {
	if tokenInfo == nil {
		return fmt.Errorf("令牌信息不能为空")
	}

	// 临时内存存储
	s.tokenStore[userID] = tokenInfo
	return nil
}

// GetToken 获取用户当前活跃令牌
func (s *SessionServiceImpl) GetToken(ctx context.Context, userID string) (*TokenInfo, error) {
	if token, ok := s.tokenStore[userID]; ok {
		return token, nil
	}
	return nil, fmt.Errorf("用户令牌不存在")
}

// RemoveToken 移除用户令牌
func (s *SessionServiceImpl) RemoveToken(ctx context.Context, userID string) error {
	delete(s.tokenStore, userID)
	return nil
}

// AddTokenToBlacklist 将令牌加入黑名单
func (s *SessionServiceImpl) AddTokenToBlacklist(ctx context.Context, jti string, expiresAt time.Time) error {
	s.blacklistStore[jti] = true
	return nil
}

// IsTokenInBlacklist 检查令牌是否在黑名单中
func (s *SessionServiceImpl) IsTokenInBlacklist(ctx context.Context, jti string) (bool, error) {
	return s.blacklistStore[jti], nil
}

// RevokeUserSessions 撤销用户所有会话
func (s *SessionServiceImpl) RevokeUserSessions(ctx context.Context, userID string) error {
	// 删除用户令牌
	delete(s.tokenStore, userID)
	return nil
}

// GetUserSessions 获取用户所有活跃会话
func (s *SessionServiceImpl) GetUserSessions(ctx context.Context, userID string) ([]*SessionInfo, error) {
	if token, ok := s.tokenStore[userID]; ok {
		sessionInfo := &SessionInfo{
			SessionID:  token.JTI,
			UserID:     userID,
			DeviceID:   token.DeviceID,
			UserAgent:  token.UserAgent,
			IPAddress:  token.IPAddress,
			LoginTime:  token.CreatedAt,
			LastActive: token.LastUsedAt,
			ExpiresAt:  token.ExpiresAt,
			IsActive:   true,
		}
		return []*SessionInfo{sessionInfo}, nil
	}
	return []*SessionInfo{}, nil
}

// GetActiveUsers 获取活跃用户数量
func (s *SessionServiceImpl) GetActiveUsers(ctx context.Context) (int64, error) {
	return int64(len(s.tokenStore)), nil
}

// CleanupExpiredSessions 清理过期会话
func (s *SessionServiceImpl) CleanupExpiredSessions(ctx context.Context) error {
	now := time.Now()

	// 清理过期令牌
	for userID, token := range s.tokenStore {
		if now.After(token.ExpiresAt) {
			delete(s.tokenStore, userID)
		}
	}

	return nil
}