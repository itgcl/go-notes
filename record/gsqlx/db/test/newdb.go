package test

import (
	"go-notes/record/gsqlx/db"
)

func TestNewDB() db.DataBase {
	db := db.NewDataBase()
	return db
}
