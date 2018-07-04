package main

import (
	"fmt"
	"sync"
	"image"
	"os"
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

        if len(images) == 0 {
            return
        }

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

	    if _, err := os.Stat(dir + "/result.png"); err != nil {
	        continue
	    }

	    img := util.DecodeImage(dir + "/result.png")
	    //var x,y int
	    bounds := img.Bounds()
        w, h := bounds.Max.X, bounds.Max.Y
        fmt.Println("w,h=", w, h)

        var x,y,i,j int
        i = 0
        j = 0

        for y<h {
            for x=0;  x<w; x+=256 {
                r := image.Rect(x, y, x + 256, y + 256)
                fmt.Println(r)

                //bounds := img.Bounds()
                nrgba := img.(*image.NRGBA)
                //image.NewRGBA(image.Rect(0, 0, 256, 256))
                sub := nrgba.SubImage(r)
                //sub = nrgba

                util.CreateImage(fmt.Sprintf("%s/sub-%d-%d.png", dir, i, j), sub)

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
