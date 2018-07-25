package pipes

import (
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

	for i := 0; i < count; i++ {
		out := make(chan Container)
		result = append(result, out)

		go func(out chan Container, id int) {
			for n := range in {
				w(&n, i)
				out <- n
			}
			close(out)
		}(out, i)
	}

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
