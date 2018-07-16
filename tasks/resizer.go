package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"sync"
)

func ResizeImages(input, output string) {
	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(output, dirs)

	wg := sync.WaitGroup{}

	for w := 0; w < 5; w++ {
		wg.Add(1)
		go func(w int) {
			dirWorker(w, output, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func dirWorker(worker int, outputdir string, dirs <-chan string) {
	for inputdir := range dirs {
		fmt.Printf("dirWorker=%v: dirWorker=%v\n", worker, inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", outputdir, inputdir)); err == nil {
			// fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		files := walker.WalkFiles(inputdir)

		c := make(chan util.ProcessedImage)
		wg := sync.WaitGroup{}

		for w := 0; w < 3; w++ {
			wg.Add(1)
			go func(dirWorker, fileWorker int) {
				filesWorker(dirWorker, fileWorker, outputdir, files, c)
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

		if len(processedFiles) > 0 {
			// create manifest file
			fmt.Printf("dirWorker=%v: createManifest=%v\n", worker, inputdir)
			manifest.Create(processedFiles, inputdir, outputdir)
		}
	}
}

func filesWorker(dirWorker, fileWorker int, output string, files <-chan util.File, c chan<- util.ProcessedImage) {
	for file := range files {
		fmt.Printf("dirWorker=%v, fileWorker=%v: filesWorkers=%v\n", dirWorker, fileWorker, file.Name)
		file2 := util.ResizeFile(file, output)

		c <- file2
	}
}
