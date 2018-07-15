package main

import (
	"github.com/jvdanker/go-test/tasks"
	"os"
)

func main() {
	os.RemoveAll("./output/")
	os.MkdirAll("./output", os.ModePerm)

	tasks.ResizeImages("images3/", "output")
	tasks.SliceImages("output/images3/", "output/slices/")
	tasks.CreateBottomLayer("output/images3/", "output/slices/", "output/layers/")
	tasks.CreateZoomLayers("output/layers/")
}
