package nopreempt

import (
	"bytes"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestGoId(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// t.Log("slow", goid())
			// t.Log("quick", GetGID())
			if goid() != GetGID() {
				t.Fail()
			}
		}()
	}
	wg.Wait()
}

func TestMId(t *testing.T) {
	const mp = 1
	omp := runtime.GOMAXPROCS(mp)
	defer runtime.GOMAXPROCS(omp)
	set := map[int64]struct{}{}
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			set[GetMID()] = struct{}{}
		}()
	}
	wg.Wait()
	t.Log(set)
	if len(set) != mp {
		t.Fail()
	}
}

func TestMId2(t *testing.T) {
	const mp = 2
	omp := runtime.GOMAXPROCS(mp)
	defer runtime.GOMAXPROCS(omp)
	set := map[int64]struct{}{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			set[GetMID()] = struct{}{}
		}()
	}
	wg.Wait()
	t.Log(set)
	if len(set) > mp {
		t.Fail()
	}
}

func TestMId3(t *testing.T) {
	const mp = 3
	omp := runtime.GOMAXPROCS(mp)
	defer runtime.GOMAXPROCS(omp)
	t.Log("old GOMAXPROCS", omp)
	set := map[int64]struct{}{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < omp*3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			set[GetMID()] = struct{}{}
			mu.Unlock()

			var s float64
			for k := 0; k < 1000000; k++ {
				s *= rand.Float64()
			}
		}()
	}
	wg.Wait()
	t.Log(set)
	if len(set) > mp {
		t.Fail()
	}
}

func TestMId10(t *testing.T) {
	const mp = 10
	omp := runtime.GOMAXPROCS(mp)
	defer runtime.GOMAXPROCS(omp)
	t.Log("old GOMAXPROCS", omp)
	set := map[int64]struct{}{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < omp*3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			set[GetMID()] = struct{}{}
			mu.Unlock()

			var s float64
			for k := 0; k < 1000000; k++ {
				s *= rand.Float64()
			}
		}()
	}
	wg.Wait()
	t.Log(set)
	if len(set) > mp {
		t.Fail()
	}
}

func TestNoPreempt(t *testing.T) {
	omp := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(omp)

	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now()
	go func() {
		for time.Since(start) < 2*time.Second {
			t.Log("busy")
			var s float64
			for k := 0; k < 100000000; k++ {
				s *= rand.Float64()
			}
		}
		wg.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	t.Log("exit")
	if time.Since(start) < time.Second {
		t.Log("preempted")
	} else {
		t.Error("cannot preempt")
		t.Fail()
	}
	wg.Wait()
}

func TestPreempt(t *testing.T) {
	omp := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(omp)

	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now()
	go func() {
		mp := AcquireM()
		defer mp.Release()
		for time.Since(start) < 2*time.Second {
			t.Log("busy")
			var s float64
			for k := 0; k < 100000000; k++ {
				s *= rand.Float64()
			}
		}
		wg.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	t.Log("exit")
	if time.Since(start) < 2*time.Second {
		t.Error("preempted")
		t.Fail()
	} else {
		t.Log("cannot preempt")
	}
	wg.Wait()
}

func goid() int64 {
	buf := make([]byte, len("goroutine ddddddddd"))
	runtime.Stack(buf, false)
	buf = buf[len("goroutine "):]
	buf = buf[:bytes.IndexByte(buf, ' ')]
	gid, _ := strconv.ParseInt(string(buf), 10, 64)
	return gid
}
