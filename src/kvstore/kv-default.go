// This is a very ugly text-based kv store
// Do not support change existing keys
// But it works...
package kvstore

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type KVText struct {
	meta    Meta
	rwmutex sync.RWMutex
}

func (store *KVText) Init(conf string) {
	store.meta.DBName = conf + ".txt"
	store.meta.Capacity = 0
}

func (store *KVText) Get(key string) (value string, err error) {
	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	f, err := os.Open(store.meta.DBName)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var text string
	for scanner.Scan() {
		text = scanner.Text()
		kvpair := strings.Split(text, " ")
		if len(kvpair) != 2 || kvpair[0] != key {
			continue
		}
		value = kvpair[1]
		break
	}

	return
}

func (store *KVText) Set(key string, value string) (err error) {
	store.rwmutex.Lock()
	defer store.rwmutex.Unlock()

	f, err := os.OpenFile(store.meta.DBName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Open file failed")
		return
	}
	defer f.Close()
	writestr := key + " " + value + "\n"
	_, err = f.Write(([]byte)(writestr))
	if err != nil {
		fmt.Println("Set key failed")
		return
	}

	return
}

func (store *KVText) Count() (count int64) {
	count = 0

	store.rwmutex.RLock()
	defer store.rwmutex.RUnlock()

	f, err := os.Open(store.meta.DBName)
	if err != nil {
		return 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		_ = scanner.Text()
		count++
	}

	return count
}

func (store *KVText) GetMeta() (meta *Meta) {
	return &store.meta
}
