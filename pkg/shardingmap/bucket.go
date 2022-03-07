package shardingmap

import (
	"errors"
	"sync"

	"github.com/x-junkang/connected/pkg/ciface"
)

type Bucket struct {
	data []*Map
	cap  int
	mask uint64
}

func NewBucket(cap int) (*Bucket, error) {
	if !checkUin64IsPowerof2(uint64(cap)) {
		return nil, errors.New("cap must 2^n")
	}
	mask := uint64(cap - 1)
	b := &Bucket{
		data: make([]*Map, cap),
		cap:  cap,
		mask: mask,
	}
	for i := 0; i < cap; i++ {
		b.data[i] = NewMap()
	}
	return b, nil
}

func (b *Bucket) Set(key uint64, value ciface.IConnection) {
	index := key & b.mask
	b.data[index].Set(key, value)
}

func (b *Bucket) Get(key uint64) (ciface.IConnection, bool) {
	index := key & b.mask
	return b.data[index].Get(key)
}

type Map struct {
	sync.RWMutex
	data map[uint64]ciface.IConnection
}

func NewMap() *Map {
	return &Map{
		data: make(map[uint64]ciface.IConnection),
	}
}

func (m *Map) Set(key uint64, value ciface.IConnection) {
	m.Lock()
	m.data[key] = value
	m.Unlock()
}

func (m *Map) Get(key uint64) (ciface.IConnection, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func checkUin64IsPowerof2(n uint64) bool {
	cnt := 0
	for i := 0; i < 64; i++ {
		if n&(1<<i) != 0 {
			cnt++
		}
	}
	return cnt == 1
}
