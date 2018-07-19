package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Pipelines Test")

	c := numsToChan(1, 2, 3, 4, 5)
	cs := createWorkers(c, 1)
	for n := range mergeWorkers(cs...) {
		fmt.Printf("merged output n=%v\n", n)
	}
}

func numsToChan(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, j := range nums {
			fmt.Printf("output to channel, n=%v\n", j)
			out <- j
		}

		close(out)
	}()

	return out
}

func createWorkers(c <-chan int, count int) []<-chan int {
	result := make([]<-chan int, 0)

	for i := 0; i < count; i++ {
		fmt.Printf("create worker %v\n", i)
		result = append(result, channelWorker(c, i))
	}

	return result
}

func channelWorker(in <-chan int, id int) <-chan int {
	out := make(chan int)

	go func(id int) {
		for n := range in {
			fmt.Printf("worker=%v, n=%v\n", id, n)
			out <- n
		}
		close(out)
	}(id)

	return out
}

func mergeWorkers(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
