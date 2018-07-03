package main

import (
	"fmt"
	"sync"
	"image"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/merger"
	"github.com/jvdanker/go-test/walker"
)

func filesWorker(id int, files <-chan util.File) []util.ProcessedImage {
    var result = []util.ProcessedImage{}

	for file := range files {
        file2 := util.ResizeFile(file)
        result = append(result, file2)
	}

	return result
}

func dirWorker(id int, dirs <-chan string) {
	for dir := range dirs {
		files := walker.WalkFiles(dir)
        images := filesWorker(id, files)

        manifest := manifest.Create(images, dir)

        bounds := manifest.Bounds()
        fmt.Println(bounds)

        lm := layout.CreateBoxLayout()
        lm.Layout(bounds)
        manifest.Layout = lm

        image := merger.MergeImages(manifest)
        util.CreateImage(manifest.OutputDir + "/result.png", image)

        manifest.Update()

        fmt.Println()
	}
}

func resizeImages() {
    dirs := walker.WalkDirectories("images/")
	dirs = walker.CreateDirectories(dirs)

	wg := sync.WaitGroup{}

	for w := 0; w < 1; w++ {
		wg.Add(1)
		go func(w int) {
			dirWorker(w, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func sliceImages() {
    dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
	    fmt.Println(dir)

	    img := util.DecodeImage(dir + "/result.png")
	    //var x,y int
	    bounds := img.Bounds()
        w, h := bounds.Max.X, bounds.Max.Y
        fmt.Println("w,h=", w, h)

        var x,y,i,j int
        for y<h {
            for x=0;  x<w; x+=256 {
                r := image.Rect(x, y, x + 256, y + 256)
                fmt.Println(r)

                a := image.RGBA(img)
                sub := img.SubImage(r)
                util.CreateImage(dir + "/sub-" + string(i) + "-" + string(j) + ".png", sub)

                i++
            }

            y += 256
            j++
        }

        fmt.Println()
	}
}

func main() {
	fmt.Println("test")

	//resizeImages()
	sliceImages()
}
