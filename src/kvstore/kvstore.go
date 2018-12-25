package kvstore

import "pluginmgr"

const KV_PLUGIN_NAME = "kvstore"

type Meta struct {
	DBName   string
	Capacity int
}

type KVStore interface {
	Init(string)
	Get(string) (string, error)
	Set(string, string) error
	Count() int64
	GetMeta() *Meta
}

func NewKVStore(pm *pluginmgr.PluginMgr, subType string) (store KVStore, err error) {
	var iface interface{}
	iface, err = pm.CreatePluginInstance(KV_PLUGIN_NAME, subType)
	if err != nil {
		return
	}

	if iface != nil {
		// registered external plugin
		store = iface.(KVStore)
	} else {
		// no plugin registered, use default one
		store = new(KVText)
	}
	return
}
