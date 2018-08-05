package tasks

import (
	"context"
	"fmt"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"strings"
)

func ResizeImages(ctx context.Context, dirs <-chan string, output string) {
	imagesOutput := strings.TrimSuffix(output, "/") + "/images"

	for inputdir := range dirs {
		fmt.Printf("ResizeImages, dirWorker=%v\n", inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", imagesOutput, inputdir)); err == nil {
			continue
		}

		files := walker.WalkFiles(inputdir)

		var processedDirectory = util.Create(inputdir, output, imagesOutput+"/"+inputdir)
		for file := range files {
			select {
			case <-ctx.Done():
				fmt.Println("Aborting ResizeImages...")
				return
			default:
				// do nothing
			}

			// newName := output + "/" + file.Name + ".png"
			dest := fmt.Sprintf("%v/%v.png", processedDirectory.OutputDir, file.Name)
			if _, err := os.Stat(dest); err == nil {
				//fmt.Printf("ResizeImages, skipping=%v\n", dest)
				continue
			}

			pi, err := util.ResizeFile(file, processedDirectory.OutputDir)
			fmt.Printf("ResizeImages, filesWorkers=%v, existing=%v\n", file.Name, pi.Existing)
			if err != nil {
				fmt.Println("ERROR resizing file: ", file)
				//panic(err)
				continue
			}

			//processedDirectory.ProcessedImages = append(processedDirectory.ProcessedImages, pi)
		}
	}
}
