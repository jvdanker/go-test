package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
)

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
