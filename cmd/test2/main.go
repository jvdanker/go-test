package main

import (
	"fmt"
	"github.com/jvdanker/go-test/tasks"
)

func main() {
	input := "images/"
	output := "output/"

	fmt.Printf("input=%v, output=%v\n", input, output)

	//os.RemoveAll(output)
	//os.MkdirAll(output, os.ModePerm)

	tasks.ResizeImages(input, "output")
	tasks.SliceImages("output/images/", "output/slices/")
	tasks.CreateBottomLayer("output/images/", "output/slices/", "output/layers/")
	tasks.CreateZoomLayers("output/layers/")
}
