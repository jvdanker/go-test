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
     "image/draw"
)

type File struct {
    Dir string
    Name string
    W int
    H int
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


func MergeImages(files []File, maxWidth, maxHeight, itemsPerRow int) image.Image {
    var x, y int
    var curr int

    img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxWidth, maxHeight}})
    maxHeight = 0

    for _, f := range files {
        if strings.HasPrefix(f.Name, ".") {
            continue
        }

        src := DecodeImage("images/" + f.Name)
        b := src.Bounds()

        if b.Max.Y > maxHeight {
            maxHeight = b.Max.Y
        }

        var a = image.Point{x, y}

        draw.Draw(
            img,
            image.Rectangle{a, a.Add(image.Point{b.Max.X, b.Max.Y})},
            src,
            image.ZP,
            draw.Src)

        x += b.Max.X
        curr++

        if curr % itemsPerRow == 0 {
            x = 0
            y += maxHeight
            maxHeight = 0
        }
    }

    return img
}

func CalculateMaxWidthAndHeight(files []File, itemsPerRow int) (int, int) {
    var maxWidth, maxHeight int
    var curr, currWidth, currHeight int

    for _, f := range files {
        currWidth += f.W

        if currWidth > maxWidth {
            maxWidth = currWidth
        }

        if curr % itemsPerRow == 0 {
            currWidth = 0
            currHeight += f.H
        }

        if currHeight > maxHeight {
            maxHeight = currHeight
        }

        curr++
    }

    return maxWidth, maxHeight
}