package main

import (
	"encoding/json"
	"fmt"
	"github.com/jvdanker/go-test/util"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"math"
	"image/png"
)

func GetFilesInDir(dir string) <-chan util.File {
	out := make(chan util.File)

	go func() {
		// fmt.Printf("Walkfiles, dir = %v\n", dir)

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			file := util.File{
				Dir:  dir,
				Name: f.Name(),
			}

			// fmt.Printf("Walk file: %v\n", file)
			out <- file
		}

		close(out)
	}()

	return out
}

func createManifest(files []util.ProcessedImage, id int, dir string) util.Manifest {
	fmt.Println("create manifest", dir)

	manifest := util.Manifest{
	    InputDir: dir,
	    OutputDir: "./output/" + dir,
	    Files: files,
	}

	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	outfile, err := os.Create("./output/" + dir + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(b)

	return manifest
}

func mergeFiles(manifest util.Manifest) {
    var files []util.File
    for _, file := range manifest.Files {
        files = append(files, file.Processed)
    }

    var itemsPerRow = int(math.Ceil(math.Sqrt(float64(len(files)))))
    fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(files), itemsPerRow)

    maxWidth, maxHeight := util.CalculateMaxWidthAndHeight(files, itemsPerRow)
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", maxWidth, maxHeight)

    image := util.MergeImages(files, maxWidth, maxHeight, itemsPerRow)

    outfilename := manifest.OutputDir + "/result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
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
		files := GetFilesInDir(dir)
        images := filesWorker(id, files)
        manifest := createManifest(images, id, dir)
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
