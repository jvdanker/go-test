package main

import (
	"fmt"
	_"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"github.com/jvdanker/go-test/manifest"
)

func main() {
	fmt.Println("test")

	//tasks.ResizeImages()
	//tasks.SliceImages()

    var total int
	dirs := walker.WalkDirectories("output/images/")
    for dir := range dirs {
        fmt.Println(dir)

        m, err := manifest.Read(dir + "/manifest.json")
        if (err == nil) {
            total += len(m.Files)
        }
    }
    fmt.Printf("Total %d\n", total)
}

// maxzoom = Math.ceil( Math.log( (cw/iw > ch/ih ? iw/cw : ih/ch) ) / Math.log(2) );