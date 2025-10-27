package repository

import (
	"testing"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/entity"
	"github.com/LucienVen/golang-auth-service/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接测试数据库: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

// 创建测试用户
func createTestUser(t *testing.T, db *gorm.DB, username, email, phone, password string) *entity.User {
	user, err := entity.NewUser(username, "测试用户", password, phone, email)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("保存测试用户失败: %v", err)
	}

	return user
}

// ========== 基础CRUD操作测试 ==========

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user, err := entity.NewUser("testuser", "测试用户", "Password123", "13812345678", "test@example.com")
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}

	// 测试创建
	err = repo.Create(user)
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
	}

	// 验证用户已创建
	if user.ID == "" {
		t.Error("用户ID应该被设置")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	testUser := createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试查询
	user, err := repo.GetByID(1) // SQLite自增ID从1开始
	if err != nil {
		t.Errorf("查询用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("查询结果为空")
	}

	if user.Username != "testuser" {
		t.Errorf("用户名不匹配，期望: testuser, 实际: %s", user.Username)
	}

	// 测试不存在的用户
	_, err = repo.GetByID(999)
	if err == nil {
		t.Error("查询不存在的用户应该返回错误")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试查询
	user, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Errorf("查询用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("查询结果为空")
	}

	if user.Username != "testuser" {
		t.Errorf("用户名不匹配，期望: testuser, 实际: %s", user.Username)
	}

	// 测试不存在的邮箱
	_, err = repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("查询不存在的邮箱应该返回错误")
	}
}

func TestUserRepository_GetByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试查询
	user, err := repo.GetByPhone("13812345678")
	if err != nil {
		t.Errorf("查询用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("查询结果为空")
	}

	if user.Username != "testuser" {
		t.Errorf("用户名不匹配，期望: testuser, 实际: %s", user.Username)
	}

	// 测试不存在的手机号
	_, err = repo.GetByPhone("19999999999")
	if err == nil {
		t.Error("查询不存在的手机号应该返回错误")
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试查询
	user, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Errorf("查询用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("查询结果为空")
	}

	if user.Email != "test@example.com" {
		t.Errorf("邮箱不匹配，期望: test@example.com, 实际: %s", user.Email)
	}

	// 测试不存在的用户名
	_, err = repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("查询不存在的用户名应该返回错误")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	testUser := createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 修改用户信息
	testUser.NickName = "新昵称"
	testUser.Status = entity.UserStatusActive

	// 测试更新
	err := repo.Update(testUser)
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
	}

	// 验证更新结果
	user, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Errorf("查询用户失败: %v", err)
	}

	if user.NickName != "新昵称" {
		t.Errorf("昵称更新失败，期望: 新昵称, 实际: %s", user.NickName)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	testUser := createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试软删除
	err := repo.Delete(1)
	if err != nil {
		t.Errorf("删除用户失败: %v", err)
	}

	// 验证用户已被软删除
	_, err = repo.GetByUsername("testuser")
	if err == nil {
		t.Error("软删除的用户应该无法被查询到")
	}
}

// ========== 查询操作测试 ==========

func TestUserRepository_GetByAccount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试邮箱登录
	user, err := repo.GetByAccount("test@example.com")
	if err != nil {
		t.Errorf("邮箱查询用户失败: %v", err)
	}
	if user.Username != "testuser" {
		t.Error("邮箱查询结果不正确")
	}

	// 测试手机号登录
	user, err = repo.GetByAccount("13812345678")
	if err != nil {
		t.Errorf("手机号查询用户失败: %v", err)
	}
	if user.Username != "testuser" {
		t.Error("手机号查询结果不正确")
	}

	// 测试用户名登录
	user, err = repo.GetByAccount("testuser")
	if err != nil {
		t.Errorf("用户名查询用户失败: %v", err)
	}
	if user.Username != "testuser" {
		t.Error("用户名查询结果不正确")
	}

	// 测试不存在的账户
	_, err = repo.GetByAccount("nonexistent")
	if err == nil {
		t.Error("查询不存在的账户应该返回错误")
	}
}

