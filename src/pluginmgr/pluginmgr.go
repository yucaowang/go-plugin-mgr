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
	pluginConf map[string]map[string]confNode
	//pluginInstances map[string]interface{}
}

type confNode struct {
	SubType  string `json:"subtype"`
	Path     string `json:"path"`
	Version  string `json:"version"`
	OnDemand bool   `json:"ondemand"`
}

// public functions start from here

// CreateMgr returns instance of PluginMgr
func CreateMgr(confPath string) (pm *PluginMgr, err error) {
	pm = new(PluginMgr)
	// init config struct
	pm.pluginConf = make(map[string]map[string]confNode)
	//pm.pluginInstances = make(map[string]interface{})

	// Read conf file to get all plugins
	err = pm.readPluginConfig(confPath)
	if err != nil {
		return
	}

	return
}

// CreatePluginInstance always create new plugin instance
func (pm *PluginMgr) CreatePluginInstance(name string, subtype string) (pluginInstance interface{}, err error) {
	if _, ok := pm.pluginConf[name]; !ok {
		return nil, errors.New("Invalid plugin name")
	}

	if _, ok := pm.pluginConf[name][subtype]; !ok {
		return nil, errors.New("Invalid plugin subtype")
	}

	return pm.loadOnePlugin(name, subtype)
}

// internal functions start from here
func (pm *PluginMgr) readPluginConfig(confPath string) error {
	// get config file
	f, err := os.Open(confPath)
	if err != nil {
		fmt.Println("load plugin conf file failed, cannot open file")
	}
	defer f.Close()

	tmpConf := make(map[string][]confNode)

	// read config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&tmpConf)
	if err != nil {
		fmt.Println("load plugin conf file failed while parsing")
		return err
	}

	for pname, pconfs := range tmpConf {
		if pm.pluginConf[pname] == nil {
			pm.pluginConf[pname] = make(map[string]confNode)
		}

		for _, pconf := range pconfs {
			pm.pluginConf[pname][pconf.SubType] = pconf
		}
	}

	pluginCount := len(pm.pluginConf)
	fmt.Println("Plugin conf number is:", pluginCount, "data=", pm.pluginConf)
	return nil
}

func (pm *PluginMgr) loadOnePlugin(name string, subtype string) (pi interface{}, err error) {
	// open plugin
	conf := pm.pluginConf[name][subtype]
	pg, err := plugin.Open(conf.Path)
	if err != nil {
		fmt.Println("Warn: plugin", name, "open failed, err=", err)
		err = errors.New("plugin open failed!")
	}

	// get instance
	iSymbol, err := pg.Lookup("GetInstance")
	if err != nil {
		fmt.Println("Warn: plugin", name, "don't have func named GetInstance, err=", err)
		err = errors.New("Invalid plugin, it doesn't meet our requirements")
	}
	pi = iSymbol.(func() interface{})()

	// TODO: verify the plugin's signature, make sure it's authorized by us

	return
}
