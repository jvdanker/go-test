package main

import (
	"encoding/json"
	"fmt"
	"github.com/jvdanker/go-test/util"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func WalkFiles(dir string) <-chan util.File {
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

func createManifest(out <-chan util.ProcessedImage, id int, dir string) {
	fmt.Println("should create manifest", dir)
	var result []util.ProcessedImage
	for file := range out {
		result = append(result, file)
	}

	fmt.Println("create manifest", dir)

	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	outfile, err := os.Create("./output/" + dir + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(b)
}

func filesWorker(id int, files <-chan util.File, out chan<- util.ProcessedImage) {
	for file := range files {
		// fmt.Println(id, files, file)
		newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"
		if _, err := os.Stat(newName); err == nil {
			// out <- file
			continue
		}

		file2 := util.ResizeFile(file)
		out <- file2
	}
}

func dirWorker(id int, dirs <-chan string) {
	for dir := range dirs {
		files := WalkFiles(dir)

		out := make(chan util.ProcessedImage)
		wg2 := sync.WaitGroup{}

		// create workers
		for w := 1; w <= 3; w++ {
			wg2.Add(1)
			go func(w int) {
				filesWorker(w, files, out)
				wg2.Done()
			}(w)
		}

		go createManifest(out, id, dir)

		fmt.Println("wait")
		wg2.Wait()
		close(out)
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
