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


func TestWhereSql(t *testing.T)  {
	db := TestNewDB()
	db.Where("id", 1)
	db.Where("name", "张三").Where("password","123456")
	sql, values := db.QB.WhereSql()
	assert.Equal(t, " WHERE id=? AND name=? AND password=?", sql)
	assert.Equal(t, []interface {}{1, "张三", "123456"}, values)
}


func TestOrWhereSql(t *testing.T)  {
	db := TestNewDB()
	db.OrWhere("name", "张三").OrWhere("password","123456")
	sql, values := db.QB.WhereSql()
	fmt.Println(sql)
	assert.Equal(t, " WHERE name=? OR password=?", sql)
	assert.Equal(t, []interface {}{"张三", "123456"}, values)
}

func TestWhereAndOrWhereSql(t *testing.T)  {
	db := TestNewDB()
	db.Where("id", 10).
		OrWhere("name", "张三").
		Where("address", "上海").
		OrWhere("password","123456")
	sql, values := db.QB.WhereSql()
	fmt.Println(sql)
	assert.Equal(t, " WHERE id=? AND address=? OR name=? OR password=?", sql)
	assert.Equal(t, []interface {}{10, "上海", "张三", "123456"}, values)
}


func TestSelect(t *testing.T)  {
	db := TestNewDB()
	db.Select("id", "name", "created_at")
	assert.Equal(t, db.QB.Select, "`id`,`name`,`created_at`")
}