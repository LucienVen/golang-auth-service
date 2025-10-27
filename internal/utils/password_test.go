package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用的全局变量
var (
	testPassword = "TestPassword123!"
	testHasher   = NewPasswordHasher()
)

// TestNewPasswordHasher 测试密码哈希器创建
func TestNewPasswordHasher(t *testing.T) {
	hasher := NewPasswordHasher()
	assert.NotNil(t, hasher)
	assert.IsType(t, &BcryptPasswordHasher{}, hasher)
}

// TestBcryptPasswordHasher_HashPassword 测试密码哈希生成
func TestBcryptPasswordHasher_HashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		errorType   error
	}{
		{
			name:        "有效密码",
			password:    testPassword,
			expectError: false,
		},
		{
			name:        "密码太短",
			password:    "short",
			expectError: true,
			errorType:   ErrPasswordTooShort,
		},
		{
			name:        "密码太长",
			password:    strings.Repeat("a", 100),
			expectError: true,
			errorType:   ErrPasswordTooLong,
		},
		{
			name:        "弱密码 - 只有小写字母",
			password:    "weakpassword",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "弱密码 - 只有大写字母",
			password:    "WEAKPASSWORD",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "弱密码 - 只有数字",
			password:    "12345678",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "有效密码 - 包含特殊字符",
			password:    "ValidPass123@",
			expectError: false,
		},
		{
			name:        "空密码",
			password:    "",
			expectError: true,
			errorType:   ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &BcryptPasswordHasher{}
			hash, err := hasher.HashPassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
				if tt.errorType != nil {
					assert.Contains(t, err.Error(), tt.errorType.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.Equal(t, 60, len(hash)) // bcrypt哈希固定长度
				assert.True(t, strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$"))
			}
		})
	}
}

