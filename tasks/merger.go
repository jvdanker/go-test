package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/merger"
	"github.com/jvdanker/go-test/util"
	"os"
)

func MergeImages(in <-chan manifest.ManifestFile) <-chan manifest.ManifestFile {
	out := make(chan manifest.ManifestFile)

	go func() {
		for m := range in {
			fmt.Printf("MergeImages. input=%v\n", m.ImagesDir)
			mergeImageWorker(&m)
			out <- m
		}

		close(out)
	}()

	return out
}

func mergeImageWorker(m *manifest.ManifestFile) {
	dir := m.ImagesDir

	result := fmt.Sprintf("%v/result.png", dir)
	if _, err := os.Stat(result); err == nil {
		fmt.Printf("Skip dir: dir=%v\n", dir)
		return
	}

	fmt.Printf("inputDir=%v\n", dir)

	// merge images into one image if there are images to process
	if len(m.Files) > 0 {
		mergeImages(m)
	}
}

func mergeImages(m *manifest.ManifestFile) {
	//fmt.Printf("dirWorker=%v: mergeImages\n", dirWorker)

	bounds := m.Bounds()
	//fmt.Println(bounds)

	lm := layout.CreateBoxLayout()
	lm.Layout(bounds)
	m.Layout = lm

	image := merger.MergeImages(*m)
	util.CreateImage(m.ImagesDir+"/result.png", image)
	m.Update()
}
