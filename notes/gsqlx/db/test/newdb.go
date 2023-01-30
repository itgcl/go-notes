package test

import (
	"go-notes/notes/gsqlx/db"
)

func TestNewDB() db.DataBase {
	db := db.NewDataBase()
	return db
}
