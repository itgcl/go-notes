package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Territory string

func (t *Territory) UnmarshalJSON(b []byte) error {
	fmt.Println(222)
	return nil
}

type Data struct {
	Territory Territory `json:"territory"`
	ID        int64     `json:"id"`
}

func UnMarshal(data []byte) (*Message, error) {
	var entry entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &Message{entry: &entry}, nil
}

/*
	switch msg.Type {
		case canal.UPDATE:
			...
		case canal.INSERT:
			...
		case canal.DELETE:
			...
		default:
			...
	}
*/
const (
	UPDATE = "UPDATE"
	INSERT = "INSERT"
	DELETE = "DELETE"
)

// Box 不知道起什么名字
type Box interface {
	Get(key string) string
	GetInt(key string) (int64, error)
	GetBool(key string) (bool, error)
	Has(key string) bool
	set(key, val string)
	Bind(obj interface{}) error
}

type box map[string]string

var _ Box = box(nil)

func (r box) Has(key string) bool {
	_, ok := r[key]
	return ok
}

func (r box) set(key, val string) {
	r[key] = val
}

func (r box) Get(key string) string {
	return r[key]
}
func (r box) GetInt(key string) (int64, error) {
	x := r.Get(key)
	return strconv.ParseInt(x, 10, 64)
}
func (r box) GetBool(key string) (bool, error) {
	x := r.Get(key)
	return strconv.ParseBool(x)
}

func (r box) Bind(obj interface{}) error {
	return fillMapToStruct(r, reflect.ValueOf(obj))
}

type entry struct {
	Data      []box    `json:"data,omitempty"`
	Database  string   `json:"database,omitempty"`
	ES        int64    `json:"es,omitempty"` // binlog executeTime，单位：ms
	ID        int64    `json:"id,omitempty"` // batch_id
	IsDDL     bool     `json:"isDdl,omitempty"`
	MysqlType box      `json:"mysqlType,omitempty"`
	Old       []box    `json:"old,omitempty"`
	PkNames   []string `json:"pkNames,omitempty"`
	SQL       string   `json:"sql,omitempty"`
	Table     string   `json:"table,omitempty"`
	TS        int64    `json:"ts,omitempty"` // 操作执行当前时间，单位：ms
	Type      string   `json:"type,omitempty"`
}

// Message 是 canal 发送到 kafka / mq 等外部队列后的消息格式在 Go 中的映射。
// 如果是直连 canal 的话，这个类型不适用
// 内部包装了一层 entry 是为了使数据不可写，外部调用 Message 只允许读。
type Message struct {
	entry *entry
}

// DMLType DML 操作类型，包括 update、insert、delete
func (msg *Message) DMLType() string {
	return msg.entry.Type
}

// IsDDL 是否为 DDL 操作
func (msg *Message) IsDDL() bool {
	return msg.entry.IsDDL
}

// Table 表名
func (msg *Message) Table() string {
	return msg.entry.Table
}

// Database 数据库名
func (msg *Message) Database() string {
	return msg.entry.Database
}

// RowChanged 受影响的行树
func (msg *Message) RowsCount() int {
	return len(msg.entry.Data)
}

// PKNames 主键列表
func (msg *Message) PKNames() []string {
	return msg.entry.PkNames
}

func (msg *Message) Rows() []RowChange {
	var result []RowChange
	for i, row := range msg.entry.Data {
		rowchange := &rowChange{row: row}
		if len(msg.entry.Old) > i {
			rowchange.before = msg.entry.Old[i]
			rowchange.after = make(box, len(rowchange.before))
			for k := range rowchange.before {
				rowchange.after.set(k, row.Get(k))
			}
		}
		rowchange.pk = make(box, len(msg.entry.PkNames))
		for _, filed := range msg.entry.PkNames {
			rowchange.pk.set(filed, row.Get(filed))
		}
		result = append(result, rowchange)
	}
	return result
}

type RowChange interface {
	Get(key string) string
	// Row 获取变更后的全部数据
	Row() Box
	// After 仅发生变更的字段在更新之前的值
	After() Box
	// Before 仅发生变更的字段在更新之前的值
	Before() Box
	PK() Box
}

type rowChange struct {
	row    box
	after  box
	before box
	pk     box
}

func (rc *rowChange) Get(key string) string {
	return rc.row.Get(key)
}

func (rc *rowChange) Row() Box {
	return rc.row
}

func (rc *rowChange) After() Box {
	return rc.after
}

func (rc *rowChange) Before() Box {
	return rc.before
}

func (rc *rowChange) PK() Box {
	return rc.pk
}

func fillMapToStruct(src map[string]string, val reflect.Value) error {
	for val.Kind() == reflect.Pointer {
		if val.IsNil() && val.CanSet() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("struct only, not %v", val.Kind())
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldName := field.Name
		fieldName = string(unicode.ToLower(rune(fieldName[0]))) + fieldName[1:]
		if tag, ok := field.Tag.Lookup("json"); ok {
			// 跳过 - 字段
			if tag == "-" {
				continue
			}
			fieldName, _, _ = strings.Cut(tag, ",")
		}
		v, ok := src[fieldName]
		if !ok {
			continue
		}
		fv := val.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		if err := setField(fv, v); err != nil {
			return err
		}
	}
	return nil
}

func setField(val reflect.Value, v string) error {
	t := val.Interface()
	if _, ok := t.(time.Time); ok {
		vv, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(vv))
		return nil
	}
	if _, ok := t.(*time.Time); ok {
		vv, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(&vv))
		return nil
	}
	switch val.Kind() {
	case reflect.Bool:
		d, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		val.SetBool(d)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(d)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(d)
	case reflect.Float32, reflect.Float64:
		d, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		val.SetFloat(d)
	case reflect.String:
		val.SetString(v)
	case reflect.Slice:
		rec := reflect.MakeSlice(reflect.SliceOf(val.Type().Elem()), 10, 10) //nolint:gomnd
		val.Set(rec)
		if err := json.Unmarshal([]byte(v), val.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Struct:
		rec := reflect.New(val.Type())
		if err := json.Unmarshal([]byte(v), rec.Interface()); err != nil {
			return err
		}
		val.Set(rec.Elem())
	case reflect.Pointer:
		for val.Kind() == reflect.Pointer {
			if val.IsNil() && val.CanSet() {
				val.Set(reflect.New(val.Type().Elem()))
			}
			val = val.Elem()
		}
		return setField(val, v)
	default:
		return errors.New("unsupported type")
	}

	return nil
}

func camelToSnake(s string) string {
	var result string
	for i, v := range s {
		if v >= 'A' && v <= 'Z' {
			if i > 0 {
				result += "_"
			}
			result += string(v + 32) //nolint:gomnd
		} else {
			result += string(v)
		}
	}
	return result
}
