package main

import (
	"fmt"
	"kvstore"
	"pluginmgr"
)

var pluginMgr *pluginmgr.PluginMgr

func main() {
	pluginMgr = new(pluginmgr.PluginMgr)
	pluginMgr.Init()

	store, err := kvstore.NewKVStore(pluginMgr)
	if err != nil {
		fmt.Println("NewKVStore failed,err=", err)
	}

	store.Init("data")

	// set several keys
	store.Set("key1", "test")
	store.Set("key2", "test2")
	store.Set("key3", "test3")
	store.Set("key4", "test4")

	// output kv pairs count
	fmt.Println("now kvstore has", store.Count(), "keys")

	// get the value of key2
	data, err := store.Get("key2")
	if err != nil {
		fmt.Println("KVStore read failed,err=", err)
	}
	fmt.Println("key2=", string(data))
}
