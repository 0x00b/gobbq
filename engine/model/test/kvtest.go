package main

import (
	"fmt"

	"github.com/0x00b/gobbq/engine/model"
	"github.com/0x00b/gobbq/example/exampb"
	"google.golang.org/protobuf/encoding/protowire"
)

func main() {

	// str := `{
	// 	"B": "cc"
	// 	}`

	// type A struct {
	// 	A string
	// 	B string
	// }

	// a := A{
	// 	A: "xxx",
	// }
	// e := json.Unmarshal([]byte(str), &a)
	// if e != nil {
	// 	fmt.Println(e)
	// }

	// fmt.Println(a)

	// e := EchoEntity{SayHelloRequest: SayHelloRequest{}}
	// _ = e.CLientID

	v := &exampb.EchoProperty{
		Text: "111",
		Test: &exampb.SayHelloRequest{
			Text:     "222",
			CLientID: 333,
		},
	}

	m, err := model.MarshalToMap(v, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(m)

	// proto.Merge(v, v)

	rf := v.ProtoReflect()
	desc := rf.Descriptor()

	fds := desc.Fields()

	for i := 1; i < fds.Len(); i++ {

		f := fds.ByNumber(protowire.Number(i))
		fmt.Println(f.Name())
	}
}
