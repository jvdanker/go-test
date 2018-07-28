package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
)

// TODO return channel of slice of processedImages to create manifest
func ResizeImages(dirs <-chan string, output string) {
	for inputdir := range dirs {
		fmt.Printf("dirWorker=%v\n", inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, inputdir)); err == nil {
			return
		}

		files := walker.WalkFiles(inputdir)

		var processedImages []util.ProcessedImage
		for file := range files {
			fmt.Printf("filesWorkers=%v\n", file.Name)
			pi := util.ResizeFile(file, output)
			processedImages = append(processedImages, pi)
		}

		createManifestOfProcessedFiles(processedImages, inputdir, output)
	}
}

func createManifestOfProcessedFiles(processedImages []util.ProcessedImage, inputdir string, outputdir string) {
	if len(processedImages) > 0 {
		// create manifest file
		fmt.Printf("dirWorker=%v: count=%v\n", inputdir, len(processedImages))
		manifest.Create(processedImages, inputdir, outputdir)
	}
}
