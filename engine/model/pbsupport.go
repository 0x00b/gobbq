// Package pbsupport protobuf support used in db package.
package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

var (
	_BlobKind = map[protoreflect.Kind]struct{}{
		protoreflect.BytesKind:   {},
		protoreflect.MessageKind: {},
		protoreflect.GroupKind:   {},
	}
	_marshalOptions        = &proto.MarshalOptions{}
	_unmarshalMergeOptions = &proto.UnmarshalOptions{
		Merge: true,
	}
)

// MarshalToMap recode split and marshal data to map, scalar will be marshal to string.
func MarshalToMap(msg proto.Message, fields []FieldName) (map[string]any, error) {
	rf := msg.ProtoReflect()
	desc := rf.Descriptor()
	rawStrMap := map[string]string{}
	rf.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if !IsFdMarshalToBlob(fd) {
			rawStrMap[string(fd.Name())] = marshalScalar(fd, v)
		}
		return true
	})

	fdFilter := DirtyFilter(desc, fields)
	fds := desc.Fields()
	// TODO: 这里需要优化，不能调用marshalField函数还是会进行多余的编码.
	wireMap := map[string]any{}
	var cur int32
	marshalBytes, err := _marshalOptions.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal to map err:%w", err)
	}
	for int(cur) < len(marshalBytes) {
		num, _, n := protowire.ConsumeField(marshalBytes[cur:])
		if n < 0 {
			return nil, fmt.Errorf("wire consume field ret=%d", n)
		}
		next := cur + int32(n)
		fd := fds.ByNumber(num)
		if fd == nil {
			return nil, fmt.Errorf("unknow field num=%d", num)
		}
		fdName := string(fd.Name())
		if fdFilter != nil {
			if _, exist := fdFilter[fdName]; !exist {
				cur = next
				continue
			}
		}
		rstr, exist := rawStrMap[fdName]
		if exist {
			wireMap[fdName] = rstr
		} else {
			cb, exist := wireMap[fdName]
			if !exist {
				wireMap[fdName] = marshalBytes[cur:next]
			} else {
				wireMap[fdName] = append(cb.([]byte), marshalBytes[cur:next]...)
			}
		}
		cur = next
	}
	return wireMap, nil
}

// DirtyFilter Create dirty filter and ignore deep level fields.
func DirtyFilter(desc protoreflect.MessageDescriptor, fields []FieldName) map[string]struct{} {
	if len(fields) == 0 {
		return nil
	}
	fds := desc.Fields()
	retFilter := map[string]struct{}{}
	for _, fs := range fields {
		retKey := fs
		fd := fds.ByName(protoreflect.Name(fs))
		if fd == nil {
			// TODO:need optimize.
			ir := strings.Index(string(fs), ".")
			if ir < 0 {
				fmt.Printf("dirty field=%s not in msg=%s", fs, desc.FullName())
				continue
			}
			retKey = fs[:ir]
			fd = fds.ByName(protoreflect.Name(retKey))
			if fd == nil {
				fmt.Printf("dirty field=%s not in msg=%s", fs, desc.FullName())
				continue
			}
		}
		retFilter[string(retKey)] = struct{}{}
	}
	return retFilter
}

// UnmarshalFromMap unmarsh data loaded from db.
func UnmarshalFromMap(msg proto.Message, bytesMap map[string]string) (err error) {
	buf := bytes.Buffer{}
	rf := msg.ProtoReflect()
	desc := rf.Descriptor()
	fds := desc.Fields()
	for k, s := range bytesMap {
		if strings.HasPrefix(k, "_") {
			continue
		}
		fd := fds.ByName(protoreflect.Name(k))
		if fd == nil {
			fmt.Printf("cannot find field num=%s type=%s", k, rf.Descriptor().FullName())
			continue
		}
		if IsFdMarshalToBlob(fd) {
			_, e := buf.WriteString(s)
			if e != nil {
				err = fmt.Errorf("merge buf err:%w", e)
				return
			}
		} else {
			v, e1 := unmarshalScalarByStr(fd, s)
			if e1 != nil {
				err = e1
				return
			}
			rf.Set(fd, v)
		}
	}
	err = _unmarshalMergeOptions.Unmarshal(buf.Bytes(), msg)
	return
}

// FindFds find field descriptor, keep order of param KeyNames.
func FindFds(msgDesc protoreflect.MessageDescriptor, keyNames []string) []protoreflect.FieldDescriptor {
	fds := make([]protoreflect.FieldDescriptor, len(keyNames))
	fields := msgDesc.Fields()
	for i, key := range keyNames {
		fds[i] = fields.ByTextName(key)
	}
	return fds
}

func unmarshalScalarByStr(fd protoreflect.FieldDescriptor, str string) (protoreflect.Value, error) {
	const b32 int = 32
	const b64 int = 64
	const base10 = 10

	kind := fd.Kind()
	switch kind {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(str), nil

	case protoreflect.BoolKind:
		switch str {
		case "true", "1":
			return protoreflect.ValueOfBool(true), nil
		case "false", "0", "":
			return protoreflect.ValueOfBool(false), nil
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if n, err := strconv.ParseInt(str, base10, b32); err == nil {
			return protoreflect.ValueOfInt32(int32(n)), nil
		}

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if n, err := strconv.ParseInt(str, base10, b64); err == nil {
			return protoreflect.ValueOfInt64(n), nil
		}

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if n, err := strconv.ParseUint(str, base10, b32); err == nil {
			return protoreflect.ValueOfUint32(uint32(n)), nil
		}

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if n, err := strconv.ParseUint(str, base10, b64); err == nil {
			return protoreflect.ValueOfUint64(n), nil
		}
	case protoreflect.DoubleKind:
		if n, err := strconv.ParseFloat(str, b64); err == nil {
			return protoreflect.ValueOfFloat64(n), nil
		}
	case protoreflect.FloatKind:
		if n, err := strconv.ParseFloat(str, b64); err == nil {
			return protoreflect.ValueOfFloat32(float32(n)), nil
		}
	}

	return protoreflect.Value{}, fmt.Errorf("invalid value for fd=%s value=%s", fd.Name(), str)
}

func marshalScalar(fd protoreflect.FieldDescriptor, v protoreflect.Value) (retstr string) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		if v.Bool() {
			retstr = "1"
		} else {
			retstr = "0"
		}
	default:
		retstr = v.String()
	}
	return
}

// IsFdMarshalToBlob if a field will marshal to blob.
func IsFdMarshalToBlob(fd protoreflect.FieldDescriptor) bool {
	if fd.Cardinality() == protoreflect.Repeated {
		return true
	}
	_, exist := _BlobKind[fd.Kind()]
	return exist
}
