package kv

import "context"

// KVDBEngine defines the interface of a KVDB engine implementation
type KVDBEngine interface {
	Get(ctx context.Context, key string) (val string, err error)
	Put(ctx context.Context, key string, val string) (err error)
	Find(ctx context.Context, beginKey string, endKey string) (Iterator, error)
	Close()
	IsConnectionError(err error) bool
}

// Iterator is the interface for iterators for KVDB
//
// Next should returns the next item with error=nil whenever has next item
// otherwise returns KVItem{}, io.EOF
// When failed, returns KVItem{}, error
type Iterator interface {
	Next(ctx context.Context) (KVItem, error)
}

// KVItem is the type of KVDB item
type KVItem struct {
	Key string
	Val string
}
