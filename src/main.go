package main

import (
	"fmt"
	"kvstore"
	"pluginmgr"
	"sync"
)

var pluginMgr *pluginmgr.PluginMgr
var wg sync.WaitGroup

func main() {
	var err error
	pluginMgr, err = pluginmgr.CreateMgr("conf/plugin.conf")
	if err != nil {
		fmt.Println("plugin load failed,err=", err)
		return
	}

	wg.Add(2)
	go runKvStore("Memory")
	go runKvStore("Json")

	wg.Wait()
}

func runKvStore(storetype string) {
	defer wg.Done()
	/* test for kvstore*/
	store, err := kvstore.NewKVStore(pluginMgr, storetype)
	if err != nil {
		fmt.Println("NewKVStore failed,err=", err)
		return
	}

	store.Init("data")

	// set several keys
	store.Set("key1", "test")
	store.Set("key2", "test2")

	if storetype == "Memory" {
		store.Set("key2", "test2-2")
		store.Set("key3", "test3")
		store.Set("key4", "test4")
	}

	// output kv pairs count
	fmt.Println("Store", storetype, "now kvstore has", store.Count(), "keys")

	// get the value of key2
	data, err := store.Get("key2")
	if err != nil {
		fmt.Println("KVStore read failed,err=", err)
		return
	}
	fmt.Println("Store", storetype, "key2=", string(data))
}
