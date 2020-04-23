package db

type Model interface {
	TableName() string
	BeforeCreate()
}
