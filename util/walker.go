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
    // fmt.Println("Walk directories")

    out := make(chan string)

    go func() {
        filepath.Walk("images/", func (path string, info os.FileInfo, err error) error {
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

func WalkFiles(in <- chan string) <- chan File {
    fmt.Println("Walk files")

    out := make(chan File)

    go func() {
        for dir := range in {
            fmt.Printf("Walkfiles, dir = %v\n", dir)

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

                fmt.Printf("Walk file: %v\n", file)
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

func ResizeFiles(in <- chan File) <- chan ProcessedImage {
    out := make(chan ProcessedImage)

    go func() {
        fmt.Println("Start resize files")
        for file := range in {
            newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"
            if _, err := os.Stat(newName); err == nil {
                continue
            }

            pi := ResizeFile(file)

            fmt.Printf("Resized file %v\n", pi)
            out <- pi
        }

        close(out)
        fmt.Println("End resize files")
    }()

    return out
}

func ResizeFile(file File) ProcessedImage {
    newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"

    image := DecodeImage(file.Dir + "/" + file.Name)
    bounds := image.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y

    image2 := resize.Thumbnail(400, 300, image, resize.NearestNeighbor)
    bounds2 := image2.Bounds()
    w2, h2 := bounds2.Max.X, bounds2.Max.Y

    outfile, err := os.Create(newName)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image2)

    file.W = w
    file.H = h

    processed := File{
        Dir: "./output/" + file.Dir + "/",
        Name: file.Name + "_400x300.png",
        W: w2,
        H: h2,
    }

    pi := ProcessedImage{
        Original: file,
        Processed: processed,
    }

    return pi
}