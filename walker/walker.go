package walker

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetDirMax(dir string) int {
	var max int

	files3, _ := ioutil.ReadDir(dir)
	for _, f3 := range files3 {
		name := f3.Name()
		if strings.Contains(name, ".") {
			name = name[:strings.Index(name, ".")]
		}

		i, _ := strconv.Atoi(name)
		if i > max {
			max = i
		}
	}

	return max
}

func WalkDirectories(quit <-chan bool, dir string) <-chan string {
	out := make(chan string)

	go func() {
		stop := false

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			if stop {
				fmt.Println("Stop = true")
				return filepath.SkipDir
			}

			if info.IsDir() {
				out <- path
			}

			select {
			case <-quit:
				stop = true
				fmt.Println("Quitting WalkDirectories...")
				return filepath.SkipDir
			default:
				// do nothing
			}

			return nil
		})

		close(out)
	}()

	return out
}

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

			if f.Name() == ".DS_Store" {
				continue
			}
			if f.Name() == "_thumbdata" {
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

func WalkSlicedFiles(dir string) <-chan util.File {
	out := make(chan util.File)

	go func() {
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

func CreateDirectories(output string, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		for dir := range in {
			if _, err := os.Stat(output + "/images/" + dir); os.IsNotExist(err) {
				os.MkdirAll(output+"/images/"+dir, os.ModePerm)
			}
			out <- dir
		}

		close(out)
	}()

	return out
}
