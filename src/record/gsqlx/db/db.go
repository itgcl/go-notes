package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"reflect"
	"src/record/gsqlx/config"
)

type DataBase struct {
	Core *sqlx.DB
	QB *QB
}

func NewDataBase(params config.DBConnectParams)(dataBase DataBase){
	dbConnectParams := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s",
		params.UserName,
		params.Password,
		params.Host + ":" + config.Get().Port,
		params.DBName,
		params.Charset)
	sqlQB, err := sqlx.Open(config.Get().DriverName, dbConnectParams)
	if err != nil {
		panic(errors.New("mysql connect error"))
	}
	qb := &QB{}
	search := &Search{}
	qb.Search = search
	dataBase.QB = qb
	dataBase.Core = sqlQB
	return
}

func GetPtrElem(ptr interface{}) (prtValue reflect.Value, isPtr bool) {
	valueOf := reflect.ValueOf(ptr)
	if valueOf.Kind() != reflect.Ptr {
		return
	}
	prtValue = valueOf.Elem()
	isPtr = true
	return
}
func (db DataBase) Create(modelPtr Model){
	db.coreCreate(modelPtr)
}

func (db DataBase) coreCreate(modelPtr Model){
	//获取指针值
	valueOf, isPtr := GetPtrElem(modelPtr)
	if isPtr == false {
		panic(errors.New("param not pointer"))
	}
	typeOf := reflect.TypeOf(modelPtr)
	insertData := make(map[string]interface{})
	//遍历结构体参数
	for i:= 0; i < valueOf.NumField(); i++ {
		//获取标签
		tagValue := typeOf.Field(i).Tag.Get("db")
		if tagValue == "" {
			panic(errors.New("tag not exists"))
		}
		insertData[tagValue] = valueOf.Field(i).Interface()
	}
	qb := QB{}
	qb.Insert = insertData
	sqlStr,sqlValueList := qb.BindModel(modelPtr).GetInsert()
	newResult, err := db.Core.Exec(sqlStr, sqlValueList...)
	if err != nil {
		panic(err)
	}
	var result sql.Result
	result = newResult

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	if  lastInsertID != 0 {
		valueOf.FieldByName("ID").SetInt(lastInsertID)
	}
}



func (db DataBase) Update(modelPtr Model){
	db.CoreUpdate(modelPtr)
}

func (db DataBase) CoreUpdate(modelPtr Model){
	valueOf, isPtr := GetPtrElem(modelPtr)
	if isPtr == false {
		panic(errors.New("param not pointer"))
	}
	typeOf := reflect.TypeOf(modelPtr)
	updateData := make(map[string]interface{})
	//判断id是否存在
	var ID interface{}
	if valueOf.FieldByName("ID").Interface() == "" {
		panic(errors.New("update id not exists"))
	}else{
		ID = valueOf.FieldByName("ID").Interface()
	}
	//遍历结构体参数
	for i:= 0; i < valueOf.NumField(); i++ {
		//获取标签
		tagValue := typeOf.Field(i).Tag.Get("db")
		if tagValue == "" {
			panic(errors.New("tag not exists"))
		}
		updateData[tagValue] = valueOf.Field(i).Interface()
	}
	qb := QB{}
	qb.Update = updateData
	sqlStr, sqlValueList := qb.BindModel(modelPtr).GetUpdate()
	//TODO 还没有实现where 这里id没有使用
	_ = ID
	newResult, err := db.Core.Exec(sqlStr, sqlValueList...)
	if err != nil {
		panic(err)
	}
	_ = newResult


}

func (DataBase) Delete(){

}

func (DataBase) OneID( id interface{}){

}

func (db DataBase) Find(modelPtr Model, ID interface{}){

	_, isPtr := GetPtrElem(modelPtr)
	if isPtr == false {
		panic(errors.New("param not pointer"))
	}
	db.Where("id", ID)
	fmt.Println(db.QB.BindModel(modelPtr).WhereSql())
}

func (db DataBase) Select(fieldNameParams ...string) DataBase{
	fieldNameList := []string{}
	for i:= 0; i < len(fieldNameParams) ; i++ {
		fieldNameList = append(fieldNameList, "`" + fieldNameParams[i] + "`")
	}
	db.QB.CoreSelect(fieldNameList)
	return db
}


func (db DataBase) Where(condition ...interface{}) DataBase {
	db.QB.Search.Where(condition...)
	return db
}

func (db DataBase) OrWhere(condition ...interface{}) DataBase {
	 db.QB.Search.OrWhere(condition...)
	return db
}