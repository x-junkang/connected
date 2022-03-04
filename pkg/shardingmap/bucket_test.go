package shardingmap

import (
	"sync"
	"testing"
)

const ConcurrentWrite = 100
const ConcurrentRead = 2000

func BenchmarkBucket(b *testing.B) {
	bucket := NewBucket(64)
	wg := sync.WaitGroup{}
	fn := func() {
		for i := 0; i < b.N; i++ {
			bucket.Set(i, "hello")
		}
		wg.Done()
	}
	fn2 := func() {
		for i := 0; i < b.N; i++ {
			bucket.Get(i)
		}
		wg.Done()
	}
	wg.Add(ConcurrentRead + ConcurrentWrite)

	loop := func(cnt int, fn func()) {
		for i := 0; i < cnt; i++ {
			go fn()
		}
	}
	loop(ConcurrentWrite, fn)
	loop(ConcurrentRead, fn2)

	wg.Wait()
}

func BenchmarkMap(b *testing.B) {
	bucket := NewMap()
	wg := sync.WaitGroup{}
	fn := func() {
		for i := 0; i < b.N; i++ {
			bucket.Set(i, "hello")
		}
		wg.Done()
	}
	fn2 := func() {
		for i := 0; i < b.N; i++ {
			bucket.Get(i)
		}
		wg.Done()
	}
	wg.Add(ConcurrentRead + ConcurrentWrite)

	loop := func(cnt int, fn func()) {
		for i := 0; i < cnt; i++ {
			go fn()
		}
	}
	loop(ConcurrentWrite, fn)
	loop(ConcurrentRead, fn2)
	wg.Wait()
}

func BenchmarkSTDMap(b *testing.B) {
	bucket := sync.Map{}
	wg := sync.WaitGroup{}
	fn := func() {
		for i := 0; i < b.N; i++ {
			bucket.Store(i, "hello")
		}
		wg.Done()
	}
	fn2 := func() {
		for i := 0; i < b.N; i++ {
			value, ok := bucket.Load(i)
			if ok {
				_ = value.(string)
			}
		}
		wg.Done()
	}
	wg.Add(ConcurrentRead + ConcurrentWrite)

	loop := func(cnt int, fn func()) {
		for i := 0; i < cnt; i++ {
			go fn()
		}
	}
	loop(ConcurrentWrite, fn)
	loop(ConcurrentRead, fn2)
	wg.Wait()
}

// read:write   2000:500
// pkg: github.com/x-junkang/connected/pkg/ciface/shardingmap
// cpu: AMD Ryzen 7 4800H with Radeon Graphics
// BenchmarkBucket-16    	   18242	     60305 ns/op	     118 B/op	       0 allocs/op
// BenchmarkMap-16       	    8276	    219452 ns/op	     133 B/op	       0 allocs/op
// BenchmarkSTDMap-16    	   23397	     59087 ns/op	   13911 B/op	     996 allocs/op

// read:write   2000:100
// cpu: AMD Ryzen 7 4800H with Radeon Graphics
// BenchmarkBucket-16    	   25700	     51920 ns/op	      85 B/op	       0 allocs/op
// BenchmarkMap-16       	    7929	    132563 ns/op	     138 B/op	       0 allocs/op
// BenchmarkSTDMap-16    	   37898	     46835 ns/op	    4818 B/op	     201 allocs/op
