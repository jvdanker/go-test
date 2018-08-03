package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func random(quit <-chan bool) <-chan int {
	out := make(chan int)

	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			select {
			case out <- r.Int():
				// do nothing
			case <-quit:
				fmt.Println("Stopping output function...")
				time.Sleep(3 * time.Second)
				close(out)
				fmt.Println("Stopped output function")
				return
			}
		}
	}()

	return out
}
func main() {
	quit := make(chan bool)

	// stop after pressing ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		fmt.Println("Press ctrl+c to interrupt...")
		<-c
		fmt.Println("Stopping...")
		quit <- true
	}()

	// stop after 30 seconds
	t := time.After(30 * time.Second)
	go func() {
		<-t
		quit <- true
	}()

	in := random(quit)
	for range in {
		// fmt.Println(i)
		time.Sleep(100 * time.Millisecond)
	}
}
