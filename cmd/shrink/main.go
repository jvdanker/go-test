package main

import (
	"flag"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"os"
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

	os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(output, dirs)
	pi := tasks.ResizeImages(dirs, output)
	for range tasks.CreateManifest(pi) {

	}
}
