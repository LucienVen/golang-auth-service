package app

import (
	"context"
	"fmt"
	"time"
)

// Shutdownable 可关闭接口
type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

// ShutdownManager 关闭管理器
type ShutdownManager struct {
	components []Shutdownable
}

// NewShutdownManager 创建关闭管理器
func NewShutdownManager() *ShutdownManager {
	return &ShutdownManager{
		components: make([]Shutdownable, 0),
	}
}

// Register 注册可关闭组件
func (m *ShutdownManager) Register(component Shutdownable) {
	m.components = append(m.components, component)
}

// Shutdown 关闭所有组件
func (m *ShutdownManager) Shutdown(ctx context.Context) error {
	// 创建超时上下文
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 按顺序关闭组件
	for _, component := range m.components {
		if err := component.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("组件关闭失败: %w", err)
		}
	}

	return nil
}