func TestUserRepository_ListByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建多个不同状态的用户
	user1 := createTestUser(t, db, "user1", "user1@example.com", "13812345678", "Password123")
	user2 := createTestUser(t, db, "user2", "user2@example.com", "13812345679", "Password123")

	// 修改用户状态
	user1.Status = entity.UserStatusActive
	user2.Status = entity.UserStatusInactive
	repo.Update(user1)
	repo.Update(user2)

	// 测试查询活跃用户
	users, err := repo.ListByStatus(entity.UserStatusActive, 10, 0)
	if err != nil {
		t.Errorf("查询活跃用户失败: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("活跃用户数量不正确，期望: 1, 实际: %d", len(users))
	}
}

func TestUserRepository_SearchUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test.search@example.com", "13812345678", "Password123")
	createTestUser(t, db, "searchuser", "user.search@example.com", "13812345679", "Password123")

	// 测试搜索
	users, err := repo.SearchUsers("search", 10, 0)
	if err != nil {
		t.Errorf("搜索用户失败: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("搜索结果数量不正确，期望: 2, 实际: %d", len(users))
	}
}

// ========== 业务相关查询测试 ==========

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试存在的邮箱
	exists, err := repo.ExistsByEmail("test@example.com")
	if err != nil {
		t.Errorf("检查邮箱是否存在失败: %v", err)
	}
	if !exists {
		t.Error("存在的邮箱应该返回true")
	}

	// 测试不存在的邮箱
	exists, err = repo.ExistsByEmail("nonexistent@example.com")
	if err != nil {
		t.Errorf("检查邮箱是否存在失败: %v", err)
	}
	if exists {
		t.Error("不存在的邮箱应该返回false")
	}
}

func TestUserRepository_ExistsByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试存在的手机号
	exists, err := repo.ExistsByPhone("13812345678")
	if err != nil {
		t.Errorf("检查手机号是否存在失败: %v", err)
	}
	if !exists {
		t.Error("存在的手机号应该返回true")
	}

	// 测试不存在的手机号
	exists, err = repo.ExistsByPhone("19999999999")
	if err != nil {
		t.Errorf("检查手机号是否存在失败: %v", err)
	}
	if exists {
		t.Error("不存在的手机号应该返回false")
	}
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")

	// 测试存在的用户名
	exists, err := repo.ExistsByUsername("testuser")
	if err != nil {
		t.Errorf("检查用户名是否存在失败: %v", err)
	}
	if !exists {
		t.Error("存在的用户名应该返回true")
	}

	// 测试不存在的用户名
	exists, err = repo.ExistsByUsername("nonexistent")
	if err != nil {
		t.Errorf("检查用户名是否存在失败: %v", err)
	}
	if exists {
		t.Error("不存在的用户名应该返回false")
	}
}

// ========== 统计功能测试 ==========

func TestUserRepository_CountByStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := createTestUser(t, db, "testuser", "test@example.com", "13812345678", "Password123")
	user.Status = entity.UserStatusActive
	repo.Update(user)

	// 测试统计
	count, err := repo.CountByStatus(entity.UserStatusActive)
	if err != nil {
		t.Errorf("统计用户数量失败: %v", err)
	}
	if count != 1 {
		t.Errorf("活跃用户数量不正确，期望: 1, 实际: %d", count)
	}
}

func TestUserRepository_GetStatusDistribution(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建不同状态的用户
	user1 := createTestUser(t, db, "user1", "user1@example.com", "13812345678", "Password123")
	user2 := createTestUser(t, db, "user2", "user2@example.com", "13812345679", "Password123")

	user1.Status = entity.UserStatusActive
	user2.Status = entity.UserStatusInactive
	repo.Update(user1)
	repo.Update(user2)

	// 测试状态分布
	distribution, err := repo.GetStatusDistribution()
	if err != nil {
		t.Errorf("获取状态分布失败: %v", err)
	}

	if distribution[entity.UserStatusActive] != 1 {
		t.Errorf("活跃用户数量不正确，期望: 1, 实际: %d", distribution[entity.UserStatusActive])
	}

	if distribution[entity.UserStatusInactive] != 1 {
		t.Errorf("非活跃用户数量不正确，期望: 1, 实际: %d", distribution[entity.UserStatusInactive])
	}
}

