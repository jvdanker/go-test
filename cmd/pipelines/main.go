package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

func main() {
	//fmt.Println("Pipelines Test")

	var data []int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 20; i++ {
		data = append(data, r.Int())
	}

	c := numsToChan(data...)
	cs := createWorkers(c, 2)
	merged := mergeWorkers(cs...)
	for m := range merged {
		//fmt.Printf("merged output n=%v\n", n)
		a := m
		m = a
	}

	fmt.Printf("\nmin=%v, max=%v\n", min, max)
}

var min int64 = math.MaxInt64
var max int64

func timings(f string, start int64) {
	end := time.Now().UnixNano()
	if start < min {
		min = start
	}
	if end > max {
		max = end
	}

	fmt.Printf("%v, %v, %v\n", f, start, end)
}

func numsToChan(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		start := time.Now().UnixNano()
		for _, j := range nums {
			//fmt.Printf("output to channel, n=%v\n", j)
			out <- j
		}

		close(out)
		timings("numsToChan", start)
	}()

	return out
}

func createWorkers(c <-chan int, count int) []<-chan int {
	result := make([]<-chan int, 0)

	for i := 0; i < count; i++ {
		//fmt.Printf("create worker %v\n", i)
		result = append(result, channelWorker(c, i))
	}

	return result
}

func channelWorker(in <-chan int, id int) <-chan int {
	out := make(chan int)

	go func(id int) {
		start := time.Now().UnixNano()
		for n := range in {
			start2 := time.Now().UnixNano()
			time.Sleep(1 * time.Millisecond)
			//fmt.Printf("worker=%v, n=%v\n", id, n)
			out <- n
			timings(fmt.Sprintf("channelWorker--(%v)", id), start2)
		}
		close(out)
		timings(fmt.Sprintf("channelWorker(%v)", id), start)
	}(id)

	return out
}

func mergeWorkers(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	digest := func(c <-chan int) {
		start := time.Now().UnixNano()
		for n := range c {
			out <- n
		}
		wg.Done()
		timings("mergeWorkers", start)
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
