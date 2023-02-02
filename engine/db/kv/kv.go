package kv

import "github.com/0x00b/gobbq/engine/entity"

// KVDBEngine defines the interface of a KVDB engine implementation
type KVDBEngine interface {
	Get(ctx *entity.Context, key string) (val string, err error)
	Put(ctx *entity.Context, key string, val string) (err error)
	Find(ctx *entity.Context, beginKey string, endKey string) (Iterator, error)
	Close()
	IsConnectionError(err error) bool
}

// Iterator is the interface for iterators for KVDB
//
// Next should returns the next item with error=nil whenever has next item
// otherwise returns KVItem{}, io.EOF
// When failed, returns KVItem{}, error
type Iterator interface {
	Next(ctx *entity.Context) (KVItem, error)
}

// KVItem is the type of KVDB item
type KVItem struct {
	Key string
	Val string
}
