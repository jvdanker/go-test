package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"strings"
)

func ResizeImages(dirs <-chan string, output string) <-chan []util.ProcessedImage {
	out := make(chan []util.ProcessedImage)

	output = strings.TrimSuffix(output, "/") + "/images"

	for inputdir := range dirs {
		fmt.Printf("dirWorker=%v\n", inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, inputdir)); err == nil {
			return nil
		}

		files := walker.WalkFiles(inputdir)

		var processedImages []util.ProcessedImage
		for file := range files {
			fmt.Printf("filesWorkers=%v\n", file.Name)
			pi := util.ResizeFile(file, output)
			processedImages = append(processedImages, pi)
		}

		go func() {
			out <- processedImages
		}()
	}

	return out
}
