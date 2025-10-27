package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/LucienVen/golang-auth-service/internal/entity"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户Repository实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户Repository实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// ========== 基础CRUD操作 ==========

// Create 创建用户记录
func (r *UserRepositoryImpl) Create(user *entity.User) error {
	if user == nil {
		return errors.New("用户对象不能为空")
	}

	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// GetByID 根据用户ID获取用户信息
func (r *UserRepositoryImpl) GetByID(id uint) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("id = ? AND is_delete = 0", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// GetByEmail 根据邮箱获取用户信息
func (r *UserRepositoryImpl) GetByEmail(email string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("email = ? AND is_delete = 0", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// GetByPhone 根据手机号获取用户信息
func (r *UserRepositoryImpl) GetByPhone(phone string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("phone = ? AND is_delete = 0", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// GetByUsername 根据用户名获取用户信息
func (r *UserRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where("username = ? AND is_delete = 0", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(user *entity.User) error {
	if user == nil {
		return errors.New("用户对象不能为空")
	}

	if user.ID == "" {
		return errors.New("用户ID不能为空")
	}

	result := r.db.Model(user).Where("id = ? AND is_delete = 0", user.ID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已被删除")
	}

	return nil
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(id uint) error {
	result := r.db.Model(&entity.User{}).Where("id = ? AND is_delete = 0", id).Update("is_delete", 1)
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在或已被删除")
	}

	return nil
}

// HardDelete 硬删除用户（物理删除）
func (r *UserRepositoryImpl) HardDelete(id uint) error {
	result := r.db.Delete(&entity.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("硬删除用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// ========== 查询操作 ==========

// GetByAccount 统一账户查询接口
func (r *UserRepositoryImpl) GetByAccount(account string) (*entity.User, error) {
	if account == "" {
		return nil, errors.New("账户不能为空")
	}

	var user entity.User

	// 使用OR条件查询邮箱、手机号或用户名
	query := r.db.Where("(email = ? OR phone = ? OR username = ?) AND is_delete = 0", account, account, account)

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// ListByStatus 根据用户状态批量查询用户
func (r *UserRepositoryImpl) ListByStatus(status int8, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User

	query := r.db.Where("status = ? AND is_delete = 0", status)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("create_time DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}

	return users, nil
}

// ListActive 查询所有活跃用户
func (r *UserRepositoryImpl) ListActive(limit, offset int) ([]*entity.User, error) {
	return r.ListByStatus(entity.UserStatusActive, limit, offset)
}

// CountByStatus 统计指定状态的用户数量
func (r *UserRepositoryImpl) CountByStatus(status int8) (int64, error) {
	var count int64

	if err := r.db.Model(&entity.User{}).Where("status = ? AND is_delete = 0", status).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计用户数量失败: %w", err)
	}

	return count, nil
}

// CountAll 统计所有用户数量
func (r *UserRepositoryImpl) CountAll() (int64, error) {
	var count int64

	if err := r.db.Model(&entity.User{}).Where("is_delete = 0").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计用户总数失败: %w", err)
	}

	return count, nil
}

// SearchUsers 搜索用户（支持用户名、昵称、邮箱模糊搜索）
func (r *UserRepositoryImpl) SearchUsers(keyword string, limit, offset int) ([]*entity.User, error) {
	if keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}

	var users []*entity.User

	// 构建模糊搜索条件
	searchPattern := "%" + keyword + "%"
	query := r.db.Where("(username LIKE ? OR nick_name LIKE ? OR email LIKE ?) AND is_delete = 0",
		searchPattern, searchPattern, searchPattern)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("create_time DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	return users, nil
}

// ========== 业务相关查询 ==========

// ExistsByEmail 检查邮箱是否已存在
func (r *UserRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	if email == "" {
		return false, nil
	}

	var count int64
	if err := r.db.Model(&entity.User{}).Where("email = ? AND is_delete = 0", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ExistsByPhone 检查手机号是否已存在
func (r *UserRepositoryImpl) ExistsByPhone(phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}

	var count int64
	if err := r.db.Model(&entity.User{}).Where("phone = ? AND is_delete = 0", phone).Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查手机号是否存在失败: %w", err)
	}

	return count > 0, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepositoryImpl) ExistsByUsername(username string) (bool, error) {
	if username == "" {
		return false, nil
	}

	var count int64
	if err := r.db.Model(&entity.User{}).Where("username = ? AND is_delete = 0", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查用户名是否存在失败: %w", err)
	}

	return count > 0, nil
}

// GetUsersByLastLoginTime 根据最后登录时间查询用户
func (r *UserRepositoryImpl) GetUsersByLastLoginTime(startTime, endTime int64, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User

	query := r.db.Where("last_login_at >= ? AND last_login_at <= ? AND is_delete = 0", startTime, endTime)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("last_login_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return users, nil
}

// GetInactiveUsers 获取长期未登录的用户
func (r *UserRepositoryImpl) GetInactiveUsers(days int, limit, offset int) ([]*entity.User, error) {
	if days <= 0 {
		return nil, errors.New("天数必须大于0")
	}

	var users []*entity.User

	// 计算截止时间（days天前）
	cutoffTime := time.Now().AddDate(0, 0, -days).Unix()

	query := r.db.Where("(last_login_at < ? OR last_login_at IS NULL) AND is_delete = 0", cutoffTime)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Order("last_login_at ASC NULLS FIRST").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询不活跃用户失败: %w", err)
	}

	return users, nil
}

// ========== 批量操作 ==========

// CreateBatch 批量创建用户
func (r *UserRepositoryImpl) CreateBatch(users []*entity.User) error {
	if len(users) == 0 {
		return nil
	}

	if err := r.db.CreateInBatches(users, 100).Error; err != nil {
		return fmt.Errorf("批量创建用户失败: %w", err)
	}

	return nil
}

// UpdateBatch 批量更新用户状态
func (r *UserRepositoryImpl) UpdateBatch(userIDs []uint, status int8) error {
	if len(userIDs) == 0 {
		return nil
	}

	result := r.db.Model(&entity.User{}).Where("id IN ? AND is_delete = 0", userIDs).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("批量更新用户状态失败: %w", result.Error)
	}

	return nil
}

// DeleteBatch 批量软删除用户
func (r *UserRepositoryImpl) DeleteBatch(userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	result := r.db.Model(&entity.User{}).Where("id IN ? AND is_delete = 0", userIDs).Update("is_delete", 1)
	if result.Error != nil {
		return fmt.Errorf("批量删除用户失败: %w", result.Error)
	}

	return nil
}

// ========== 事务相关 ==========

// WithTransaction 开始事务
func (r *UserRepositoryImpl) WithTransaction(tx *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: tx,
	}
}

// ========== 统计和分析 ==========

// GetRegistrationStats 获取用户注册统计
func (r *UserRepositoryImpl) GetRegistrationStats(days int) (map[string]int64, error) {
	if days <= 0 {
		return nil, errors.New("天数必须大于0")
	}

	// 计算开始时间
	startTime := time.Now().AddDate(0, 0, -days).Unix()

	var results []struct {
		Date  string
		Count int64
	}

	// 按日期分组统计注册用户数
	query := `
		SELECT DATE(FROM_UNIXTIME(create_time)) as date, COUNT(*) as count
		FROM users
		WHERE create_time >= ? AND is_delete = 0
		GROUP BY DATE(FROM_UNIXTIME(create_time))
		ORDER BY date DESC
	`

	if err := r.db.Raw(query, startTime).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("获取注册统计失败: %w", err)
	}

	// 转换为map格式
	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Date] = result.Count
	}

	return stats, nil
}

// GetStatusDistribution 获取用户状态分布
func (r *UserRepositoryImpl) GetStatusDistribution() (map[int8]int64, error) {
	var results []struct {
		Status int8
		Count  int64
	}

	if err := r.db.Model(&entity.User{}).
		Select("status, COUNT(*) as count").
		Where("is_delete = 0").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("获取状态分布失败: %w", err)
	}

	// 转换为map格式
	distribution := make(map[int8]int64)
	for _, result := range results {
		distribution[result.Status] = result.Count
	}

	return distribution, nil
}