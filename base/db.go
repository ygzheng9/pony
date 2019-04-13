package base

import (
	"database/sql"
	"log"
	"sync"

	"github.com/gobuffalo/pop"
	"github.com/jmoiron/sqlx"
)

// 全局
var _conn *pop.Connection
var dbOnce sync.Once

// Store copy from pop
type Store interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	NamedExec(string, interface{}) (sql.Result, error)
	Exec(string, ...interface{}) (sql.Result, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Transaction() (*pop.Tx, error)
	Rollback() error
	Commit() error
	Close() error
}

// DB 返回 pop.connection 中的 Store
func DB() Store {
	if _conn == nil {
		Conn()
	}

	return _conn.Store
}

// Conn 数据库连接，这里使用了 pop，不是直接的 *sqlx.DB
func Conn() *pop.Connection {
	dbOnce.Do(func() {
		if _conn != nil {
			panic("initDB twice!")
		}

		var err error

		// 连接数据库
		// must the same as database.yml
		// connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		// 	"postgres", "postgres", "pony_development")
		// _db, err = sqlx.Connect("postgres", connStr)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// Test the connection to the database
		// err = _db.Ping()
		// if err != nil {
		// 	log.Panic("can not ping DB.")
		// }

		// env := envy.Get("GO_ENV", "development")
		_conn, err = pop.Connect("development")
		if err != nil {
			log.Fatal(err)
		}
	})

	return _conn
}
