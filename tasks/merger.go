package tasks

import (
	"context"
	"fmt"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/merger"
	"github.com/jvdanker/go-test/util"
	"os"
	"strings"
)

func MergeImages(ctx context.Context, dirs <-chan string, output string) {
	manifestOutput := strings.TrimSuffix(output, "/") + "/manifest"

	for inputdir := range dirs {
		fmt.Printf("MergeImages. input=%v\n", inputdir)

		mf := fmt.Sprintf("%v/%v/manifest.json", manifestOutput, inputdir)
		if _, err := os.Stat(mf); err != nil {
			continue
		}

		m, err := manifest.Read(mf)
		if err != nil {
			fmt.Printf("ERROR MergeImages, file=%v\n", mf)
			continue
		}

		dir := m.ImagesDir

		result := fmt.Sprintf("%v/result.png", dir)
		if _, err := os.Stat(result); err == nil {
			fmt.Printf("Skip dir: dir=%v\n", dir)
			continue
		}

		fmt.Printf("inputDir=%v\n", dir)

		// merge images into one image if there are images to process
		if len(m.Files) > 0 {
			bounds := m.Bounds()
			lm := layout.CreateBoxLayout()
			lm.Layout(bounds)
			m.Layout = lm

			// TODO write to separate output dir
			image := merger.MergeImages(m)
			util.CreateImage(m.ImagesDir+"/result.png", image)

			m.Update()
		}
	}
}