// TestBcryptPasswordHasher_VerifyPassword 测试密码验证
func TestBcryptPasswordHasher_VerifyPassword(t *testing.T) {
	hasher := &BcryptPasswordHasher{}

	// 先生成一个哈希
	hash, err := hasher.HashPassword(testPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectError bool
		errorType   error
	}{
		{
			name:        "正确密码",
			password:    testPassword,
			hash:        hash,
			expectError: false,
		},
		{
			name:        "错误密码",
			password:    "WrongPassword123!",
			hash:        hash,
			expectError: true,
			errorType:   ErrPasswordVerifyFailed,
		},
		{
			name:        "空哈希",
			password:    testPassword,
			hash:        "",
			expectError: true,
			errorType:   ErrPasswordVerifyFailed,
		},
		{
			name:        "无效哈希格式",
			password:    testPassword,
			hash:        "invalid_hash",
			expectError: true,
			errorType:   ErrPasswordVerifyFailed,
		},
		{
			name:        "空密码",
			password:    "",
			hash:        hash,
			expectError: true,
			errorType:   ErrPasswordVerifyFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.VerifyPassword(tt.password, tt.hash)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.Contains(t, err.Error(), tt.errorType.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBcryptPasswordHasher_CheckPasswordStrength 测试密码强度检查
func TestBcryptPasswordHasher_CheckPasswordStrength(t *testing.T) {
	hasher := &BcryptPasswordHasher{}

	tests := []struct {
		name        string
		password    string
		expectError bool
		errorType   error
	}{
		{
			name:        "有效密码",
			password:    "ValidPass123",
			expectError: false,
		},
		{
			name:        "有效密码带特殊字符",
			password:    "ValidPass123@",
			expectError: false,
		},
		{
			name:        "密码太短",
			password:    "short",
			expectError: true,
			errorType:   ErrPasswordTooShort,
		},
		{
			name:        "密码太长",
			password:    strings.Repeat("a", 100),
			expectError: true,
			errorType:   ErrPasswordTooLong,
		},
		{
			name:        "只有小写字母",
			password:    "lowercaseonly",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "只有大写字母",
			password:    "UPPERCASEONLY",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "只有数字",
			password:    "12345678",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "缺少大写字母",
			password:    "lowercase123",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "缺少小写字母",
			password:    "UPPERCASE123",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "缺少数字",
			password:    "NoDigitsHere",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "边界长度测试 - 最小长度",
			password:    "Abcd1234",
			expectError: false,
		},
		{
			name:        "边界长度测试 - 最小长度减1",
			password:    "Abc1234",
			expectError: true,
			errorType:   ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.CheckPasswordStrength(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.Contains(t, err.Error(), tt.errorType.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGlobalFunctions 测试全局函数
func TestGlobalFunctions(t *testing.T) {
	t.Run("全局HashPassword函数", func(t *testing.T) {
		hash, err := HashPassword(testPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Equal(t, 60, len(hash))
	})

	t.Run("全局VerifyPassword函数", func(t *testing.T) {
		hash, err := HashPassword(testPassword)
		require.NoError(t, err)

		err = VerifyPassword(testPassword, hash)
		assert.NoError(t, err)

		err = VerifyPassword("wrong", hash)
		assert.Error(t, err)
	})

	t.Run("全局CheckPasswordStrength函数", func(t *testing.T) {
		err := CheckPasswordStrength(testPassword)
		assert.NoError(t, err)

		err = CheckPasswordStrength("weak")
		assert.Error(t, err)
	})
}

// TestIsPasswordHashValid 测试密码哈希格式验证
func TestIsPasswordHashValid(t *testing.T) {
	tests := []struct {
		name     string
		hash     string
		expected bool
	}{
		{
			name:     "有效的bcrypt哈希",
			hash:     "$2b$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW",
			expected: true,
		},
		{
			name:     "无效格式 - 错误前缀",
			hash:     "$2c$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW",
			expected: false,
		},
		{
			name:     "无效格式 - 长度错误",
			hash:     "$2b$12$short",
			expected: false,
		},
		{
			name:     "空字符串",
			hash:     "",
			expected: false,
		},
		{
			name:     "纯文本",
			hash:     "plaintext",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPasswordHashValid(tt.hash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGenerateRandomPassword 测试随机密码生成
func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "生成8位密码",
			length: 8,
		},
		{
			name:   "生成12位密码",
			length: 12,
		},
		{
			name:   "生成16位密码",
			length: 16,
		},
		{
			name:   "长度小于最小值",
			length: 4,
		},
		{
			name:   "长度大于最大值",
			length: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GenerateRandomPassword(tt.length)
			assert.NoError(t, err)
			assert.NotEmpty(t, password)

			// 检查密码长度
			expectedLength := tt.length
			if expectedLength < minPasswordLength {
				expectedLength = minPasswordLength
			}
			if expectedLength > maxPasswordLength {
				expectedLength = maxPasswordLength
			}
			assert.Equal(t, expectedLength, len(password))

			// 检查密码强度
			err = CheckPasswordStrength(password)
			assert.NoError(t, err)
		})
	}
}

// BenchmarkHashPassword 密码哈希生成性能测试
func BenchmarkHashPassword(b *testing.B) {
	hasher := &BcryptPasswordHasher{}
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hasher.HashPassword(password)
	}
}

// BenchmarkVerifyPassword 密码验证性能测试
func BenchmarkVerifyPassword(b *testing.B) {
	hasher := &BcryptPasswordHasher{}
	password := "BenchmarkPassword123!"
	hash, _ := hasher.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hasher.VerifyPassword(password, hash)
	}
}

// TestPasswordConsistency 测试密码哈希的一致性
func TestPasswordConsistency(t *testing.T) {
	hasher := &BcryptPasswordHasher{}
	password := "ConsistentPassword123!"

	// 多次生成哈希，应该得到不同的结果（因为每次使用不同的盐）
	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := hasher.HashPassword(password)
		require.NoError(t, err)
		hashes[i] = hash
	}

	// 所有哈希都应该不同
	for i := 1; i < len(hashes); i++ {
		assert.NotEqual(t, hashes[0], hashes[i])
	}

	// 但所有哈希都应该能验证同一个密码
	for _, hash := range hashes {
		err := hasher.VerifyPassword(password, hash)
		assert.NoError(t, err)
	}
}

// TestPasswordCharacterValidation 测试密码字符验证逻辑
func TestPasswordCharacterValidation(t *testing.T) {
	hasher := &BcryptPasswordHasher{}

	tests := []struct {
		name        string
		password    string
		expectError bool
		errorType   error
	}{
		{
			name:        "包含所有必需字符",
			password:    "ValidPass123",
			expectError: false,
		},
		{
			name:        "包含特殊字符",
			password:    "ValidPass123@",
			expectError: false,
		},
		{
			name:        "缺少大写字母",
			password:    "lowercase123",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "缺少小写字母",
			password:    "UPPERCASE123",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "缺少数字",
			password:    "NoDigitsHere",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "只有数字",
			password:    "12345678",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "只有小写字母",
			password:    "onlylowercase",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "只有大写字母",
			password:    "ONLYUPPERCASE",
			expectError: true,
			errorType:   ErrPasswordTooWeak,
		},
		{
			name:        "边界情况 - 最小有效密码",
			password:    "Abc12345",
			expectError: false,
		},
		{
			name:        "边界情况 - 包含Unicode字符",
			password:    "测试密码123ABC", // 包含中文，但实际实现可能无法识别为有效字符
			expectError: true,  // 由于当前实现只检查ASCII字母数字，Unicode字符会被忽略
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.CheckPasswordStrength(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.Contains(t, err.Error(), tt.errorType.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}