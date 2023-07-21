package db

import (
	"context"

	"github.com/0x00b/gobbq/engine/model"
	"google.golang.org/protobuf/proto"
)

// Record db plugin support record type.
type Record = proto.Message

type DBName string

const (
	DBMySQL DBName = "mysql"
	DBMongo DBName = "mongo"
	DBRedis DBName = "redis"
)

// IDatabase DB接口类.
type IDatabase interface {
	Name() DBName

	Connect(config any) error
	Table(name string) (IDatabase, error)

	Load(c context.Context, record Record) error // get by id
	Update(c context.Context, record Record, fields []model.FieldName) error
	Insert(c context.Context, record Record) error
	Delete(c context.Context, record Record) error
	Save(c context.Context, record Record) error     // insert or update
	AutoSave(c context.Context, record Record) error // just save updated field

	model.Watcher
}
