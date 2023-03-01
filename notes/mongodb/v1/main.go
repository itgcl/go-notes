package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
	M：一张无序的map。它和D是一样的，只是它不保持顺序。
	A：一个BSON数组。
	E：D里面的一个元素。
*/

type Service struct {
	mgClient *mongo.Client
}

func NewService(mgClient *mongo.Client) *Service {
	return &Service{mgClient: mgClient}
}

func main() {
	ctx := context.Background()
	mgClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 检测连接
	err = mgClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	s := NewService(mgClient)

	{
		//if err := s.InsertOne(ctx); err != nil {
		//	log.Fatal(err)
		//}
	}

	{
		//if err := s.BatchInsert(ctx); err != nil {
		//	log.Fatal(err)
		//}
	}

	{
		//record, err := s.FindOne(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("%+v\n", record)
	}

	{
		//recordList, err := s.FindMany(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("%+v\n", recordList)
	}

	{
		//recordList, err := s.FindManyWithIn(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("%+v\n", recordList)
	}
	{
		//err := s.UpdateOne(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	{
		//err := s.UpdateByID(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	{
		//err := s.UpdateMany(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	{
		//err := s.DeleteOne(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	{
		//err := s.DeleteMany(ctx)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	{
		recordList, err := s.FindManyGroup(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(recordList, err)
	}

}

func (s *Service) Collection() *mongo.Collection {
	return s.mgClient.Database("test").Collection("test")
}

// InsertOne 插入单条数据
func (s *Service) InsertOne(ctx context.Context) error {
	result, err := s.Collection().InsertOne(ctx, Record{
		Name:        "record2",
		CreatedTime: time.Now(),
		Count:       10,
		IsTest:      true,
	})
	if err != nil {
		return err
	}
	fmt.Printf("insertOne: %s", result.InsertedID)
	return nil
}

// BatchInsert 批量插入
func (s *Service) BatchInsert(ctx context.Context) error {
	recordList := []interface{}{
		Record{
			Name:        "a1",
			CreatedTime: time.Now(),
			Count:       1,
			IsTest:      false,
		},
		Record{
			Name:        "a2",
			CreatedTime: time.Now(),
			Count:       1,
			IsTest:      true,
		},
		Record{
			Name:        "a3",
			CreatedTime: time.Now(),
			Count:       0,
			IsTest:      false,
		},
	}
	result, err := s.Collection().InsertMany(ctx, recordList)
	if err != nil {
		return err
	}
	fmt.Println(result.InsertedIDs)
	return nil
}

// FindOne 查询单挑数据
func (s *Service) FindOne(ctx context.Context) (interface{}, error) {
	var record Record
	filter := bson.D{{Key: "isTest", Value: true}}
	if err := s.Collection().FindOne(ctx, filter).Decode(&record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) FindMany(ctx context.Context) (interface{}, error) {
	filter := bson.D{{Key: "isTest", Value: true}}
	return s.FindManyBase(ctx, filter)
}

func (s *Service) FindManyWithIn(ctx context.Context) (interface{}, error) {
	filter := bson.D{{
		"name", bson.D{{
			"$in", bson.A{"a1", "a2"},
		}},
	}}
	return s.FindManyBase(ctx, filter)
}

func (s *Service) FindManyGroup(ctx context.Context) (interface{}, error) {
	var recordList []Record
	filter := mongo.Pipeline{
		bson.D{{
			"$group", bson.D{
				{"count", "$count"},
				{"total", bson.D{
					{"$sum", 1},
				}},
			},
		}},
	}
	// 分页
	cursor, err := s.Collection().Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	fmt.Println(results)
	return recordList, err
}

func (s *Service) FindManyBase(ctx context.Context, d bson.D) (interface{}, error) {
	var (
		op         = options.Find()
		recordList []Record
	)
	// 分页
	cursor, err := s.Collection().Find(ctx, d, op.SetSkip(0), op.SetLimit(2))
	if err != nil {
		return nil, err
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	if err := cursor.All(ctx, &recordList); err != nil {
		return nil, err
	}
	return recordList, err
}

func (s *Service) UpdateOne(ctx context.Context) error {
	var (
		filter = bson.D{{
			Key:   "name",
			Value: "a1",
		}}
		update = bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "name",
				Value: "aa11",
			}},
		}}
	)
	result, err := s.Collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("ModifiedCount: %d, MatchedCount: %d, UpsertedCount: %d",
		result.ModifiedCount, result.MatchedCount, result.UpsertedCount)
	return nil
}

func (s *Service) UpdateByID(ctx context.Context) error {
	var (
		update = bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "name",
					Value: "record1223333333",
				},
				{
					Key:   "createdTime",
					Value: time.Now(),
				},
			},
		}}
	)
	id, err := primitive.ObjectIDFromHex("63dfa7d1ea8c7203feacf8cd")
	if err != nil {
		return err
	}
	result, err := s.Collection().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	fmt.Printf("ModifiedCount: %d, MatchedCount: %d, UpsertedCount: %d",
		result.ModifiedCount, result.MatchedCount, result.UpsertedCount)
	return nil
}

func (s *Service) UpdateMany(ctx context.Context) error {
	var (
		filter = bson.D{{
			Key:   "name",
			Value: "qq",
		}}
		update = bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "name",
				Value: "qq1",
			}},
		}}
	)
	result, err := s.Collection().UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("ModifiedCount: %d, MatchedCount: %d, UpsertedCount: %d",
		result.ModifiedCount, result.MatchedCount, result.UpsertedCount)
	return nil
}

func (s *Service) DeleteOne(ctx context.Context) error {
	var (
		filter = bson.D{
			{
				Key:   "name",
				Value: "record2",
			},
			{
				Key:   "isTest",
				Value: true,
			},
		}
	)
	result, err := s.Collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Printf("DeletedCount: %d", result.DeletedCount)
	return nil
}

func (s *Service) DeleteMany(ctx context.Context) error {
	var (
		filter = bson.D{
			{
				Key:   "name",
				Value: "record2",
			},
			{
				Key:   "isTest",
				Value: true,
			},
		}
	)
	result, err := s.Collection().DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Printf("DeletedCount: %d", result.DeletedCount)
	return nil
}
