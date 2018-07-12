package main

import (
	"github.com/jvdanker/go-test/tasks"
)

func main() {
	tasks.ResizeImages()
	tasks.SliceImages()
	tasks.CreateBottomLayer()
	tasks.CreateZoomLayers()
}
