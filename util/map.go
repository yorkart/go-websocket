package util

import (
	"sync"
)

// Key for map key
type Key interface{}

// Val for map value
type Val interface{}

// Map custom collection
type Map struct {
	m    map[Key]Val
	lock sync.RWMutex
}

// NewMap create Map instance
func NewMap() *Map {
	return &Map{
		m: map[Key]Val{},
	}
}

// Put put key,value
func (p *Map) Put(key Key, val Val) {
	p.lock.Lock()
	p.m[key] = val
	p.lock.Unlock()
}

// Get get value by key
func (p *Map) Get(key Key) (Val, bool) {
	p.lock.RLock()
	val, ok := p.m[key]
	p.lock.RUnlock()

	return val, ok
}

func (p *Map) Del(key Key) {
	p.lock.RLock()
	delete(p.m, key)
	p.lock.RUnlock()
}

// Len return map length
func (p *Map) Len() int {
	p.lock.RLock()
	var l = len(p.m)
	p.lock.RUnlock()

	return l
}

// Range just range map
func (p *Map) Range(action func(key Key, val Val)) {
	defer p.lock.RUnlock()
	p.lock.RLock()

	for key, val := range p.m {
		action(key, val)
	}
}
