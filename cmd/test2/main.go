package main

import (
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"os"
)

func main() {
	//input := "/Volumes/Juan/Photos/Diversen/"
	//output := "/Volumes/App/output/"

	input := "test"
	output := "output"

	fmt.Printf("input=%v, output=%v\n", input, output)

	os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(output, dirs)
	pi := tasks.ResizeImages(dirs, output)
	manifests := tasks.CreateManifest(pi)
	mergedImages := tasks.MergeImages(manifests)
	tasks.SliceImages(mergedImages)

	//CreateBottomLayer(output+"/images/", output+"/slices/", output+"/layers/")
	//tasks.CreateZoomLayers(output + "/layers/")
}
