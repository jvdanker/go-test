package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"os"
	"os/signal"
)

func main() {
	//input := "/Volumes/Juan/Photos/Diversen/"
	//output := "/Volumes/App/output/"

	var (
		input  = "input"
		output = "output"
	)

	flag.StringVar(&input, "i", input, "Input directory to process")
	flag.StringVar(&output, "o", output, "Output directory to write results to")
	flag.Parse()

	fmt.Printf("input=%v, output=%v\n", input, output)

	ctx, cancel := setupExitChannel()
	defer cancel()

	//os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(ctx, input)
	dirs = walker.CreateDirectories(output, dirs)
	tasks.ResizeImages(ctx, dirs, output)
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
