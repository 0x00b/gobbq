package redis

import (
	"io"

	"github.com/0x00b/gobbq/engine/db/kv"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	keyPrefix = "_KV_"
)

type redisKVDB struct {
	c *redis.Client
}

// OpenRedisKVDB opens Redis for KVDB backend
func OpenRedisKVDB(url string, dbindex int) (kv.KVDBEngine, error) {
	// c, err := redis.(url)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "redis dail failed")
	// }

	opt := &redis.Options{}
	cli := redis.NewClient(opt)

	db := &redisKVDB{
		c: cli,
	}

	if err := db.initialize(dbindex); err != nil {
		panic(errors.Wrap(err, "redis kvdb initialize failed"))
	}

	return db, nil
}

func (db *redisKVDB) initialize(dbindex int) error {
	if dbindex >= 0 {
		if err := db.c.Do(context.Background(), "SELECT", dbindex); err != nil {
			return err.Err()
		}
	}

	return nil
}

func (db *redisKVDB) isZeroCursor(c any) bool {
	return string(c.([]byte)) == "0"
}

func (db *redisKVDB) Get(ctx entity.Context, key string) (val string, err error) {
	cmd := db.c.Do(ctx, "GET", keyPrefix+key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	if cmd.Val() == nil {
		return "", nil
	}
	return string(cmd.Val().([]byte)), err
}

func (db *redisKVDB) Put(ctx entity.Context, key string, val string) error {
	err := db.c.Do(ctx, "SET", keyPrefix+key, val)
	return err.Err()
}

type redisKVDBIterator struct {
	db       *redisKVDB
	leftKeys []string
}

func (it *redisKVDBIterator) Next(ctx entity.Context) (kv.KVItem, error) {
	if len(it.leftKeys) == 0 {
		return kv.KVItem{}, io.EOF
	}

	key := it.leftKeys[0]
	it.leftKeys = it.leftKeys[1:]
	val, err := it.db.Get(ctx, key)
	if err != nil {
		return kv.KVItem{}, err
	}

	return kv.KVItem{Key: key, Val: val}, nil
}

func (db *redisKVDB) Find(ctx entity.Context, beginKey string, endKey string) (kv.Iterator, error) {
	return nil, errors.Errorf("operation not supported on redis")
}

func (db *redisKVDB) Close() {
	db.c.Close()
}

func (db *redisKVDB) IsConnectionError(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}
