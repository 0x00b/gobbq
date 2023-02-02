package es_test

// func TestEs(t *testing.T) {
// 	c := context.Background()
// 	cli, _ := elasticsearch.NewClient(elasticsearch.Config{
// 		Logger: es.Monitor{DumpHttp: true},
// 	})

// 	docID := ""

// 	var body strings.Builder
// 	body.Reset()
// 	body.WriteString(`{"index" : { "_index" : "test", "_type" : "_doc", "_id" : "` + docID + `" }}`)
// 	body.WriteString(`{"foo" : "bar `)
// 	body.WriteString(docID)
// 	body.WriteString(`	" }`)

// 	_, err := cli.Bulk(
// 		strings.NewReader(body.String()),
// 		cli.Bulk.WithRefresh("true"),
// 		cli.Bulk.WithPretty(),
// 		cli.Bulk.WithTimeout(100),
// 		cli.Bulk.WithContext(c),
// 	)
// 	if err != nil {
// 		//
// 	}

// }
