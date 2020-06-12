package cache

import (
	"sync"
)

type descriptorList struct {
	lock        sync.RWMutex
	descriptors map[string][]byte
}

func NewDescriptorList() ServentCache {
	return &descriptorList{
		lock:        sync.RWMutex{},
		descriptors: make(map[string][]byte),
	}
}

func (d *descriptorList) Set(key string, value []byte) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.descriptors[key] = value
}

func (d *descriptorList) Remove(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.descriptors, key)
}

func (d *descriptorList) Exists(key string) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	_, ok := d.descriptors[key]
	return ok
}

func (d *descriptorList) Get(key string) []byte {
	d.lock.RLock()
	defer d.lock.RUnlock()
	val, ok := d.descriptors[key]
	if !ok {
		return nil
	}
	return val
}

func (d *descriptorList) GetAll() map[string][]byte {
	d.lock.RLock()
	defer d.lock.RUnlock()

	out := make(map[string][]byte)

	for k, v := range d.descriptors {
		out[k] = v
	}

	return out
}
