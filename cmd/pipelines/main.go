package main

import (
	"fmt"
	"github.com/jvdanker/go-test/pipes"
	"math/rand"
	"time"
)

func main() {
	//fmt.Println("Pipelines Test")

	var data []int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 10; i++ {
		data = append(data, r.Int())
	}

	c := numsToChan(data)
	merged := pipes.FanoutAndMerge(c, 1, func(c *pipes.Container, id int) {
		time.Sleep(100 * time.Millisecond)
		//fmt.Printf("worker=%v, n=%v\n", id, n)
	})

	for m := range merged {
		x := m.(int)
		fmt.Printf("merged output n=%v\n", x)
	}

	//fmt.Printf("\nmin=%v, max=%v\n", util.TMin, util.TMax)
}

func numsToChan(v []int) <-chan pipes.Container {
	out := make(chan pipes.Container)

	go func() {
		for _, j := range v {
			//fmt.Printf("output to channel, n=%v\n", j)
			out <- pipes.Container{Payload: j}
		}

		close(out)
	}()

	return out
}
