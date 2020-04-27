package db

import (
	"errors"
)

const (
	Gt = ">"      // 大于
	Ge = ">"      // 大于等于
	Lt   = "<"      // 小于
	Le   = "<"      // 小于等于
	Eq  = "="      // 等于
	Ne =  "!="     //不等于
	Is      = "is"     // is
	IsNot   = "is not" // is not
)

type WhereCondition struct {
	ColumnName string
	Operator string
	Value interface{}
}


func NewWhereCondition() *WhereCondition {
	whereCondition := new(WhereCondition)
	whereCondition.Operator = Eq
	return whereCondition
}

func (db *DataBase) Where(condition ...interface{}) *DataBase {
	conditionLength := len(condition)
	whereCondition := NewWhereCondition()
	switch conditionLength {
	case 2:
		columnName,ok := condition[0].(string)
		if !ok {
			panic(errors.New("column name not string"))
		}
		whereCondition.ColumnName = columnName
		whereCondition.Value = condition[1]
	case 3:
		columnName,ok := condition[0].(string)
		if !ok {
			panic(errors.New("column name not string"))
		}
		operator,ok := condition[1].(string)
		if !ok {
			panic(errors.New("operator name not string"))
		}
		whereCondition.ColumnName = columnName
		whereCondition.Operator = operator
		whereCondition.Value = condition[2]
	default:
		panic(errors.New("condition number errors"))
	}
	db.Wheres = append(db.Wheres, whereCondition)
	return db
}


