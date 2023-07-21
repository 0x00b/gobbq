package mongo

import (
	"errors"
	"reflect"
	"strings"

	"github.com/0x00b/gobbq/engine/db"
	"github.com/0x00b/gobbq/engine/model"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/0x00b/gobbq/tool/utils"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

func GetMongoID(r db.Record) (any, error) {
	rf := r.ProtoReflect()
	desc := rf.Descriptor()
	fds := desc.Fields()
	for i := 1; i < fds.Len(); i++ {

		f := fds.ByNumber(protowire.Number(i))

		v, ok := proto.GetExtension(f.Options(), bbq.E_Field).(*bbq.Field)
		if !ok {
			continue
		}

		switch v.GetMgo() {
		case bbq.MONGO_MGO_NONE:
		case bbq.MONGO_MGO_FIELD:
		case bbq.MONGO_MGO_ID:
			v := rf.Get(f).Interface()
			vv := reflect.ValueOf(v)
			if !vv.IsValid() || vv.IsZero() {
				return nil, errors.New("empty id")
			}
			return v, nil
		}
	}

	return nil, errors.New("not found")
}

func (m *mongoDB) PartialMarshalToMap(msg proto.Message, fields []model.FieldName) (map[string]any, error) {

	rt := reflect.TypeOf(msg)
	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	rv := reflect.ValueOf(msg)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	wireMap := map[string]any{}

	for _, fname := range fields {

		ft, ok := rt.FieldByName(string(fname))
		if !ok {
			// fmt.Println("no filed", string(fname))
			continue
		}

		fv := rv.FieldByName(string(fname))

		// fmt.Println("get filed", string(fname))

		tag := utils.GetFiledTag(ft, "bson")
		k := strings.ToLower(string(fname))
		if len(tag) > 0 {
			if tag[0] == "-" {
				continue
			}
			k = tag[0]
		}

		wireMap[k] = fv.Interface()
	}

	return wireMap, nil
}