// ========== 事务测试 ==========

func TestUserRepository_WithTransaction(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 开始事务
	err := db.Transaction(func(tx *gorm.DB) error {
		txRepo := repo.WithTransaction(tx)

		// 在事务中创建用户
		user, err := entity.NewUser("txuser", "事务用户", "Password123", "13812345678", "tx@example.com")
		if err != nil {
			return err
		}

		if err := txRepo.Create(user); err != nil {
			return err
		}

		// 在事务中更新用户
		user.Status = entity.UserStatusActive
		if err := txRepo.Update(user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf("事务执行失败: %v", err)
	}

	// 验证事务提交后的结果
	user, err := repo.GetByUsername("txuser")
	if err != nil {
		t.Errorf("查询事务中的用户失败: %v", err)
	}

	if user.Status != entity.UserStatusActive {
		t.Errorf("用户状态不正确，期望: %d, 实际: %d", entity.UserStatusActive, user.Status)
	}
}

// ========== 集成测试 ==========

func TestUserRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 完整的用户生命周期测试
	t.Run("用户完整生命周期", func(t *testing.T) {
		// 1. 创建用户
		user, err := entity.NewUser("lifecycle", "生命周期测试", "Password123", "13812345678", "lifecycle@example.com")
		if err != nil {
			t.Fatalf("创建用户失败: %v", err)
		}

		err = repo.Create(user)
		if err != nil {
			t.Fatalf("保存用户失败: %v", err)
		}

		// 2. 验证用户存在
		exists, err := repo.ExistsByEmail("lifecycle@example.com")
		if err != nil {
			t.Fatalf("检查邮箱存在失败: %v", err)
		}
		if !exists {
			t.Error("用户应该存在")
		}

		// 3. 更新用户状态
		user.Status = entity.UserStatusActive
		user.UpdateLastLogin()
		err = repo.Update(user)
		if err != nil {
			t.Fatalf("更新用户失败: %v", err)
		}

		// 4. 验证更新结果
		updatedUser, err := repo.GetByEmail("lifecycle@example.com")
		if err != nil {
			t.Fatalf("查询用户失败: %v", err)
		}
		if !updatedUser.IsActive() {
			t.Error("用户应该是活跃状态")
		}
		if updatedUser.LastLoginAt == nil {
			t.Error("最后登录时间应该被设置")
		}

		// 5. 软删除用户
		err = repo.Delete(1)
		if err != nil {
			t.Fatalf("删除用户失败: %v", err)
		}

		// 6. 验证用户已被删除
		_, err = repo.GetByEmail("lifecycle@example.com")
		if err == nil {
			t.Error("已删除的用户不应该被查询到")
		}
	})
}

// ========== 性能测试 ==========

func BenchmarkUserRepository_Create(b *testing.B) {
	db := setupTestDB(&testing.T{})
	repo := NewUserRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user, _ := entity.NewUser(
			fmt.Sprintf("user%d", i),
			"测试用户",
			"Password123",
			fmt.Sprintf("1381234%04d", i),
			fmt.Sprintf("user%d@example.com", i),
		)
		repo.Create(user)
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	db := setupTestDB(&testing.T{})
	repo := NewUserRepository(db)

	// 创建测试数据
	for i := 0; i < 100; i++ {
		user, _ := entity.NewUser(
			fmt.Sprintf("user%d", i),
			"测试用户",
			"Password123",
			fmt.Sprintf("1381234%04d", i),
			fmt.Sprintf("user%d@example.com", i),
		)
		repo.Create(user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("user%d@example.com", i%100)
		repo.GetByEmail(email)
	}
}