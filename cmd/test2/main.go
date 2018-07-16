package main

import (
	"fmt"
	"github.com/jvdanker/go-test/tasks"
)

func main() {
	//input := "/Volumes/Juan/Photos/Diversen/"
	//output := "/Volumes/App/output/"

	input := "Diversen"
	output := "output"

	fmt.Printf("input=%v, output=%v\n", input, output)

	//os.RemoveAll(output)
	//os.MkdirAll(output, os.ModePerm)

	//tasks.ResizeImages(input, output+"/images")
	// TODO move merge images to task
	tasks.SliceImages(output+"/images/", output+"/slices/")
	tasks.CreateBottomLayer(output+"/images/", output+"/slices/", output+"/layers/")
	tasks.CreateZoomLayers(output + "/layers/")
}
