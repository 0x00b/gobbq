package gorm

import (
	"reflect"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Tabler interface {
	TableName(c *entity.Context, i any) string
}

// func TableName(c *entity.Context, i any) string {
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

func (gdb *GormDB) Table(c *entity.Context, i any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(gdb.TableName(c, i))
	}
}

func RawTableName(i any) string {
	v := reflect.TypeOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return schema.NamingStrategy{}.TableName(v.Name())
}

type defaultTabler struct{}

func (defaultTabler) TableName(c *entity.Context, i any) string {
	return RawTableName(i)
}
