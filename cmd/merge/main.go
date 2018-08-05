package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"os"
	"os/signal"
	"sync"
)

func main() {
	var (
		input   = "input"
		output  = "output"
		workers = 1
	)

	flag.StringVar(&input, "i", input, "Input directory to process")
	flag.StringVar(&output, "o", output, "Output directory to write results to")
	flag.IntVar(&workers, "w", workers, "Number of concurrent workers")
	flag.Parse()

	fmt.Printf("input=%v, output=%v, workers=%v\n", input, output, workers)

	ctx, cancel := setupExitChannel()
	defer cancel()

	os.RemoveAll(output + "/merged")
	os.MkdirAll(output+"/merged", os.ModePerm)

	dirs := walker.WalkDirectories(ctx, output+"/manifest")

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			tasks.MergeImages(ctx, dirs, output)
			wg.Done()
		}()
	}

	wg.Wait()
}

func setupExitChannel() (context.Context, context.CancelFunc) {
	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// stop after pressing ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		fmt.Println("Press ctrl+c to interrupt...")
		<-c
		fmt.Println("Shutting down...")
		cancel()
	}()

	return ctx, cancel
}
