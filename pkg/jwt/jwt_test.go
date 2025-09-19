package jwt

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	// 测试用例
	tests := []struct {
		name     string
		userID   string
		username string
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "正常生成Token",
			userID:   "123",
			username: "testuser",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "空用户ID",
			userID:   "",
			username: "testuser",
			duration: time.Hour,
			wantErr:  false,
		},
		{
			name:     "空用户名",
			userID:   "123",
			username: "",
			duration: time.Hour,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() token is empty")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	// 先生成一个有效的token
	userID := "123"
	username := "testuser"
	duration := time.Hour
	validToken, _ := GenerateToken(userID, username, duration)

	// 测试用例
	tests := []struct {
		name         string
		token        string
		wantUserID   string
		wantUsername string
		wantErr      bool
	}{
		{
			name:         "有效Token",
			token:        validToken,
			wantUserID:   userID,
			wantUsername: username,
			wantErr:      false,
		},
		{
			name:         "无效Token",
			token:        "invalid.token.here",
			wantUserID:   "",
			wantUsername: "",
			wantErr:      true,
		},
		{
			name:         "空Token",
			token:        "",
			wantUserID:   "",
			wantUsername: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims.UserID != tt.wantUserID {
					t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, tt.wantUserID)
				}
				if claims.Username != tt.wantUsername {
					t.Errorf("ParseToken() Username = %v, want %v", claims.Username, tt.wantUsername)
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// 先生成一个有效的token
	userID := "123"
	username := "testuser"
	oldDuration := time.Hour
	validToken, _ := GenerateToken(userID, username, oldDuration)
	newDuration := 2 * time.Hour

	// 测试用例
	tests := []struct {
		name     string
		token    string
		duration time.Duration
		wantErr  bool
	}{
		{
			name:     "刷新有效Token",
			token:    validToken,
			duration: newDuration,
			wantErr:  false,
		},
		{
			name:     "刷新无效Token",
			token:    "invalid.token.here",
			duration: newDuration,
			wantErr:  true,
		},
		{
			name:     "刷新空Token",
			token:    "",
			duration: newDuration,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := RefreshToken(tt.token, tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if newToken == "" {
					t.Error("RefreshToken() new token is empty")
				}
				// 验证新token
				claims, err := ParseToken(newToken)
				if err != nil {
					t.Errorf("ParseToken() error = %v", err)
					return
				}
				if claims.UserID != userID {
					t.Errorf("RefreshToken() UserID = %v, want %v", claims.UserID, userID)
				}
				if claims.Username != username {
					t.Errorf("RefreshToken() Username = %v, want %v", claims.Username, username)
				}
			}
		})
	}
}
