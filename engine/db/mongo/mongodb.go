package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/model"
	"github.com/0x00b/gobbq/xlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	_DEFAULT_DB_NAME = "gobbq"
)

type Config struct {
	URL            string `yaml:"url"`
	DBName         string `yaml:"db"`
	CollectionName string `yaml:"collection"`
}

var _ db.IDatabase = &mongoDB{}

type mongoDB struct {
	client   *mongo.Client
	database *mongo.Database

	collection *mongo.Collection

	watchModels map[proto.Message]*watchMessage
}

func NewMongoDB() *mongoDB {
	d := &mongoDB{
		watchModels: make(map[protoreflect.ProtoMessage]*watchMessage),
	}

	return d
}

func (m *mongoDB) Name() db.DBName {
	return db.DBMongo
}

func (m *mongoDB) Connect(cfg any) error {

	mcfg, ok := cfg.(*Config)
	if !ok {
		mtcfg, ok := cfg.(Config)
		if !ok {
			return errors.New("err config")
		}
		mcfg = &mtcfg
	}

	xlog.Debugln(nil, "Connecting mongoDB ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mcfg.URL))
	if err != nil {
		return err
	}

	if mcfg.DBName == "" {
		mcfg.DBName = _DEFAULT_DB_NAME
	}

	m.client = client
	m.database = client.Database(mcfg.DBName)

	if mcfg.CollectionName != "" {
		m.collection = m.database.Collection(mcfg.CollectionName)
	}

	return nil
}

func (m *mongoDB) Table(name string) (db.IDatabase, error) {
	if m == nil || m.database == nil {
		return nil, errors.New("nil db")
	}
	if name == "" {
		return nil, errors.New("empty name")
	}
	tm := NewMongoDB()

	tm.client = m.client
	tm.database = m.database

	tm.collection = m.database.Collection(name)

	return tm, nil
}

func (m *mongoDB) Load(c context.Context, record db.Record) error { // get by id

	id, err := GetMongoID(record)
	if err != nil {
		return err
	}

	res := m.collection.FindOne(c, bson.M{"_id": id})
	if res == nil {
		return errors.New("nil res")
	}
	if res.Err() != nil {
		return res.Err()
	}

	return res.Decode(record)
}

func (m *mongoDB) Update(c context.Context, record db.Record, fields []model.FieldName) error {
	if len(fields) == 0 {
		return m.Save(c, record)
	}

	fieldMap, err := m.PartialMarshalToMap(record, fields)
	if err != nil {
		return err
	}

	id, err := GetMongoID(record)
	if err != nil {
		return err
	}

	// map to set param
	param := bson.D{{Key: "$set", Value: bson.M(fieldMap)}}

	res, err := m.collection.UpdateByID(c, id, param)
	if err != nil {
		return err
	}
	_ = res

	return nil
}

func (m *mongoDB) Insert(c context.Context, record db.Record) error {

	res, err := m.collection.InsertOne(c, record)

	if err != nil {
		return err
	}
	_ = res.InsertedID

	return nil
}

func (m *mongoDB) Delete(c context.Context, record db.Record) error {

	id, err := GetMongoID(record)
	if err != nil {
		return err
	}

	res, err := m.collection.DeleteOne(c, bson.M{"_id": id})

	if err != nil {
		return err
	}
	_ = res

	return nil
}

func (m *mongoDB) Save(c context.Context, record db.Record) error {

	id, err := GetMongoID(record)
	if err != nil {
		return err
	}

	opts := options.Update().SetUpsert(true)
	param := bson.D{{Key: "$set", Value: record}}
	res, err := m.collection.UpdateOne(c, bson.M{"_id": id}, param, opts)
	if err != nil {
		return err
	}
	_ = res
	return nil
}

func (m *mongoDB) AutoSave(c context.Context, record db.Record) error { // just save updated field
	return m.updateDirtyField(c, record)
}

func (m *mongoDB) GetIncrementID(c context.Context, idKey string) (uint64, error) {

	if idKey == "" {
		return 0, nil
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := m.database.Collection("IncrementID").FindOneAndUpdate(c, bson.M{"_id": idKey}, bson.M{"$inc": bson.M{"seq": 1}}, opts)
	if res == nil {
		return 0, errors.New("nil res")
	}
	if res.Err() != nil {
		return 0, res.Err()
	}
	type ID struct {
		ID  string `bson:"_id"`
		Seq uint64 `bson:"seq"`
	}
	var id ID

	err := res.Decode(&id)
	if err != nil {
		return 0, err
	}

	return id.Seq, nil
}
