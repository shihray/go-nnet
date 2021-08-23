package cache

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

// InMemory 記憶體緩存
func InMemory() Contract {
	return &MemoryMap{
		lock:   sync.Mutex{},
		cached: make(map[string]string),
	}
}

type MemoryMap struct {
	lock   sync.Mutex
	cached map[string]string
}

func (i *MemoryMap) GetMarshal(key string, unMarshal interface{}) error {
	cached, err := i.GetOrErr(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cached), unMarshal)
}

func (i *MemoryMap) SetMarshal(key string, canMarshalVal interface{}, seconds int) error {
	bytes, err := json.Marshal(canMarshalVal)
	if err != nil {
		return err
	}
	i.Set(key, string(bytes), seconds)
	return nil
}

func (i *MemoryMap) Exist(key string) bool {
	_, ok := i.cached[key]
	return ok
}

func (i *MemoryMap) GetOrErr(key string) (string, error) {
	v, ok := i.cached[key]
	if !ok {
		return "", errors.New("data not found")
	}
	return v, nil
}

func (i *MemoryMap) Get(key string, fallback string) string {
	v, err := i.GetOrErr(key)
	if err != nil {
		return fallback
	}
	return v
}

func (i *MemoryMap) Set(key string, value string, seconds int) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.cached[key] = value
	go func() {
		<-time.After(time.Duration(seconds) * time.Second)
		i.Remove(key)
	}()
}

func (i *MemoryMap) Remove(key string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	delete(i.cached, key)
	return nil
}

func (i *MemoryMap) RemovePrefix(prefix string) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	for k := range i.cached {
		if strings.HasPrefix(k, prefix) {
			delete(i.cached, k)
		}
	}
	return nil
}

func (i *MemoryMap) Reset() {
	i.cached = make(map[string]string)
	return
}

func (i *MemoryMap) SetIfNotExist(key string, value string, seconds int) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	if i.Exist(key) {
		return false
	}

	// set val
	i.cached[key] = value
	go func() {
		<-time.After(time.Duration(seconds) * time.Second)
		i.Remove(key)
	}()
	return true
}

func (i *MemoryMap) Incr(key string) (result int64, err error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	v, err := strconv.Atoi(i.cached[key])
	if err != nil {
		return -1, err
	}
	i.cached[key] = strconv.Itoa(v + 1)
	return int64(v + 1), nil
}

func (i *MemoryMap) GetInt64(key string, fallback int64) int64 {
	val := i.Get(key, strconv.FormatInt(fallback, 10))
	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return result
}