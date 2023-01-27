package es_test

import (
	"context"
	"strings"
	"testing"

	"github.com/0x00b/gobbq/engine/db/es"
	"github.com/elastic/go-elasticsearch/v8"
)

func TestEs(t *testing.T) {
	c := context.Background()
	cli, _ := elasticsearch.NewClient(elasticsearch.Config{
		Logger: es.Monitor{DumpHttp: true},
	})

	docID := ""

	var body strings.Builder
	body.Reset()
	body.WriteString(`{"index" : { "_index" : "test", "_type" : "_doc", "_id" : "` + docID + `" }}`)
	body.WriteString(`{"foo" : "bar `)
	body.WriteString(docID)
	body.WriteString(`	" }`)

	_, err := cli.Bulk(
		strings.NewReader(body.String()),
		cli.Bulk.WithRefresh("true"),
		cli.Bulk.WithPretty(),
		cli.Bulk.WithTimeout(100),
		cli.Bulk.WithContext(c),
	)
	if err != nil {
		//
	}

}
