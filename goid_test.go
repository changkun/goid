package goid_test

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/changkun/goid"
)

func slowGetGoID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseInt(idField, 10, 64) // very unlikely to be failed
	return id
}

func ExampleGet() {
	cnums := make(chan int, 100)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			cnums <- int(goid.Get()) // down cast, wrong in large goid
			wg.Done()
		}()
	}
	wg.Wait()
	close(cnums)

	nums := []int{int(goid.Get())}
	for v := range cnums {
		nums = append(nums, v)
	}
	sort.Ints(nums)
	fmt.Printf("%v", nums)
}

func TestGet(t *testing.T) {
	got := goid.Get()
	want := slowGetGoID()
	if got != uint64(want) {
		t.Errorf("want %d, got: %d", want, got)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = goid.Get()
	}
}
