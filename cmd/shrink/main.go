package main

import (
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

	quit := setupExitChannel()

	os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(quit, input)
	dirs = walker.CreateDirectories(output, dirs)
	pi := tasks.ResizeImages(dirs, output)
	for range tasks.CreateManifest(pi) {

	}
}

func setupExitChannel() <-chan bool {
	quit := make(chan bool)

	// stop after pressing ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		fmt.Println("Press ctrl+c to interrupt...")
		<-c
		fmt.Println("Shutting down...")
		quit <- true
	}()

	return quit
}
