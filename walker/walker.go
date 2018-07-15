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

func WalkDirectories(dir string) <-chan string {
	// fmt.Println("Walk directories")

	out := make(chan string)

	go func() {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return nil
			}

			if info.IsDir() {
				//fmt.Printf("Walk dirs: %v\n", path)
				out <- path
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

			if !strings.HasPrefix(f.Name(), "sub-") {
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

func CreateDirectories(in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		for dir := range in {
			if _, err := os.Stat("output/" + dir); os.IsNotExist(err) {
				os.MkdirAll("output/"+dir, os.ModePerm)
			}
			out <- dir
		}

		close(out)
	}()

	return out
}
