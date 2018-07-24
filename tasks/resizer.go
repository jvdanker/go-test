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

	for w := 0; w < 1; w++ {
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
		w := fanoutFilesWorker(files, 2, worker, outputdir)
		processedImages := mergeFilesWorkers(w...)
		createManifestOfProcessedFiles(processedImages, worker, inputdir, outputdir)
	}
}

func createManifestOfProcessedFiles(processedImages <-chan util.ProcessedImage, worker int, inputdir string, outputdir string) {
	var processedFiles []util.ProcessedImage

	for file := range processedImages {
		processedFiles = append(processedFiles, file)
	}

	if len(processedFiles) > 0 {
		// create manifest file
		fmt.Printf("dirWorker=%v: createManifest=%v, count=%v\n", worker, inputdir, len(processedFiles))
		manifest.Create(processedFiles, inputdir, outputdir)
	}
}

func mergeFilesWorkers(cs ...<-chan util.ProcessedImage) <-chan util.ProcessedImage {
	var wg sync.WaitGroup
	out := make(chan util.ProcessedImage)

	digest := func(c <-chan util.ProcessedImage) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go digest(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func fanoutFilesWorker(in <-chan util.File, count int, dirWorkerId int, outputDir string) []<-chan util.ProcessedImage {
	result := make([]<-chan util.ProcessedImage, 0)

	for i := 0; i < count; i++ {
		result = append(result, filesWorker(dirWorkerId, i, outputDir, in))
	}

	return result
}

func filesWorker(dirWorkerId, fileWorkerId int, outputDir string, files <-chan util.File) <-chan util.ProcessedImage {
	out := make(chan util.ProcessedImage)

	go func(dirWorkerId, fileWorkerId int) {
		for file := range files {
			fmt.Printf("dirWorker=%v, fileWorker=%v: filesWorkers=%v\n", dirWorkerId, fileWorkerId, file.Name)
			file2 := util.ResizeFile(file, outputDir)
			out <- file2
		}
		close(out)
	}(dirWorkerId, fileWorkerId)

	return out
}
