package main

import (
	"github.com/jvdanker/go-test/tasks"
)

func main() {
	tasks.ResizeImages("images/", "./output")
	//tasks.SliceImages()
	//tasks.CreateBottomLayer()
	//tasks.CreateZoomLayers("output/parts")
}
