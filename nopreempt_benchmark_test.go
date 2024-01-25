package nopreempt

import (
	"sync"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		s += i // unix 0.3376 ns/op
	}
}

func BenchmarkAddLock(b *testing.B) {
	var s int
	var mu sync.Mutex
	for i := 0; i < b.N; i++ {
		mu.Lock()
		s += i // unix 5.359 ns/op
		mu.Unlock()
	}
}

func BenchmarkAddDisablePreempt(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		mp := AcquireM()
		s += i // unix 3.498 ns/op
		mp.Release()
	}
}
