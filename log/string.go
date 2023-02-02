package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/0x00b/gobbq/tool/utils"
	"google.golang.org/protobuf/proto"
)

//打印请求包或者回包时，对报文内容进行处理：
// 1、如果pb协议没有过长，直接使用pb.String()
// 2、过长则对结构体中的字段进行处理，参考下面的可配置变量

// use this
//
// String 入参换成可打印的json string，包括protobuf和普通结构体，或者其他类型
// 返回 fmt.Stringer， 防止日志级别不够时，还执行这个很消耗资源的操作
func String(i any) fmt.Stringer {
	return &defaultStringer{i}
}

var (
	// PrintProtoLen 设置打印pb报文的长度，超过这个长度将会对报文中的超长字段进行处理（截断）
	PrintProtoLen = 1024

	// PrintStringLen 设置打印报文中字符串的长度，超过这个长度将会对报文中的超长字段进行处理（截断）
	PrintStringLen = 64

	// PrintSliceLen 设置打印报文中的slice或者array的长度，超过这个长度将会对报文中的超长内容进行处理（截断）
	PrintSliceLen = 64

	// 处理string的方式，默认截断后面补上 "..."
	DealStringHook DealStringHookFunc = getStringSlice
)

type DealStringHookFunc func(string) string

var nilstr = "nil"

type defaultStringer struct {
	i any
}

func (d *defaultStringer) String() string {
	if d.i == nil {
		return nilstr
	}
	typ := reflect.TypeOf(d.i)
	val := reflect.ValueOf(d.i)
	for typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nilstr
		}
		val = val.Elem()
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Invalid:
		return FmtString(val.Interface())
	case reflect.Struct:
		return StructToString(val.Interface())
	default:
		return JsonString(d.i)
	}
}

// 把定义了Json字段的结构本转换为Json字符串输出
func JsonString(r any) string {
	b, _ := json.Marshal(r)
	return string(b)
}

// String 把i转换成string
func FmtString(i any) string {
	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.String {
		return reflect.ValueOf(i).String()
	}
	if typ.Kind() == reflect.Ptr {
		return String(reflect.ValueOf(i).Elem()).String()
	}
	return fmt.Sprint(i)
}

// ProtoToPrintString pb协议转换成可打印的json string，会对超长内容截断
func ProtoToPrintString(p proto.Message) string {
	if p == nil {
		return ""
	}
	data, err := proto.Marshal(p)
	str := string(data)
	if err == nil {
		if PrintProtoLen <= 0 || len(str) < PrintProtoLen {
			return str
		}
	}
	m, e := GetStructFields(p)
	if e == nil {
		data, e := json.Marshal(m)
		if e == nil {
			return string(data)
		}
	}
	return ""

}

// StructToPrintString 结构体转换成可打印的json string，会对超长内容截断
func StructToString(p any) string {
	if p == nil {
		return nilstr
	}
	if pm, ok := p.(proto.Message); ok {
		return ProtoToPrintString(pm)
	}
	m, e := GetStructFields(p)
	if e == nil {
		data, e := json.Marshal(m)
		if e == nil {
			return string(data)
		}
	}
	return ""
}

//
//func handleString(m map[string]any, key string, s any) {
//	if s, ok := s.(string); ok {
//		if len(s) < PrintStringLen {
//			m[key] = s
//		} else {
//			m[key] = s[:PrintStringLen] + "..."
//		}
//		return
//	}
//
//	if s, ok := s.(*string); ok {
//		if len(*s) < PrintStringLen {
//			m[key] = s
//		} else {
//			m[key] = (*s)[:PrintStringLen] + "..."
//		}
//		return
//	}
//}

func getStringSlice(s string) string {
	if PrintStringLen <= 0 || len(s) < PrintStringLen {
		return s
	}
	return s[:PrintStringLen] + "..."
}

