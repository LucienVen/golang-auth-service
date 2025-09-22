package appcontext

import "github.com/LucienVen/golang-auth-service/internal/db"

// AppContext 聚合全局依赖
type AppContext struct {
	DB    db.DB
	Redis *db.RedisClient
}
