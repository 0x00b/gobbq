package mongo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/engine/db/mongo"
	"github.com/0x00b/gobbq/example/exampb"
)

// func TestMain(t *testing.T) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	collection := client.Database("testing").Collection("numbers")

// 	res, err := collection.InsertOne(ctx, bson.D{bson.E{Key: "name", Value: "pi"}, bson.E{Key: "value", Value: 3.14159}})
// 	if err != nil {
// 		panic(err)
// 	}
// 	id := res.InsertedID
// 	fmt.Println("insert id", id)

// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cur.Close(ctx)
// 	for cur.Next(ctx) {
// 		var result bson.D
// 		err := cur.Decode(&result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// do something with result....
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func TestFieldMarsh(t *testing.T) {
// 	mgo := mongo.NewMongoDB()

// 	a := &exampb.EchoProperty{
// 		Test: &exampb.SayHelloRequest{
// 			Text: "xxx",
// 		},
// 		Text: "bbbb",
// 		Test3: map[int32]string{
// 			1: "cxcc",
// 		},
// 		TEST7: 1,
// 	}
// 	m, err := mgo.PartialMarshalToMap(a, []model.FieldName{
// 		exampb.EchoProperty_Text,
// 		exampb.EchoProperty_Test,
// 		exampb.EchoProperty_Test2,
// 		exampb.EchoProperty_TEST7})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(m)

// 	fmt.Println(mongo.GetMongoID(a))

// }

// func TestDB(t *testing.T) {
// 	mgo := mongo.NewMongoDB()
// 	err := mgo.Connect(&mongo.Config{
// 		URL:            "mongodb://127.0.0.1:27017",
// 		DBName:         "testing",
// 		CollectionName: "bbq",
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	a := &exampb.EchoProperty{
// 		Test: &exampb.SayHelloRequest{
// 			Text: "xxx",
// 		},
// 		Text: "bbbb",
// 		Test3: map[int32]string{
// 			1: "cxcc",
// 		},
// 		TEST7: 1,
// 	}

// 	err = mgo.Insert(context.Background(), a)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	a.TEST7 = 0
// 	a.Test5 = 111

// 	err = mgo.Update(context.Background(), a, []model.FieldName{exampb.EchoProperty_Test5, exampb.EchoProperty_TEST7})
// 	if err != nil {
// 		panic(err)
// 	}

// 	mgo.Load(context.Background(), a)

// 	fmt.Println("load:", a)
// }

func TestDB2(t *testing.T) {
	mgo := mongo.NewMongoDB()
	err := mgo.Connect(&mongo.Config{
		URL:            "mongodb://127.0.0.1:27017",
		DBName:         "testing",
		CollectionName: "bbq",
	})
	if err != nil {
		panic(err)
	}

	a := &exampb.EchoPropertyModel{}

	c := context.Background()

	err = a.ModelInit(c, mgo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)

	a.Text = "bbbb"

	err = a.ModelSave(c)
	if err != nil {
		panic(err)
	}

	a.SetTest5(5555)

	// save test5
	a.ModelAutoSave(c)

	a.SetTest6(123)

	// save test6
	a.ModelStopWatcher(c)

}

// func TestMarsh(t *testing.T) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	collection := client.Database("testing").Collection("bbq")

// 	a := &exampb.EchoProperty{
// 		Test: &exampb.SayHelloRequest{
// 			Text: "xxx",
// 		},
// 		Text: "bbbb",
// 		Test3: map[int32]string{
// 			1: "cxcc",
// 		},
// 	}

// 	res, err := collection.InsertOne(ctx, a)
// 	// res, err := collection.InsertOne(ctx, bson.M{"A": "111"})
// 	if err != nil {
// 		panic(err)
// 	}
// 	id := res.InsertedID
// 	fmt.Println("insert id", id)

// 	// _, err = collection.UpdateByID(ctx, id, bson.D{{"$set", bson.D{{Key: "A", Value: "222"}}}})
// 	// _, err = collection.ReplaceOne(ctx, bson.D{{Key: "A", Value: "111"}}, bson.D{{Key: "a", Value: "222"}})
// 	// _, err = collection.UpdateByID(ctx, id, bson.D{{"$set", bson.D{{Key: "a", Value: "222"}}}})
// 	// if err != nil {
// 	// fmt.Println("err", err)
// 	// }

// 	cur, err := collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cur.Close(ctx)
// 	for cur.Next(ctx) {
// 		err := cur.Decode(&a)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// do something with result....
// 		fmt.Println(a)

// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// }
