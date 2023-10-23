package main

import (
	"context"
	"database/sql"
	"log"

	// 注册 mysql 数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	UserID int64
}

func main() {

	// 创建 db 实例
	db, err := sql.Open("mysql", "username:passpord@(ip:port)/database")
	if err != nil {
		log.Fatal(err)
		return
	}

	// 执行 sql
	ctx := context.Background()
	row := db.QueryRowContext(ctx, "SELECT user_id FROM user WHERE ORDER BY created_at DESC limit 1")
	if row.Err() != nil {
		log.Print(err)
		return
	}

	// 解析结果
	var u user
	if err = row.Scan(&u.UserID); err != nil {
		log.Print(err)
		return
	}
	log.Print(u.UserID)
}
