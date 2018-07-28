package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/merger"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
	"strings"
	"sync"
)

func MergeImages(dir string) {
	dir = strings.TrimSuffix(dir, "/") + "/images"

	fmt.Printf("MergeImages. input=%v\n", dir)
	dirs := walker.WalkDirectories(dir)

	wg := sync.WaitGroup{}

	for w := 0; w < 1; w++ {
		wg.Add(1)
		go func(w int) {
			mergeImageWorker(w, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

func mergeImageWorker(worker int, dirs <-chan string) {
	for dir := range dirs {
		result := fmt.Sprintf("%v/manifest.json", dir)
		if _, err := os.Stat(result); err != nil {
			fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		result = fmt.Sprintf("%v/result.png", dir)
		if _, err := os.Stat(result); err == nil {
			fmt.Printf("Skip dir: worker=%v, dir=%v\n", worker, dir)
			continue
		}

		fmt.Printf("mergeImageWorker=%v: inputDir=%v\n", worker, dir)

		mf := fmt.Sprintf("%v/manifest.json", dir)
		m, err := manifest.Read(mf)
		if err != nil {
			fmt.Println(dir)
			panic(err)
		}

		// merge images into one image if there are images to process
		if len(m.Files) > 0 {
			mergeImages(worker, m)
		}
	}
}

func mergeImages(dirWorker int, m manifest.ManifestFile) {
	//fmt.Printf("dirWorker=%v: mergeImages\n", dirWorker)

	bounds := m.Bounds()
	//fmt.Println(bounds)

	lm := layout.CreateBoxLayout()
	lm.Layout(bounds)
	m.Layout = lm

	image := merger.MergeImages(m)
	util.CreateImage(m.OutputDir+"/result.png", image)

	m.Update()
}
