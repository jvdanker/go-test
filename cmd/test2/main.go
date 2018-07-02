package main

import (
	"fmt"
	"os"
	"sync"
	"image/png"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/merger"
)

func mergeFiles(lm layout.LayoutManager) {
    image := merger.MergeImages(lm)

    outfilename := lm.OutputDir + "/result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
}

func filesWorker(id int, files <-chan util.File) []util.ProcessedImage {
    var result = []util.ProcessedImage{}

	for file := range files {
		// fmt.Println(id, files, file)
        file2 := util.ResizeFile(file)
        result = append(result, file2)
	}

	return result
}

func dirWorker(id int, dirs <-chan string) {
	for dir := range dirs {
		files := util.WalkFiles(dir)
        images := filesWorker(id, files)
        manifest := manifest.Create(images, dir)

        lm := layout.CreateBoxLayout(manifest)
        lm.Layout()

        mergeFiles(lm)
        fmt.Println()

        // update manifest
        lm.Update()
	}
}

func main() {
	fmt.Println("test")

	dirs := util.WalkDirectories()
	dirs = util.CreateDirectories(dirs)

	wg := sync.WaitGroup{}

	for w := 1; w <= 1; w++ {
		wg.Add(1)
		go func(w int) {
			dirWorker(w, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}
