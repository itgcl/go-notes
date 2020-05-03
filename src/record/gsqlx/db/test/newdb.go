package test

import (
	"src/record/gsqlx/config"
	"src/record/gsqlx/db"
)

func TestNewDB() db.DataBase{
	db := db.NewDataBase(config.DBConnectParams{
		DriverName: "mysql",
		UserName:   "root",
		Password:   "root",
		Host:       "127.0.0.1",
		Port:       "3306",
		DBName:     "test",
		Charset:    "utf8",
	})
	return db
}
