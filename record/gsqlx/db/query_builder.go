package db

import (
	"errors"
	"reflect"
	"strings"
)

type QB struct {
	Table  string
	Insert map[string]interface{}
	Update map[string]interface{}
	Search *Search
	Select string
}

func (qb QB) GetInsert() (sql string, sqlValues []interface{}) {
	sql = "INSERT INTO `" + qb.Table + "` ("
	fieldList := []string{}
	replaceValueList := []string{}
	for key, value := range qb.Insert {
		fieldList = append(fieldList, "`"+key+"`")
		replaceValueList = append(replaceValueList, "?")
		sqlValues = append(sqlValues, value)
	}
	fieldStr := strings.Join(fieldList, ",")
	fieldStr = strings.TrimRight(fieldStr, ",")
	fieldStr += ") VALUES ("
	replaceValueStr := strings.Join(replaceValueList, ",")
	replaceValueStr = strings.TrimRight(replaceValueStr, ",")
	replaceValueStr += " )"
	sql += fieldStr + replaceValueStr
	return
}

func (qb QB) GetUpdate() (sql string, sqlValues []interface{}) {
	sql = "Update `" + qb.Table + "` SET "
	replaceKeyList := []string{}
	for key, value := range qb.Insert {
		replaceKeyList = append(replaceKeyList, "`"+key+"` = ?")
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
		sql = ``
		return
	}
	for _, where := range wheres {
		columnNameList = append(columnNameList, where.ColumnName+where.Operator+`?`)
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
		columnNameList = append(columnNameList, where.ColumnName+where.Operator+`?`)
		sqlValues = append(sqlValues, where.Value)
	}
	sql += strings.Join(columnNameList, ` OR `)
	return
}

func (qb QB) WhereSql() (sql string, sqlValues []interface{}) {
	if len(qb.Search.wheres) == 0 && len(qb.Search.orWheres) == 0 {
		return
	}
	whereList := []string{}
	orWhereList := []string{}
	sql = ` WHERE `
	if len(qb.Search.wheres) != 0 {
		for _, where := range qb.Search.wheres {
			whereList = append(whereList, where.ColumnName+where.Operator+`?`)
			sqlValues = append(sqlValues, where.Value)
		}
	}
	if len(qb.Search.orWheres) != 0 {
		for _, orWhere := range qb.Search.orWheres {
			orWhereList = append(orWhereList, orWhere.ColumnName+orWhere.Operator+`?`)
			sqlValues = append(sqlValues, orWhere.Value)
		}
	}
	sql += strings.Join(whereList, ` AND `)
	if len(orWhereList) != 0 && len(whereList) != 0 {
		sql += ` OR `
	}
	sql += strings.Join(orWhereList, ` OR `)
	return
}

func (qb QB) CoreFind(ID interface{}) (sql string, sqlValues []interface{}) {
	//select * from table where aa = 'aa' and bb = 'bb'
	sql = `SELECT ` + qb.Select + ` from ` + qb.Table
	qb.Search.Where("id", ID)
	//sqlWhere, values := qb.WhereSql(qb.Search.wheres)

	return
}

func (qb *QB) CoreSelect(fieldList []string) string {
	if len(fieldList) == 0 {
		qb.Select = `*`
		return qb.Select
	}
	qb.Select = strings.Join(fieldList, `,`)
	return qb.Select
}
