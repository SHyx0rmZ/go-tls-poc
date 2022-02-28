package tls_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	tls "github.com/SHyx0rmZ/go-tls-poc"
)

func TestStore(t *testing.T) {
	store := tls.Store{
		New: func() interface{} {
			return make([]int, 0)
		},
	}

	do := func(s *tls.Store, ch chan<- tls.CleanupFunc, wg *sync.WaitGroup) tls.CleanupFunc {
		if wg != nil {
			defer wg.Done()
		}
		is := s.Load().([]int)
		is = append(is, rand.Intn(10))
		is = append(is, rand.Intn(10))
		is = append(is, rand.Intn(10))
		for sum, i := 0, 0; i < len(is); i++ {
			sum += is[i]
			fmt.Println(is[i], sum)
		}
		s.Store(is)
		select {
		case ch <- s.Closer():
		default:
		}
		return s.Closer()
	}

	runtime.Gosched()

	ch := make(chan tls.CleanupFunc, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go do(&store, ch, &wg)
	go do(&store, ch, &wg)

	wg.Wait()

	do(&store, nil, nil)
	c := do(&store, nil, nil)

	fmt.Println(store)

	_ = c.Close()
	close(ch)

	fmt.Println(store)
	for c := range ch {
		c()

		fmt.Println(store)
	}
}
