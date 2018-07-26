package pipes

import (
	"fmt"
	"sync"
)

type Container struct {
	Payload interface{}
}

type Worker func(*Container, int)

func FanoutAndMerge(in <-chan Container, count int, w Worker) <-chan interface{} {
	return Merge(Fanout(in, count, w))
}

func Fanout(in <-chan Container, count int, w Worker) []<-chan Container {
	result := make([]<-chan Container, 0)

	var wg sync.WaitGroup
	wg.Add(count)

	var mutex = &sync.Mutex{}
	var stats = make(map[int]int)
	var total int

	for i := 0; i < count; i++ {
		out := make(chan Container)
		result = append(result, out)

		go func(out chan Container, id int) {
			var count int

			for n := range in {
				count++
				w(&n, id)
				out <- n
			}

			mutex.Lock()
			total += count
			stats[id] = count
			mutex.Unlock()

			close(out)
			wg.Done()
		}(out, i)
	}

	go func() {
		wg.Wait()

		fmt.Printf("Total number of messages=%v\n", total)
		for k, v := range stats {
			if v > 0 {
				fmt.Printf("channel=%v, count=%v\n", k, v)
			}
		}
		fmt.Printf("stats=%v\n", stats)
	}()

	return result
}

func Merge(in []<-chan Container) <-chan interface{} {
	var wg sync.WaitGroup
	out := make(chan interface{})

	digest := func(c <-chan Container) {
		for n := range c {
			out <- n.Payload
		}
		wg.Done()
	}

	wg.Add(len(in))
	for _, c := range in {
		go digest(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
