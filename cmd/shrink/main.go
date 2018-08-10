package main

import (
	"flag"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	//input := "/Volumes/Juan/Photos/Diversen/"
	//output := "/Volumes/App/output/"

	var (
		input   = ""
		output  = ""
		workers = 1
	)

	flag.StringVar(&input, "i", input, "Input directory to process")
	flag.StringVar(&output, "o", output, "Output directory to write results to")
	flag.IntVar(&workers, "w", workers, "Number of concurrent workers")
	flag.Parse()

	input, output = checkAndSanitizeArgs(input, output)

	fmt.Printf("input=%v, output=%v, workers=%v\n", input, output, workers)

	ctx, cancel := util.SetupExitChannel()
	defer cancel()

	//os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(ctx, input)
	dirs = walker.CreateDirectories(output, dirs)

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() { tasks.ResizeImages(ctx, dirs, output); wg.Done() }()
	}

	wg.Wait()
}

func checkAndSanitizeArgs(input, output string) (string, string) {
	if input == "" || output == "" {
		flag.Usage()
		os.Exit(1)
	}

	exit := false

	input = filepath.Clean(input)
	if _, err := os.Stat(input); err != nil {
		fmt.Fprintf(os.Stderr, "Input directory doesn't exists. Directory = %v\n", input)
		exit = true
	}

	output = filepath.Clean(output)
	if _, err := os.Stat(output); err != nil {
		fmt.Fprintf(os.Stderr, "Output directory doesn't exists. Directory = %v\n", output)
		exit = true
	}

	if exit {
		os.Exit(1)
	}

	return input, output
}
