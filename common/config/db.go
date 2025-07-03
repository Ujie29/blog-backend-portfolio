package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"         // PostgreSQL 驅動程式（pgx 版本）
	"github.com/uptrace/bun"                   // Bun ORM 主套件
	"github.com/uptrace/bun/dialect/pgdialect" // PostgreSQL 語法支援
)

// Database struct 包裹一個 *bun.DB 實體，用來執行所有 ORM 操作
type Database struct {
	DB *bun.DB
}

// InitDB 用來初始化資料庫連線，回傳一個 Database 結構體
func InitDB() *Database {
	cfg := LoadDBConfig() // 載入自定義資料庫設定（從 YAML 或 .env）

	// 組合 PostgreSQL 的 DSN（資料庫連線字串）
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=Asia/Taipei",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	// 使用 pgx 套件開啟連線，得到 *sql.DB 物件
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// 使用 Bun 包裝 sql.DB，指定使用 PostgreSQL 語法
	db := bun.NewDB(sqldb, pgdialect.New())

	// 回傳自訂的 Database 結構體，讓其他地方可以使用 db 連線
	return &Database{
		DB: db,
	}
}
