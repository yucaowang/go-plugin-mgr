// This is a very ugly text-based kv store
// Very low speed
// But it works...
package main

import (
	"encoding/json"
	"fmt"
	"kvstore"
	"os"
	"sync"
)

type KVJson struct {
	meta     kvstore.Meta
	jsondata map[string]string
	rwmutex  sync.RWMutex
}

func GetInstance() interface{} {
	kvjson := KVJson{}
	return &kvjson
}

func (store *KVJson) Init(conf string) {
	store.meta.DBName = conf + ".json"
	store.meta.Capacity = 0
	store.jsondata = make(map[string]string)
	store.loadData()
}

func (store *KVJson) Get(key string) (value string, err error) {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	return store.jsondata[key], nil
}

func (store *KVJson) Set(key string, value string) (err error) {
	store.rwmutex.Lock()
	defer store.rwmutex.Unlock()

	store.jsondata[key] = value
	store.flush()

	return
}

func (store *KVJson) Count() (count int64) {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	return (int64)(len(store.jsondata))
}

func (store *KVJson) GetMeta() (meta *kvstore.Meta) {
	return &store.meta
}

func (store *KVJson) loadData() {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	// get config file
	f, err := os.Open(store.meta.DBName)
	if err != nil {
		fmt.Println("load kv data file failed, cannot open file, err=", err)
	}
	defer f.Close()

	// read config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&store.jsondata)
	if err != nil {
		fmt.Println("load kv data file failed while parsing, err=", err)
	}

}

func (store *KVJson) flush() (err error) {
	f, err := os.OpenFile(store.meta.DBName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Open file failed")
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(store.jsondata)
	return
}
