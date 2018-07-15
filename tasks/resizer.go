package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/merger"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"sync"
)

func ResizeImages(input, output string) {
	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(dirs)

	wg := sync.WaitGroup{}

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go func(w int) {
			fmt.Printf("ResizeImages, worker=%v\n", w)
			dirWorker(w, output, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func dirWorker(worker int, output string, dirs <-chan string) {
	for dir := range dirs {
		fmt.Printf("Processing dir: worker=%v, dir=%v\n", worker, dir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, dir)); err == nil {
			// fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		files := walker.WalkFiles(dir)
		images := filesWorker(files)

		if len(images) == 0 {
			continue
		}

		createManifest(images, dir)

		fmt.Println()
	}
}

func filesWorker(files <-chan util.File) []util.ProcessedImage {
	var result []util.ProcessedImage

	for file := range files {
		fmt.Printf("Resize file %v\n", file.Name)
		file2 := util.ResizeFile(file)
		result = append(result, file2)
	}

	return result
}

func createManifest(images []util.ProcessedImage, dir string) {
	m := manifest.Create(images, dir)

	bounds := m.Bounds()
	fmt.Println(bounds)

	lm := layout.CreateBoxLayout()
	lm.Layout(bounds)
	m.Layout = lm

	image := merger.MergeImages(m)
	util.CreateImage(m.OutputDir+"/result.png", image)

	m.Update()
}
