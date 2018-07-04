package main

import (
	"fmt"
	_"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
)

func main() {
	fmt.Println("test")

	//tasks.ResizeImages()
	//tasks.SliceImages()

	dirs := walker.WalkDirectories("output/images/")
    for dir := range dirs {
        fmt.Println(dir)
    }
}
