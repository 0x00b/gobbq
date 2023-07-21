package gogen

import (
	"fmt"
	"strings"

	"github.com/0x00b/gobbq/proto/bbq"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type tagItem struct {
	key   string
	value string
}

type tagItems []tagItem

func (ti tagItems) format() string {
	tags := []string{}
	for _, item := range ti {
		tags = append(tags, fmt.Sprintf(`%s:%s`, item.key, item.value))
	}
	return strings.Join(tags, " ")
}

func (ti tagItems) override(nti tagItems) tagItems {
	overrided := []tagItem{}
	for i := range ti {
		dup := -1
		for j := range nti {
			if ti[i].key == nti[j].key {
				dup = j
				break
			}
		}
		if dup == -1 {
			overrided = append(overrided, ti[i])
		} else {
			overrided = append(overrided, nti[dup])
			nti = append(nti[:dup], nti[dup+1:]...)
		}
	}
	return append(overrided, nti...)
}

func newTagItems(tag string) tagItems {
	items := []tagItem{}
	splitted := rTags.FindAllString(tag, -1)

	for _, t := range splitted {
		sepPos := strings.Index(t, ":")
		items = append(items, tagItem{
			key:   t[:sepPos],
			value: t[sepPos+1:],
		})
	}
	return items
}

func injectTag(contents []byte, area textArea) (injected []byte) {
	expr := make([]byte, area.End-area.Start)
	copy(expr, contents[area.Start-1:area.End-1])
	cti := newTagItems(area.CurrentTag)
	iti := newTagItems(area.InjectTag)
	ti := cti.override(iti)
	expr = rInject.ReplaceAll(expr, []byte(fmt.Sprintf("`%s`", ti.format())))

	injected = append(injected, contents[:area.Start-1]...)
	injected = append(injected, expr...)
	injected = append(injected, contents[area.End-1:]...)

	return
}

func FieldTag(f *protogen.Field) string {

	tag := "" //`json:"` + f.Desc.JSONName() + `,omitempty"`

	v, ok := proto.GetExtension(f.Desc.Options(), bbq.E_Field).(*bbq.Field)
	if !ok {
		return "" //"`" + tag + "`"
	}

	switch v.GetMgo() {
	case bbq.MONGO_MGO_NONE:
		tag += ` bson:"-"`
	case bbq.MONGO_MGO_FIELD:
		tag += ` bson:"` + f.Desc.JSONName() + `"`
	case bbq.MONGO_MGO_ID:
		tag += ` bson:"_id"`
	}

	switch v.GetMysql() {
	case bbq.MySQL_MYSQL_NONE:
	case bbq.MySQL_MYSQL_FIELD:
	case bbq.MySQL_MYSQL_PRIMARY_KEY:
	case bbq.MySQL_MYSQL_UNIQUE:
	case bbq.MySQL_MYSQL_INDEX:
	}

	return "`" + tag + "`"
}
