package base

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
)

// 全局
var _db *sqlx.DB
var dbOnce sync.Once

// DB 返回全局指针
func DB() *sqlx.DB {
	dbOnce.Do(func() {
		if _db != nil {
			panic("initDB twice!")
		}

		var err error

		// 连接数据库
		// p, err := pop.Connect("development")

		// must the same as database.yml
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			"postgres", "postgres", "pony_development")
		_db, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Fatalln(err)
		}

		// Test the connection to the database
		err = _db.Ping()
		if err != nil {
			log.Panic("can not ping DB.")
		}
	})

	return _db
}
