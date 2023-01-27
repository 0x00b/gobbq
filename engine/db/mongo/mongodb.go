package kvdbmongo

import (
	"context"

	"github.com/0x00b/gobbq/engine/db/kv"
	"github.com/0x00b/gobbq/log"
	"gopkg.in/mgo.v2"

	"io"

	"gopkg.in/mgo.v2/bson"
)

const (
	_DEFAULT_DB_NAME = "goworld"
	_VAL_KEY         = "_"
)

type mongoKVDB struct {
	s *mgo.Session
	c *mgo.Collection
}

// OpenMongoKVDB opens mongodb as KVDB engine
func OpenMongoKVDB(url string, dbname string, collectionName string) (kv.KVDBEngine, error) {
	log.Debugln(context.Background(), "Connecting MongoDB ...")
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	if dbname == "" {
		// if db is not specified, use default
		dbname = _DEFAULT_DB_NAME
	}
	db := session.DB(dbname)
	c := db.C(collectionName)
	return &mongoKVDB{
		s: session,
		c: c,
	}, nil
}

func (kvdb *mongoKVDB) Put(ctx context.Context, key string, val string) error {
	_, err := kvdb.c.UpsertId(key, map[string]string{
		_VAL_KEY: val,
	})
	return err
}

func (kvdb *mongoKVDB) Get(ctx context.Context, key string) (val string, err error) {
	q := kvdb.c.FindId(key)
	var doc map[string]string
	err = q.One(&doc)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = nil
		}
		return
	}
	val = doc[_VAL_KEY]
	return
}

type mongoKVIterator struct {
	it *mgo.Iter
}

func (it *mongoKVIterator) Next(ctx context.Context) (kv.KVItem, error) {
	var doc map[string]string
	ok := it.it.Next(&doc)
	if ok {
		return kv.KVItem{
			Key: doc["_id"],
			Val: doc["_"],
		}, nil
	}

	err := it.it.Close()
	if err != nil {
		return kv.KVItem{}, err
	}
	return kv.KVItem{}, io.EOF
}

func (kvdb *mongoKVDB) Find(ctx context.Context, beginKey string, endKey string) (kv.Iterator, error) {
	q := kvdb.c.Find(bson.M{"_id": bson.M{"$gte": beginKey, "$lt": endKey}})
	it := q.Iter()
	return &mongoKVIterator{
		it: it,
	}, nil
}

func (kvdb *mongoKVDB) Close() {
	kvdb.s.Close()
}

func (kvdb *mongoKVDB) IsConnectionError(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}
