package util

import (
    "fmt"
    "io/ioutil"
     "log"
     "os"
     "image"
     "strings"
     _ "image/jpeg"
     _ "image/png"
)

type File struct {
    Dir string
    Name string
    W int
    H int
}

type ProcessedImage struct {
    Original File
    Processed File
}

func GetImages(dir string) []File {
    result := []File{}

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := GetImageBounds("images/" + f.Name())

        file := File{ Name: f.Name(), W: w, H: h }
        result = append(result, file)
    }

    return result
}

func DisplayImageBounds(files []File) {
    for _, f := range files {
        fmt.Printf("f=%s, w=%d, h=%d\n", f.Name, f.W, f.H)
    }
}

func DecodeImage(filename string) image.Image {
    infile, err := os.Open(filename)
    if err != nil {
        // replace this with real error handling
        panic(err)
    }
    defer infile.Close()

    // Decode will figure out what type of image is in the file on its own.
    // We just have to be sure all the image packages we want are imported.
    src, _, err := image.Decode(infile)
    if err != nil {
        // replace this with real error handling
        panic(err)
    }

    return src
}

func GetImageBounds(filename string) (int, int) {
    src := DecodeImage(filename)

    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y

    return w, h
}
