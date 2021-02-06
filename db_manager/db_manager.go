package db_manager

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"time"
)

type DBManager struct {
	User string
	Password string
	Uri string
	Database string
	Initialized bool
}

func (manager DBManager) Initialize(user string, psw string, uri string, db string) {
	if user == "" {
		panic("user cannot be empty")
	}

	if db == "" {
		panic("db cannot be empty")
	}

	if uri == "" {
		panic("uri cannot be empty")
	}

	manager.User = user
	manager.Database = db
	manager.Password = psw
	manager.Uri = uri
	manager.Initialized = true
}

func (manager DBManager) connect() *sql.DB {
	if !manager.Initialized {
		panic("DBManager is not yet initialized.")
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", manager.User, manager.Password, manager.Uri, manager.Database))

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return conn
}

func (manager DBManager) InsertRow(row GrassEntry) {
	//conn := manager.connect()
	//defer conn.Close()
}