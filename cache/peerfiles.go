package cache

import (
	"sync"
)

type peerFilesList struct {
	lock  sync.RWMutex
	peers map[string][]byte
}

func NewPeerFilesList() ServentCache {
	return &peerFilesList{
		lock:  sync.RWMutex{},
		peers: make(map[string][]byte),
	}
}

func (pf *peerFilesList) Set(key string, value []byte) {
	pf.lock.Lock()
	defer pf.lock.Unlock()
	pf.peers[key] = value
}

func (pf *peerFilesList) Remove(key string) {
	pf.lock.Lock()
	defer pf.lock.Unlock()
	delete(pf.peers, key)
}

func (pf *peerFilesList) Exists(key string) bool {
	pf.lock.RLock()
	defer pf.lock.RUnlock()
	_, ok := pf.peers[key]
	return ok
}

func (pf *peerFilesList) Get(key string) []byte {
	pf.lock.RLock()
	defer pf.lock.RUnlock()
	val, ok := pf.peers[key]
	if !ok {
		return nil
	}
	return val
}

func (pf *peerFilesList) GetAll() map[string][]byte {
	pf.lock.RLock()
	defer pf.lock.RUnlock()

	out := make(map[string][]byte)

	for k, v := range pf.peers {
		out[k] = v
	}

	return out
}
