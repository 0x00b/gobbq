package xlog_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/0x00b/gobbq/xlog"
	"google.golang.org/protobuf/proto"
)

type Inner struct {
	Inner1 string
	Inner2 []string
	Inner3 []uint32
}

type TestPint struct {
	Inner1 string
	Name   string
	PName  *string
	PNames []*string
	Age    int
	PAge   *uint32

	Obj struct {
		ObjField1 string
		ObjField2 []string
		ObjField3 []*string
	}
	Float float64

	In   Inner
	In0  []Inner
	In1  []*Inner
	In2  *[]Inner
	In3  *[]*Inner
	In4  *[]*Inner
	Data []byte
	Inner
	Interface any
	Int       []*int32
}

func TestGetStructFields(t *testing.T) {
	//PrintStringLen = 64
	//PrintSliceLen = 4
	test := &TestPint{
		Inner1: "out Inner1",
		PNames: []*string{nil, proto.String("xxxxxxxxxxxxxxxxxxxxxxxxx" +
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"), proto.String("pnames")},
		Name:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		PName: proto.String("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
		Age:   10,
		PAge:  proto.Uint32(18),
		Obj: struct {
			ObjField1 string
			ObjField2 []string
			ObjField3 []*string
		}{
			ObjField1: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ObjField2: []string{"1111", "2222", "3333", "444444", "555", "6666", "77777"},
			ObjField3: []*string{nil, proto.String("xxxxxxxxxxxxxxxxxxxxxxxxxx" +
				"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"), proto.String("test")}},

		Inner: Inner{
			Inner1: "inner22",
			Inner2: []string{"2222", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		},

		In: Inner{
			Inner1: "test",
			Inner2: []string{"1111", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		},
		In0: []Inner{{
			Inner1: "test",
			Inner2: []string{"1111", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		}},
		In1: []*Inner{{
			Inner1: "Inner1111",
			Inner2: []string{"1111", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		}, {
			Inner1: "Inner222",
		}, {
			Inner1: "Inner3333",
		}, {
			Inner1: "Inner4444",
		}, {
			Inner1: "Inner5555",
		}, {
			Inner1: "Inner566666",
		}, {
			Inner1: "Inner7777",
		}},
		In2: &([]Inner{{
			Inner1: "test",
			Inner2: []string{"1111", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		}}),
		In3: &[]*Inner{{
			Inner1: "Inner1111",
			Inner2: []string{"1111", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
		}, nil,
		},
		In4:  nil,
		Int:  []*int32{proto.Int32(1), proto.Int32(2), proto.Int32(2), proto.Int32(2), proto.Int32(2), proto.Int32(2)},
		Data: []byte("testttttttttttttttttttttttttttttttt"),
		//Interface: &[]any{proto.String("sss"),map[string]any{"m1":"test"}},
		Interface: &[]map[any]any{{"str": proto.String("sss")},
			{5: proto.String("sss")},
			{nil: map[string]any{"m1": "test"}}, nil},
		Float: 10e7,
	}
	//DealStringHook = func(i any) any {
	//	return i
	//}

	fmt.Println(xlog.String(test))

	fmt.Println("=================json==============")
	js, e := json.Marshal(test)
	fmt.Println(js, e)
	test.Interface = nil
	js, e = json.MarshalIndent(test, "", "	")
	fmt.Println(string(js), e)
}
