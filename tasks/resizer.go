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

	for w := 0; w < 100; w++ {
		wg.Add(1)
		go func(w int) {
			dirWorker(w, output, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func dirWorker(worker int, output string, dirs <-chan string) {
	for dir := range dirs {
		fmt.Printf("dirWorker=%v: dirWorker=%v\n", worker, dir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, dir)); err == nil {
			// fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		files := walker.WalkFiles(dir)

		c := make(chan util.ProcessedImage, 1)
		wg := sync.WaitGroup{}

		for w := 0; w < 100; w++ {
			wg.Add(1)
			go func(dirWorker, fileWorker int) {
				filesWorker(dirWorker, fileWorker, files, c)
				wg.Done()
			}(worker, w)
		}

		go func() {
			wg.Wait()
			close(c)
		}()

		var processedFiles []util.ProcessedImage
		for file := range c {
			processedFiles = append(processedFiles, file)
		}

		// create manifest file
		fmt.Printf("dirWorker=%v: createManifest=%v\n", worker, dir)
		m := manifest.Create(processedFiles, dir)

		// merge images into one image
		mergeImages(worker, m)
	}
}

func filesWorker(dirWorker, fileWorker int, files <-chan util.File, output chan<- util.ProcessedImage) {
	for file := range files {
		fmt.Printf("dirWorker=%v, fileWorker=%v: filesWorkers=%v\n", dirWorker, fileWorker, file.Name)
		file2 := util.ResizeFile(file)

		output <- file2
	}
}

func mergeImages(dirWorker int, m manifest.ManifestFile) {
	fmt.Printf("dirWorker=%v: mergeImages\n", dirWorker)

	bounds := m.Bounds()
	//fmt.Println(bounds)

	lm := layout.CreateBoxLayout()
	lm.Layout(bounds)
	m.Layout = lm

	image := merger.MergeImages(m)
	util.CreateImage(m.OutputDir+"/result.png", image)

	m.Update()
}
