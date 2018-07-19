package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Pipelines Test")

	c := numsToChan(1, 2, 3, 4, 5)

	cs := fanout(c, 2)

	for n := range merge(cs...) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}

func fanout(c <-chan int, count int) []<-chan int {
	result := make([]<-chan int, 0)

	for i := 0; i < count; i++ {
		result = append(result, readChan(c))
	}

	return result
}

func numsToChan(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, j := range nums {
			out <- j
		}

		close(out)
	}()

	return out
}

func readChan(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n
		}
		close(out)
	}()

	return out
}

func merge(cs ...<-chan int) <-chan int {
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
