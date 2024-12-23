package db

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"sync"
)

var badgerSession = sync.Map{}

// BadgerClient 数据库客户端
type badgerClient struct {
	RootDir string
}

func (b badgerClient) initSession(dbName string) (*badger.DB, error) {
	// 打开新实例
	db, err := badger.Open(badger.DefaultOptions(b.RootDir + "/" + dbName).WithLogger(nil))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetBadgerClient(dbName string) (bd *badger.DB, err error) {
	bdClient, ok := badgerSession.Load(dbName)
	if !ok {
		bdClient, err = badgerClient{"data"}.initSession(dbName)
		if err != nil {
			return nil, err
		}
		badgerSession.Store(dbName, bdClient)
	}
	db, ok := bdClient.(*badger.DB)
	if !ok {
		return nil, errors.New("badger db error")
	}
	return db, nil

}

func CloseAllBadgerClient() {
	badgerSession.Range(func(k, v interface{}) bool {
		if dbClient, ok := v.(*badger.DB); ok {
			if err := dbClient.Close(); err == nil {
				badgerSession.Delete(k)
				fmt.Printf("数据库客户端连接已关闭: %s\n", k)
			}
		}
		return true
	})
}
