package util

import (
	"sync"
)

// Set custom collection
type Set struct {
	m     map[interface{}]bool
	cache []interface{}
	lock  sync.RWMutex
}

// NewSet create Set instance
func NewSet() *Set {
	return &Set{
		m:     map[interface{}]bool{},
		cache: []interface{}{},
	}
}

// Add add element to set
func (p *Set) Add(item interface{}) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.m[item] {
		p.m[item] = true
		p.list()
		return true
	}
	return false
}

// Remove remove a element from set
func (p *Set) Remove(item interface{}) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.m[item] {
		delete(p.m, item)
		p.list()
		return true
	}
	return false
}

// Clear clear set
func (p *Set) Clear() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.m = map[interface{}]bool{}
	p.list()
}

// Get get a element by a func
func (p *Set) Get(fn func(value interface{}) bool) (interface{}, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for key := range p.m {
		if fn(key) {
			return key, true
		}
	}
	return nil, false
}

// Len get the set len
func (p *Set) Len() int {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return len(p.m)
}

func (p *Set) list() {
	list := []interface{}{}
	for key := range p.m {
		list = append(list, key)
	}
	p.cache = list
}

// List return the set as list
func (p *Set) List() []interface{} {
	return p.cache
}
