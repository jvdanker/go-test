package walker

import (
    "fmt"
	"os"
	"io/ioutil"
	"log"
	"path/filepath"
    "github.com/jvdanker/go-test/util"
)

func WalkDirectories(dir string) <- chan string {
    // fmt.Println("Walk directories")

    out := make(chan string)

    go func() {
        filepath.Walk(dir, func (path string, info os.FileInfo, err error) error {
            if err != nil {
                fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
                return nil
            }

            if info.IsDir() {
                // fmt.Printf("Walk dirs: %v\n", path)
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

func CreateDirectories(in <- chan string) <- chan string {
    out := make(chan string)

    go func() {
        for dir := range in {
            if _, err := os.Stat("output/" + dir); os.IsNotExist(err) {
                os.MkdirAll("output/" + dir, os.ModePerm)
            }
            out <- dir
        }

        close(out)
    }()

    return out
}
