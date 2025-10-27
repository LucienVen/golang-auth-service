package utils

import (
	"testing"
)

// ========== 邮箱验证测试 ==========

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"有效邮箱1", "test@example.com", false},
		{"有效邮箱2", "user.name@domain.co.uk", false},
		{"有效邮箱3", "user+tag@example.org", false},
		{"有效邮箱4", "123@example.com", false},
		{"空邮箱", "", true},
		{"无效格式1", "test@", true},
		{"无效格式2", "@example.com", true},
		{"无效格式3", "test.example.com", true},
		{"连续的点", "test..email@example.com", true},
		{"以点开头", ".test@example.com", true},
		{"以点结尾", "test.@example.com", true},
		{"无效域名", "test@.com", true},
		{"过长邮箱", string(make([]byte, 101))+"@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"", false},
		{"invalid-email", false},
		{"test@.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== 手机号验证测试 ==========

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"有效手机号1", "13812345678", false},
		{"有效手机号2", "15000000000", false},
		{"有效手机号3", "19999999999", false},
		{"空手机号", "", true},
		{"无效开头", "12812345678", true},
		{"长度不足", "1381234567", true},
		{"长度过长", "138123456789", true},
		{"包含非数字", "1381234567a", true},
		{"包含空格", "138 1234 5678", true},
		{"固定电话", "02112345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		phone string
		want  bool
	}{
		{"13812345678", true},
		{"", false},
		{"12812345678", false},
		{"1381234567", false},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			if got := IsValidPhone(tt.phone); got != tt.want {
				t.Errorf("IsValidPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== 用户名验证测试 ==========

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有效用户名1", "user123", false},
		{"有效用户名2", "test_user", false},
		{"有效用户名3", "abc", false},
		{"有效用户名4", "User_Name_123", false},
		{"空用户名", "", true},
		{"太短", "ab", true},
		{"太长", string(make([]byte, 51)), true},
		{"下划线开头", "_username", true},
		{"下划线结尾", "username_", true},
		{"全下划线", "___", true},
		{"包含特殊字符", "user@name", true},
		{"包含空格", "user name", true},
		{"包含中文", "用户名", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		want     bool
	}{
		{"user123", true},
		{"", false},
		{"_username", false},
		{"user@name", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := IsValidUsername(tt.username); got != tt.want {
				t.Errorf("IsValidUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== 密码验证测试 ==========

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"有效密码1", "Password123", false},
		{"有效密码2", "Test1234", false},
		{"有效密码3", "MySecurePass1", false},
		{"空密码", "", true},
		{"太短", "Pass1", true},
		{"没有大写字母", "password123", true},
		{"没有小写字母", "PASSWORD123", true},
		{"没有数字", "Password", true},
		{"太长", string(make([]byte, 73)) + "A1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ========== 账户验证测试 ==========

func TestValidateAccount(t *testing.T) {
	tests := []struct {
		name    string
		account string
		wantErr bool
	}{
		{"邮箱账户", "test@example.com", false},
		{"手机号账户", "13812345678", false},
		{"用户名账户", "testuser123", false},
		{"空账户", "", true},
		{"无效格式", "invalid@account", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccount(tt.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetectAccountType(t *testing.T) {
	tests := []struct {
		account string
		want    string
	}{
		{"test@example.com", "email"},
		{"13812345678", "phone"},
		{"testuser123", "username"},
		{"ab", "unknown"},        // 太短的用户名（少于3个字符）
		{"_invalid", "unknown"},  // 以下划线开头的用户名
		{"invalid_", "unknown"},  // 以下划线结尾的用户名
		{"user@name", "unknown"}, // 包含特殊字符的用户名
		{"", "unknown"},
		{"test@", "unknown"}, // 无效邮箱格式
		{"12345", "username"}, // 纯数字用户名是有效的
	}

	for _, tt := range tests {
		t.Run(tt.account, func(t *testing.T) {
			if got := DetectAccountType(tt.account); got != tt.want {
				t.Errorf("DetectAccountType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== 通用验证工具测试 ==========

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fieldName string
		wantErr  bool
	}{
		{"有效值", "test", "字段", false},
		{"空值", "", "字段", true},
		{"只有空格", "   ", "字段", true},
		{"包含空格的有效值", " test ", "字段", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		minLength int
		maxLength int
		wantErr   bool
	}{
		{"有效长度", "test", 3, 10, false},
		{"太短", "ab", 3, 10, true},
		{"太长", "abcdefghijk", 3, 10, true},
		{"边界最小值", "abc", 3, 10, false},
		{"边界最大值", "abcdefghij", 3, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLength(tt.value, tt.minLength, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ========== 数据清理工具测试 ==========

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"正常字符串", "test", "test"},
		{"前后有空格", "  test  ", "test"},
		{"只有空格", "   ", ""},
		{"空字符串", "", ""},
		{"中间有空格", "test value", "test value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeString(tt.value); got != tt.want {
				t.Errorf("SanitizeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"正常邮箱", "Test@Example.com", "test@example.com"},
		{"前后有空格", "  test@example.com  ", "test@example.com"},
		{"大写字母", "USER@DOMAIN.COM", "user@domain.com"},
		{"混合大小写", "Test.User@Domain.Com", "test.user@domain.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeEmail(tt.value); got != tt.want {
				t.Errorf("SanitizeEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizePhone(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"正常手机号", "13812345678", "13812345678"},
		{"带空格", "138 1234 5678", "13812345678"},
		{"带分隔符", "138-1234-5678", "13812345678"},
		{"带括号", "(138)12345678", "13812345678"},
		{"前后有空格", "  13812345678  ", "13812345678"},
		{"包含其他字符", "+86 138 1234 5678", "8613812345678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizePhone(tt.value); got != tt.want {
				t.Errorf("SanitizePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== 批量验证测试 ==========

func TestValidationErrors(t *testing.T) {
	ve := NewValidationErrors()

	// 测试初始状态
	if ve.HasErrors() {
		t.Errorf("ValidationErrors should not have errors initially")
	}

	// 添加错误
	ve.Add("email", "邮箱格式不正确")
	ve.Add("phone", "手机号格式不正确")

	// 检查是否有错误
	if !ve.HasErrors() {
		t.Errorf("ValidationErrors should have errors after adding")
	}

	// 检查错误数量
	if len(ve.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(ve.Errors))
	}

	// 检查错误内容
	if ve.Errors["email"] != "邮箱格式不正确" {
		t.Errorf("Expected '邮箱格式不正确', got '%s'", ve.Errors["email"])
	}

	// 检查错误消息
	errorMsg := ve.Error()
	if errorMsg == "" {
		t.Errorf("Error message should not be empty")
	}
}

// ========== 安全验证测试 ==========

func TestValidateNoSQLInjection(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"正常输入", "username123", false},
		{"包含单引号", "user' or '1'='1", true},
		{"包含双引号", "user\" or \"1\"=\"1", true},
		{"包含分号", "user; drop table", true},
		{"包含注释", "user-- comment", true},
		{"包含SQL关键字", "select * from users", true},
		{"包含大小写混合的关键字", "SELECT * FROM users", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoSQLInjection(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoSQLInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNoXSS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"正常输入", "normal text", false},
		{"包含script标签", "<script>alert('xss')</script>", true},
		{"包含javascript:", "javascript:alert('xss')", true},
		{"包含onclick", "onclick=\"alert('xss')\"", true},
		{"包含iframe", "<iframe src='xss'></iframe>", true},
		{"包含vbscript", "vbscript:msgbox('xss')", true},
		{"包含data URL", "data:text/html,<script>alert('xss')</script>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoXSS(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoXSS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ========== 基准测试 ==========

func BenchmarkValidateEmail(b *testing.B) {
	email := "test@example.com"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateEmail(email)
	}
}

func BenchmarkValidatePhone(b *testing.B) {
	phone := "13812345678"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePhone(phone)
	}
}

func BenchmarkValidateUsername(b *testing.B) {
	username := "testuser123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUsername(username)
	}
}

func BenchmarkValidatePassword(b *testing.B) {
	password := "Password123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePassword(password)
	}
}