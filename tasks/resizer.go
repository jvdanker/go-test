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
			fmt.Printf("worker=%v: ResizeImages\n", w)
			dirWorker(w, output, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func dirWorker(worker int, output string, dirs <-chan string) {
	for dir := range dirs {
		fmt.Printf("worker=%v: dirWorker=%v\n", worker, dir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, dir)); err == nil {
			// fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		files := walker.WalkFiles(dir)
		images := filesWorker(worker, files)

		if len(images) == 0 {
			continue
		}

		// create manifest file
		fmt.Printf("worker=%v: createManifest=%v\n", worker, dir)
		m := manifest.Create(images, dir)

		// merge bottom layer images into one image
		mergeImages(worker, m)

		fmt.Println()
	}
}

func filesWorker(worker int, files <-chan util.File) []util.ProcessedImage {
	var result []util.ProcessedImage

	for file := range files {
		fmt.Printf("worker=%v: filesWorkers=%v\n", worker, file.Name)
		file2 := util.ResizeFile(file)
		result = append(result, file2)
	}

	return result
}

func mergeImages(worker int, m manifest.ManifestFile) {
	fmt.Printf("worker=%v: mergeImages\n", worker)

	bounds := m.Bounds()
	//fmt.Println(bounds)

	lm := layout.CreateBoxLayout()
	lm.Layout(bounds)
	m.Layout = lm

	image := merger.MergeImages(m)
	util.CreateImage(m.OutputDir+"/result.png", image)

	m.Update()
}
