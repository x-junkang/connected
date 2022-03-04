package shardingmap

import "sync"

type Bucket struct {
	data []*Map
	Cap  int
}

func NewBucket(cap int) *Bucket {
	b := &Bucket{
		data: make([]*Map, cap),
		Cap:  cap,
	}
	for i := 0; i < cap; i++ {
		b.data[i] = NewMap()
	}
	return b
}

func (b *Bucket) Set(key int, value string) {
	index := key % b.Cap
	b.data[index].Set(key, value)
}

func (b *Bucket) Get(key int) (string, bool) {
	index := key % b.Cap
	return b.data[index].Get(key)
}

type Map struct {
	sync.RWMutex
	data map[int]string
}

func NewMap() *Map {
	return &Map{
		data: map[int]string{},
	}
}

func (m *Map) Set(key int, value string) {
	m.Lock()
	m.data[key] = value
	m.Unlock()
}

func (m *Map) Get(key int) (string, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.data[key]
	return value, ok
}
