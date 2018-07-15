package main

import "github.com/jvdanker/go-test/tasks"

func main() {
	//os.RemoveAll("./output/")
	//os.MkdirAll("./output", os.ModePerm)

	//tasks.ResizeImages("images3/", "./output")
	//tasks.SliceImages("./output/images3/")
	//tasks.CreateBottomLayer("output/images3/")
	tasks.CreateZoomLayers("output/parts")
}
