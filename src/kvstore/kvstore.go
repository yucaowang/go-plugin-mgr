package kvstore

import "pluginmgr"

const KV_PLUGIN_NAME = "kvstore"

type KVStore interface {
	Init(string)
	Get(string) (string, error)
	Set(string, string) error
	Count() int64
}

func NewKVStore(pm *pluginmgr.PluginMgr) (store KVStore, err error) {
	var iface interface{}
	iface, err = pm.GetPlugin(KV_PLUGIN_NAME)
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
