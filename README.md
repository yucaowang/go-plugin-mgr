# Go Plugin Manager

This is a sample code of go plugin framework

## 1. Code Layout

* `src`: source code folder of main framework
* `conf`: config folder
* `plugin-src`: source code folder of plugins
* `main.go`: main process

## 2. How To Use

The core plugin framework is called `PluginMgr`, this module would load plugin configuration and initialize all plugins' instance.

### 2.1 PluginMgr

There are two public functions in PluginMgr:

* `CreateMgr`: initialize and return a instance of PluginMgr.
* `CreatePluginInstance`: create a instance of plugin by plugin name and subtype, the plugin name and subtype should be the same with the configuration in `plugin.conf`.

### 2.2 Define interface

Abstract common functions and define a interface in the framework. All plugins for the same functional could implement this interface.

For example:

```go
package kvstore

type KVStore interface {
	Init(string)
	Get(string) (string, error)
	Set(string, string) error
}
```

### 2.3 Write your plugin

`PluginMgr` assumes all plugins have a public function named `GetInstance`, so your plugin should implement this function.

```go
func GetInstance() interface{} {
    kvmem := KVMem{}
    return &kvmem
}
```

Then you can import your interface defined in framework and implement all functions.

```go
import "kvstore"

func (ys *YourKVStore) Init(conf string) {
    // Your code here
}

func (ys *YourKVStore) Get(key string) (string, error) {
    // Your code here
}

func (ys *YourKVStore) Set(key string, value string) error {
    // Your code here
}
```

More information about how to make a golang plugin please refer to [The Official Doc](https://golang.org/pkg/plugin/).

### 2.4 Use your plugin

After implemented the interface in your plugin, you can use it in your framework.

First of all, register the plugin in the `plugin.conf`, like this:

```json
{
    "kvstore":{
        "path": "plugins/kv-store.so.1.0",
        "version": "1.0"
    }
}
```

Then in your code, call `PluginMgr.Init()` to open all registered plugins, and use `PluginMgr.GetPlugin(plugin-name)` to get the instance of your plugin.

```go
func NewKVStore(pm *pluginmgr.PluginMgr) (store KVStore, err error) {
    var iface interface{}
    iface, err = pm.GetPlugin("kvstore")
    if err != nil {
        return
    }

    if iface != nil {
        // registered external plugin
        store = iface.(KVStore)
    }
    return
}
```

## 3. Run Sample Code

First of all, clone this desopisity:

> git clone https://github.com/yucaowang/go-plugin-mgr.git

Then compile and install:

```shell
cd go-plugin-mgr
sh build.sh
```

Now you see a `bin` folder contains all binaries and other neccessary things.

```
|bin
|--|conf
|--|--|plugin.conf
|--|plugins
|--|--|kv-json.so.1.0
|--|--|kv-memory.so.1.0
|--|main
```

Finally goto `bin` folder and run `./main` to validate the result.

## 4 License
MIT