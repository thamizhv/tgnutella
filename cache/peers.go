package cache

import (
	"sync"
)

type peerList struct {
	lock  sync.RWMutex
	peers map[string][]byte
}

func NewPeerList() ServentCache {
	return &peerList{
		lock:  sync.RWMutex{},
		peers: make(map[string][]byte),
	}
}

func (p *peerList) Set(key string, value []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.peers[key] = value
}

func (p *peerList) Remove(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.peers, key)
}

func (p *peerList) Exists(key string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.peers[key]
	return ok
}

func (p *peerList) Get(key string) []byte {
	p.lock.RLock()
	defer p.lock.RUnlock()
	val, ok := p.peers[key]
	if !ok {
		return nil
	}
	return val
}

func (p *peerList) GetAll() map[string][]byte {
	p.lock.RLock()
	defer p.lock.RUnlock()

	out := make(map[string][]byte)

	for k, v := range p.peers {
		out[k] = v
	}

	return out
}
