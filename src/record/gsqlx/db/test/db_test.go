package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"src/record/gsqlx/config"
	"src/record/gsqlx/db"
	"testing"
)

func TestDBConnect(t *testing.T){
	{
		db := db.NewDataBase(config.DBConnectParams{
			DriverName: "mysql",
			UserName:   "root",
			Password:   "root",
			Host:       "127.0.0.1",
			Port:       "3306",
			DBName:     "test",
			Charset:    "utf8",
		})
		assert.Equal(t, db.Core.DriverName(), "mysql")
	}
}

func TestGetPtrElem(t *testing.T){
	type TestStudent struct {
		Name string
		Age int
		Address string
	}
	{
		testStudent := TestStudent{
			Name:    "张三",
			Age:     10,
			Address: "上海",
		}
		_, isPtr := db.GetPtrElem(testStudent)
		assert.Equal(t, isPtr, false)
	}
	{
		testStudent := &TestStudent{
			Name:    "李四",
			Age:     20,
			Address: "上海",
		}
		ptrValue, isPtr := db.GetPtrElem(testStudent)
		assert.Equal(t, isPtr, true)
		ptrValue.FieldByName("Name").SetString("王五")
		ptrValue.FieldByName("Age").SetInt(30)
		assert.Equal(t, testStudent.Name, "王五")
		assert.Equal(t, testStudent.Age, 30)
	}
}

func TestCreate(t *testing.T){
	//db := db.NewDataBase(config.DBConnectParams{
	//	DriverName: "mysql",
	//	UserName:   "root",
	//	Password:   "root",
	//	Host:       "127.0.0.1",
	//	Port:       "3306",
	//	DBName:     "test",
	//	Charset:    "utf8",
	//})
	//db.Core.MustExec(schema)
	//db.Create()
}

func TestWhere(t *testing.T){
	dba := db.NewDataBase(config.DBConnectParams{
		DriverName: "mysql",
		UserName:   "root",
		Password:   "root",
		Host:       "127.0.0.1",
		Port:       "3306",
		DBName:     "test",
		Charset:    "utf8",
	})
	fmt.Println(dba.Where("aSD", "qWEQ") . Where("qwe", "Wt").first())
	//db.Core.MustExec(schema)
	//db.Create()
}