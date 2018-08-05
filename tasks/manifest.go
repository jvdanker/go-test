package tasks

import (
	"context"
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"strings"
)

func CreateManifests(ctx context.Context, dirs <-chan string, output string) {
	imagesOutput := strings.TrimSuffix(output, "/") + "/images"
	manifestOutput := strings.TrimSuffix(output, "/") + "/manifest"

	for inputdir := range dirs {
		fmt.Printf("CreateManifests, dirWorker=%v\n", inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", manifestOutput, inputdir)); err == nil {
			continue
		}

		files := walker.WalkFiles(inputdir)

		var processedDirectory = util.Create(
			inputdir,
			output,
			imagesOutput+"/"+inputdir,
			manifestOutput+"/"+inputdir)

		for file := range files {
			select {
			case <-ctx.Done():
				fmt.Println("Aborting CreateManifests...")
				return
			default:
				// do nothing
			}

			// skip if dest image doesn't exist
			dest := fmt.Sprintf("%v/%v.png", processedDirectory.OutputDir, file.Name)
			if _, err := os.Stat(dest); err != nil {
				continue
			}

			fmt.Printf("CreateManifests, filesWorkers=%v\n", file.Name)
			pi, err := util.CreatePI(file, processedDirectory.OutputDir)
			if err != nil {
				fmt.Println("ERROR CreateManifests, file: ", file)
				continue
			}

			processedDirectory.ProcessedImages = append(processedDirectory.ProcessedImages, pi)
		}

		if len(processedDirectory.ProcessedImages) > 0 {
			manifest.Create(processedDirectory)
		}
	}
}

func CreateManifest(in <-chan util.ProcessedDirectory) <-chan manifest.ManifestFile {
	out := make(chan manifest.ManifestFile)

	go func() {
		for pd := range in {
			//fmt.Println(pd)

			if len(pd.ProcessedImages) > 0 {
				// create manifest file
				fmt.Printf("CreateManifest, dirWorker=%v: count=%v\n", pd.InputDir, len(pd.ProcessedImages))
				out <- manifest.Create(pd)
			}
		}

		close(out)
	}()

	return out
}