// nolint
func handlerField(fields map[string]any,
	fieldName string, field reflect.Value, t reflect.Type, anonymous bool) {
	// type field
	if _, ok := fields[fieldName]; ok {
		return
	}
	kind := t.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		var values []any
		moreFields := false
		fLen := field.Len()
		if PrintSliceLen > 0 && fLen > PrintSliceLen {
			field = field.Slice(0, PrintSliceLen)
			moreFields = true
		}
		defer func() {
			if len(values) > 0 {
				if moreFields {
					//多余的字段省略，使用...代替
					values = append(values, fmt.Sprintf("(more %d item)...", fLen-PrintSliceLen))
				}
				fields[fieldName] = values
			}
		}()
		var sliceKind reflect.Kind
		bPtr := false
		if field.Type().Elem().Kind() == reflect.Ptr {
			sliceKind = field.Type().Elem().Elem().Kind()
			bPtr = true
		} else {
			sliceKind = field.Type().Elem().Kind()
		}
		switch sliceKind {
		case reflect.Bool:
			fallthrough
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			//fields[fieldName] = field.Interface()
			for i := 0; i < field.Len(); i++ {
				values = append(values, field.Index(i).Interface())
			}

		case reflect.String:
			for i := 0; i < field.Len(); i++ {
				var t *string
				if bPtr {
					if !field.Index(i).IsNil() {
						str := DealStringHook(field.Index(i).Elem().String())
						t = &str
					}
				} else {
					str := DealStringHook(field.Index(i).String())
					t = &str
				}
				values = append(values, t)
			}
		case reflect.Struct:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]any)
				if field.Type().Elem().Kind() == reflect.Ptr {
					if !field.Index(i).IsNil() {
						handlerStruct(t, reflect.ValueOf(field.Index(i).Elem().Interface()),
							reflect.TypeOf(field.Index(i).Elem().Interface()))
					}
				} else {
					handlerStruct(t, reflect.ValueOf(field.Index(i).Interface()),
						reflect.TypeOf(field.Index(i).Interface()))
				}
				if len(t) > 0 {
					values = append(values, t)
				}
			}
		case reflect.Interface:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]any)
				if !field.Index(i).IsNil() {
					handlerField(t, fieldName, field.Index(i).Elem(), field.Index(i).Elem().Type(), false)
				}
				if len(t) > 0 {
					values = append(values, t[fieldName])
				}
			}
		case reflect.Map:
			for i := 0; i < field.Len(); i++ {
				t := make(map[string]any)
				if field.Type().Elem().Kind() == reflect.Ptr {
					if !field.Index(i).IsNil() {
						handlerField(t, fieldName, field.Index(i).Elem(), field.Index(i).Elem().Type(), false)
					}
				} else {
					handlerField(t, fieldName, field.Index(i), field.Index(i).Type(), false)
				}
				if len(t) > 0 {
					values = append(values, t[fieldName])
				}
			}
		default:
			fmt.Printf("type:%v", sliceKind)
			//panic("reflect.TypeOf(param).Elem().Kind() no setting")
		}
	case reflect.String:
		str := DealStringHook(field.String())
		fields[fieldName] = str
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Bool:
		fields[fieldName] = field.Interface()
	case reflect.Struct:
		if anonymous {
			handlerStruct(fields, reflect.ValueOf(field.Interface()), reflect.TypeOf(field.Interface()))
		} else {
			temp := make(map[string]any)
			handlerStruct(temp, reflect.ValueOf(field.Interface()), reflect.TypeOf(field.Interface()))
			fields[fieldName] = temp
		}
	case reflect.Ptr:
		if !field.IsNil() {
			handlerField(fields, fieldName, field.Elem(), t.Elem(), false)
		}
	case reflect.Interface:
		if !field.IsNil() {
			handlerField(fields, fieldName, field.Elem(), field.Elem().Type(), false)
		}
	case reflect.Map:
		if !field.IsNil() {
			temp := make(map[string]any)
			for _, key := range field.MapKeys() {
				handlerField(temp, formatAtom(key), field.MapIndex(key), field.MapIndex(key).Type(), false)
			}
			fields[fieldName] = temp
		}
	default:
		//panic("reflect.Type.Kind() no setting")
	}
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	case reflect.Chan, reflect.Func, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Interface, reflect.Ptr:
		return formatAtom(v.Elem())
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
func handlerStruct(fields map[string]any, v reflect.Value, t reflect.Type) {
	for i := 0; i < v.Type().NumField(); i++ {
		if !v.Field(i).CanInterface() || !v.Field(i).IsValid() ||
			(isOmitEmpty(v.Type().Field(i)) && utils.IsEmptyValue(v.Field(i))) {
			continue
		}
		fieldName := utils.GetJsonName(t.Field(i))
		if fieldName != "-" {
			handlerField(fields, fieldName, v.Field(i), v.Type().Field(i).Type, t.Field(i).Anonymous)
		}
	}
}

func isOmitEmpty(field reflect.StructField) bool {
	tag := field.Tag.Get("json")
	return strings.Contains(tag, "omitempty")
}

// GetStructFields 把结构体转换成一个map
func GetStructFields(st any) (fields map[string]any, err error) {
	fields = make(map[string]any)
	if st == nil {
		return
	}
	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	if v.Kind() == reflect.Ptr {
		if v.Elem().Type().Kind() == reflect.Struct {
			handlerStruct(fields, v.Elem(), v.Elem().Type())
			return
		}
	}
	switch v.Kind() {
	case reflect.Struct:
		handlerStruct(fields, v, t)
	default:
		err = errors.New("Can't handler type: " + v.Type().String())
	}
	return
}
