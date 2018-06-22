package util

import (
    "fmt"
	"os"
	"io/ioutil"
	"log"
	"path/filepath"
    "image/png"
    "github.com/nfnt/resize"
)

func WalkDirectories() <- chan string {
    out := make(chan string)

    go func() {
        filepath.Walk("images/", func (path string, info os.FileInfo, err error) error {
            if err != nil {
                fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
                return nil
            }

            if info.IsDir() {
                out <- path
            }

            return nil
        })

        close(out)
    }()

	return out
}

func WalkFiles(in <- chan string) <- chan File {
    out := make(chan File)

    go func() {
        for dir := range in {
            files, err := ioutil.ReadDir(dir)
            if err != nil {
                log.Fatal(err)
            }

            for _, f := range files {
                if f.IsDir() {
                    continue
                }

                file := File{
                    Dir: dir,
                    Name: f.Name(),
                }

                out <- file
            }
        }

        close(out)
    }()

    return out
}

func CreateDirectories(in <- chan string) <- chan string {
    out := make(chan string)

    go func() {
        for dir := range in {
            out <- dir

            if _, err := os.Stat("output/" + dir); os.IsNotExist(err) {
                os.MkdirAll("output/" + dir, os.ModePerm)
            }
        }

        close(out)
    }()

    return out
}

func ResizeFiles(in <- chan File) <- chan File {
    out := make(chan File)

    go func() {
        for file := range in {
            newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"

            if _, err := os.Stat(newName); err == nil {
                continue
            }

            image := DecodeImage(file.Dir + "/" + file.Name)
            image2 := resize.Thumbnail(400, 300, image, resize.NearestNeighbor)

            outfile, err := os.Create(newName)
            if err != nil {
                panic(err)
            }
            defer outfile.Close()
            png.Encode(outfile, image2)

            out <- file
        }

        close(out)
    }()

    return out
}