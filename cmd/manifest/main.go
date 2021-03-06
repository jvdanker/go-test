package main

import (
	"flag"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
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

	ctx, cancel := util.SetupExitChannel()
	defer cancel()

	os.RemoveAll(output + "/manifest")
	os.MkdirAll(output+"/manifest", os.ModePerm)

	dirs := walker.WalkDirectories(ctx, input)

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			tasks.CreateManifests(ctx, dirs, output)
			wg.Done()
		}()
	}

	wg.Wait()
}
