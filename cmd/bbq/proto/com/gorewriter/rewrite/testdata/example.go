package testdata

import (
	lol "bytes"
	"fmt"
)

// Method TODO
func Method(arg int) {
	// leading comment

	// field comment
	m.Field++

	// trailing comment
}

// Foo TODO
type Foo struct {
	Field int
}

// Method TODO
func (m *Foo) Method(arg int) {
	// leading comment

	// field comment
	m.Field++

	// trailing comment
}

// String TODO
func (m *Foo) String() string {
	var buf lol.Buffer
	buf.WriteString(fmt.Sprintf("%d", m.Field))
	return buf.String()
}
