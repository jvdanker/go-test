package pipes

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
	"sync"
	"time"
)

func Fanout2(in <-chan Container, count int, w Worker) []<-chan Container {
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
			start := time.Now().UnixNano()

			var count int

			for n := range in {
				startWorker := time.Now().UnixNano()

				count++
				w(&n, id)
				out <- n

				util.Timings(fmt.Sprintf("Fanout - Worker(%v)", id), startWorker)
			}

			mutex.Lock()
			total += count
			stats[id] = count
			mutex.Unlock()

			close(out)
			wg.Done()

			util.Timings(fmt.Sprintf("Fanout(%v)", id), start)
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
