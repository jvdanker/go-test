package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"strings"
)

func ResizeImages(quit <-chan bool, dirs <-chan string, output string) <-chan util.ProcessedDirectory {
	out := make(chan util.ProcessedDirectory)

	imagesOutput := strings.TrimSuffix(output, "/") + "/images"

	go func() {
		for inputdir := range dirs {
			fmt.Printf("ResizeImages, dirWorker=%v\n", inputdir)

			if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", imagesOutput, inputdir)); err == nil {
				continue
			}

			files := walker.WalkFiles(inputdir)

			var processedDirectory = util.Create(inputdir, output, imagesOutput+"/"+inputdir)
			for file := range files {
				select {
				case <-quit:
					//fmt.Println("Aborting ResizeImages...")
					close(out)
					return
				default:
					// do nothing
				}

				fmt.Printf("ResizeImages, filesWorkers=%v\n", file.Name)
				pi, err := util.ResizeFile(file, processedDirectory.OutputDir)
				if err != nil {
					panic(err)
				}

				processedDirectory.ProcessedImages = append(processedDirectory.ProcessedImages, pi)
			}

			out <- processedDirectory
		}

		close(out)
	}()

	return out
}
