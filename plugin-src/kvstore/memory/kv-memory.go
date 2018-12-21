// This is a very ugly in-memory kv store
// Not optimizied for memory
// But it works...
package main

import (
	"sync"
)

type KVMem struct {
	dbname  string
	data    map[string]string
	rwmutex sync.RWMutex
}

var kvmem KVMem

func GetInstance() interface{} {
	return &kvmem
}

func (store *KVMem) Init(conf string) {
	store.dbname = conf
	store.data = make(map[string]string)
}

func (store *KVMem) Get(key string) (value string, err error) {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	return store.data[key], nil
}

func (store *KVMem) Set(key string, value string) (err error) {
	store.rwmutex.Lock()
	defer store.rwmutex.Unlock()

	store.data[key] = value

	return
}

func (store *KVMem) Count() (count int64) {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	return (int64)(len(store.data))
}
