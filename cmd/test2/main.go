package main

import (
	"fmt"
	"github.com/jvdanker/go-test/tasks"
)

func main() {
	input := "images3/"
	output := "output/"

	fmt.Printf("input=%v, output=%v\n", input, output)

	//os.RemoveAll(output)
	//os.MkdirAll(output, os.ModePerm)

	tasks.ResizeImages(input, "output")
	//tasks.SliceImages("output/images3/", "output/slices/")
	//tasks.CreateBottomLayer("output/images3/", "output/slices/", "output/layers/")
	//tasks.CreateZoomLayers("output/layers/")
}
