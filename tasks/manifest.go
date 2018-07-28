package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
)

func CreateManifest(in <-chan []util.ProcessedImage) {
	fmt.Printf("Create manifest files\n")

	for processedImages := range in {
		fmt.Println(processedImages)
		return
		inputdir := ""
		output := ""
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
