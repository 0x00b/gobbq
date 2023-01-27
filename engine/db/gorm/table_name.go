package gorm

import (
	"context"
	"reflect"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Tabler interface {
	TableName(c context.Context, i interface{}) string
}

// func TableName(c context.Context, i interface{}) string {
// 	vt := reflect.TypeOf(i)
// 	vv := reflect.ValueOf(i)
// 	for vt.Kind() == reflect.Ptr {
// 		vt = vt.Elem()
// 		vv = vv.Elem()
// 	}
// 	name := ""
// 	appid := xpaas.AppID(c)
// 	key := vt.Name() + appid

// 	if cacheName, ok := tableNameCache.Get(key); ok {
// 		name, _ = cacheName.(string)
// 	} else {
// 		name = RawTableName(i)
// 		if appid != "" {
// 			name = "a" + appid + "_" + name
// 		}
// 		tableNameCache.Set(key, name, 0)
// 	}
// 	return name
// }

var tableNameCache *cache.Cache = cache.New(0, 0)

func (gdb *GormDB) Table(c context.Context, i interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(gdb.TableName(c, i))
	}
}

func RawTableName(i interface{}) string {
	v := reflect.TypeOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return schema.NamingStrategy{}.TableName(v.Name())
}

type defaultTabler struct{}

func (defaultTabler) TableName(c context.Context, i interface{}) string {
	return RawTableName(i)
}
