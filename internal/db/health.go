package db

import (
	"context"
	"time"

	"github.com/LucienVen/golang-auth-service/pkg/log"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	db     *GormDB
	ticker *time.Ticker
	done   chan struct{}
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(db *GormDB) *HealthChecker {
	return &HealthChecker{
		db:   db,
		done: make(chan struct{}),
	}
}

// Start 启动健康检查
func (h *HealthChecker) Start(ctx context.Context, interval time.Duration) {
	h.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-h.ticker.C:
				if err := h.db.Ping(); err != nil {
					log.Errorf("数据库心跳检测失败: %v", err)
					// 尝试重连
					if err := h.db.Connect(); err != nil {
						log.Errorf("数据库重连失败: %v", err)
					} else {
						log.Info("数据库重连成功")
					}
				} else {
					log.Info("数据库心跳检测正常")
				}
			case <-ctx.Done():
				h.Stop()
				return
			case <-h.done:
				return
			}
		}
	}()
}

// Stop 停止健康检查
func (h *HealthChecker) Stop() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
	close(h.done)
}

// Shutdown 实现 Shutdownable 接口
func (h *HealthChecker) Shutdown(ctx context.Context) error {
	h.Stop()
	return nil
}
