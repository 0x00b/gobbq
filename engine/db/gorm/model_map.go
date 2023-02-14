package gorm

import (
	"reflect"

	"github.com/0x00b/gobbq/engine/entity"
	"gorm.io/gorm/schema"
)

var (
	naming = schema.NamingStrategy{}
)

// ModelMap 把model转换成map， 为了能够更新空字段
func ModelMap(c entity.Context, i any) map[string]any {

	vv := reflect.ValueOf(i)
	vt := reflect.TypeOf(i)

	table := naming.TableName(vt.Name())

	m := make(map[string]any)
	mLevel := make(map[string]int)

	getFiled(vt, vv, table, m, 0, mLevel)

	return m
}

func getFiled(vt reflect.Type, vv reflect.Value, table string, m map[string]any, level int, mLevel map[string]int) {
	for vt.Kind() == reflect.Ptr {
		vv = vv.Elem()
		vt = vt.Elem()
	}
	if !vv.IsValid() {
		return
	}
	for i := 0; i < vt.NumField(); i++ {

		if vv.Field(i).Interface() != nil {
			ft := vt.Field(i).Type
			fv := vv.Field(i)
			if vt.Field(i).Anonymous {
				getFiled(ft, fv, table, m, level+1, mLevel)
				continue
			}

			column := naming.ColumnName(table, vt.Field(i).Name)
			//对象的上层成员会屏蔽匿名成员的field，因此需要这个判断
			if _, ok := m[column]; !ok || mLevel[column] > level {
				mLevel[column] = level
				m[column] = vv.Field(i).Interface()
			}
		}
	}
}
