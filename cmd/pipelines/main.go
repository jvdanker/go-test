package main

import (
	"fmt"
	"github.com/jvdanker/go-test/pipes"
	"math"
	"math/rand"
	"time"
)

func main() {
	//fmt.Println("Pipelines Test")

	var data []int
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 200000; i++ {
		data = append(data, r.Int())
	}

	c := numsToChan(data)
	merged := pipes.FanoutAndMerge(c, 200, channelWorker)

	for m := range merged {
		x := m.(int)
		fmt.Printf("merged output n=%v\n", x)
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

func numsToChan(v []int) <-chan pipes.Container {
	out := make(chan pipes.Container)

	go func() {
		start := time.Now().UnixNano()
		for _, j := range v {
			//fmt.Printf("output to channel, n=%v\n", j)
			out <- pipes.Container{Payload: j}
		}

		close(out)
		timings("numsToChan", start)
	}()

	return out
}

func channelWorker(c *pipes.Container, id int) {
	start := time.Now().UnixNano()
	time.Sleep(1 * time.Millisecond)
	//fmt.Printf("worker=%v, n=%v\n", id, n)
	timings(fmt.Sprintf("channelWorker(%v)", id), start)
}
