package db

import (
	"errors"
	"reflect"
	"strings"
)

type QB struct {
	Table	string
	Select  []string
	Insert map[string]interface{}
	Update map[string]interface{}
}


func (qb QB) GetInsert() (sql string, sqlValues []interface{}) {
	sql = "INSERT INTO `" + qb.Table + "` ("
	fieldList := []string{}
	replaceValueList := []string{}
	for key,value := range qb.Insert{
		fieldList = append(fieldList, "`"+ key + "`")
		replaceValueList = append(replaceValueList, "?")
		sqlValues = append(sqlValues, value)
	}
	fieldStr := strings.Join(fieldList, ",")
	fieldStr = strings.TrimRight(fieldStr, ",")
	fieldStr += ") VALUES ("
	replaceValueStr := strings.Join(replaceValueList, ",")
	replaceValueStr = strings.TrimRight(replaceValueStr, ",")
	replaceValueStr = " )"
	sql += fieldStr + replaceValueStr
	return
}

func (qb QB) GetUpdate() (sql string, sqlValues []interface{}) {
	sql = "Update `" + qb.Table + "` SET "
	replaceKeyList := []string{}
	for key, value := range qb.Insert{
		replaceKeyList = append(replaceKeyList, "`"+ key + "` = ?")
		sqlValues = append(sqlValues, value)
	}
	replaceKeySqlStr := strings.Join(replaceKeyList, ",")
	replaceKeySqlStr = strings.TrimRight(replaceKeySqlStr, ",")
	sql += replaceKeySqlStr
	return
}

func (qb QB) BindModel(modelPtr Model) QB {
	if qb.Table != "" {
		return qb
	}
	tableName := reflect.ValueOf(modelPtr).MethodByName("TableName").Call(nil)[0].String()
	if tableName == "" {
		panic(errors.New("tableName not exists"))
	}
	qb.Table = tableName
	return qb
}

func (qb QB) CoreWhere(wheres []*WhereCondition) (sql string, sqlValues []interface{}) {
	columnNameList := []string{}
	sql = ` where `
	if len(wheres) == 0 {
		return
	}
	for _, where := range wheres {
		columnNameList = append(columnNameList,  where.ColumnName + where.Operator + `?`)
		sqlValues = append(sqlValues, where.Value)
	}
	sql += strings.Join(columnNameList, ` AND `)
	return
}

func (qb QB) CoreOrWhere(wheres []*OrWhereCondition) (sql string, sqlValues []interface{}) {
	columnNameList := []string{}
	sql = ` where `
	if len(wheres) == 0 {
		return
	}
	for _, where := range wheres {
		columnNameList = append(columnNameList,  where.ColumnName + where.Operator + `?`)
		sqlValues = append(sqlValues, where.Value)
	}
	sql += strings.Join(columnNameList, ` OR `)
	return
}

func (qb QB) CoreFind(ID string) (sql string, sqlValues[]interface{})  {
	sql = `SELECT `
	return
}

func (qb *QB) CoreSelect(fieldList []string) *QB {
	if len(fieldList) == 0 {
		fieldList = append(fieldList, `*`)
	}
	qb.Select = fieldList
	return qb
}
