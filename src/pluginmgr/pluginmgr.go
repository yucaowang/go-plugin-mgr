/**
* This is a plugin manager which keeps all plugins' instance
* Notice: any plugin using this framework should implements
*         a func named 'GetInstance' to return there instance
**/
package pluginmgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"plugin"
)

type PluginMgr struct {
	pluginConf      map[string]confNode
	pluginInstances map[string]interface{}
}

type confNode struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

// public functions start from here

func (pm *PluginMgr) Init() {
	// init config struct
	pm.pluginConf = make(map[string]confNode)
	pm.pluginInstances = make(map[string]interface{})

	// Read conf file to get all plugins
	pm.readPluginConfig()

	// Load plugins
	pm.loadPlugins()
}

func (pm *PluginMgr) GetPlugin(name string) (pluginInstance interface{}, err error) {
	if pm.pluginInstances == nil {
		err = errors.New("Not initialized! please call Init firstly")
	}

	pluginInstance = pm.pluginInstances[name]
	return
}

// internal functions start from here
func (pm *PluginMgr) readPluginConfig() {
	// get config file
	f, err := os.Open("conf/plugin.conf")
	if err != nil {
		fmt.Println("load plugin conf file failed, cannot open file")
	}
	defer f.Close()

	// read config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&pm.pluginConf)
	if err != nil {
		fmt.Println("load plugin conf file failed while parsing")
	}

	pluginCount := len(pm.pluginConf)
	fmt.Println("Plugin conf number is:", pluginCount, "data=", pm.pluginConf)
}

func (pm *PluginMgr) loadPlugins() {
	for name, conf := range pm.pluginConf {
		// open plugin
		pg, err := plugin.Open(conf.Path)
		if err != nil {
			fmt.Println("Warn: plugin", name, "open failed, err=", err)
			continue
		}

		// get instance
		iSymbol, err := pg.Lookup("GetInstance")
		if err != nil {
			fmt.Println("Warn: plugin", name, "don't have func named GetInstance, err=", err)
			continue
		}
		pm.pluginInstances[name] = iSymbol.(func() interface{})()
	}
	fmt.Println("Successfully loaded", len(pm.pluginInstances), "plugins")
}
