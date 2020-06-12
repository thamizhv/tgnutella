package cache

import (
	"sync"
)

type findFilesList struct {
	lock      sync.RWMutex
	findFiles map[string][]byte
}

func NewFindFilesList() ServentCache {
	return &findFilesList{
		lock:      sync.RWMutex{},
		findFiles: make(map[string][]byte),
	}
}

func (ff *findFilesList) Set(key string, value []byte) {
	ff.lock.Lock()
	defer ff.lock.Unlock()
	ff.findFiles[key] = value
}

func (ff *findFilesList) Remove(key string) {
	ff.lock.Lock()
	defer ff.lock.Unlock()
	delete(ff.findFiles, key)
}

func (ff *findFilesList) Exists(key string) bool {
	ff.lock.RLock()
	defer ff.lock.RUnlock()
	_, ok := ff.findFiles[key]
	return ok
}

func (ff *findFilesList) Get(key string) []byte {
	ff.lock.RLock()
	defer ff.lock.RUnlock()
	val, ok := ff.findFiles[key]
	if !ok {
		return nil
	}
	return val
}

func (ff *findFilesList) GetAll() map[string][]byte {
	ff.lock.RLock()
	defer ff.lock.RUnlock()

	out := make(map[string][]byte)

	for k, v := range ff.findFiles {
		out[k] = v
	}

	return out
}
