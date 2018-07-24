package pipes

import (
	"sync"
)

func Merge(cs ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	out := make(chan interface{})

	digest := func(c <-chan interface{}) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go digest(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
