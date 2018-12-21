// This is a very ugly text-based kv store
// Very low speed
// But it works...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type KVJson struct {
	dbname   string
	jsondata map[string]string
	rwmutex  sync.RWMutex
}

var kvjson KVJson

func GetInstance() interface{} {
	return &kvjson
}

func (store *KVJson) Init(conf string) {
	store.dbname = conf + ".json"
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

func (store *KVJson) loadData() {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	// get config file
	f, err := os.Open(store.dbname)
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
	f, err := os.OpenFile(store.dbname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Open file failed")
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(store.jsondata)
	return
}
