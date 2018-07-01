package main

import (
	"fmt"
	"os"
	"sync"
	"math"
	"image/png"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/manifest"
)

func mergeFiles(m manifest.ManifestFile) {
    var files []util.File
    for _, file := range m.Files {
        f := util.File{
            Dir: m.OutputDir,
            Name: file.Processed.Name,
            W: file.Processed.W,
            H: file.Processed.H,
        }
        files = append(files, f)
    }

    m.ItemsPerRow = int(math.Ceil(math.Sqrt(float64(len(files)))))
    fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(files), m.ItemsPerRow)

    bounds := m.Bounds()
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", bounds.Max.X, bounds.Max.Y)

    image := util.MergeImages(files, bounds.Max.X, bounds.Max.Y, m.ItemsPerRow)

    outfilename := m.OutputDir + "/result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)

    // update manifest
    m.TotalWidth = bounds.Max.X
    m.TotalHeight = bounds.Max.Y
    m.Update()
}

func filesWorker(id int, files <-chan util.File) []util.ProcessedImage {
    var result = []util.ProcessedImage{}

	for file := range files {
		// fmt.Println(id, files, file)
        file2 := util.ResizeFile(file)
        result = append(result, file2)
	}

	return result
}

func dirWorker(id int, dirs <-chan string) {
	for dir := range dirs {
		files := util.WalkFiles(dir)
        images := filesWorker(id, files)
        manifest := manifest.Create(images, id, dir)
        mergeFiles(manifest)
	}
}

func main() {
	fmt.Println("test")

	dirs := util.WalkDirectories()
	dirs = util.CreateDirectories(dirs)

	wg := sync.WaitGroup{}

	for w := 1; w <= 1; w++ {
		wg.Add(1)
		go func(w int) {
			dirWorker(w, dirs)
			wg.Done()
		}(w)
	}

	wg.Wait()
}
