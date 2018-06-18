package main

import (
    "fmt"
    "io/ioutil"
     "log"
     "os"
     "image"
     "strings"
     _ "image/jpeg"
     "image/png"
     "math"
     "image/draw"
)

type file struct {
    name string
    w int
    h int
}

func main() {
    files := getImages("./images")

    displayImageBounds(files)

    var itemsPerRow = int(math.Ceil(math.Sqrt(float64(len(files)))))
    fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(files), itemsPerRow)

    maxWidth, maxHeight := calculateMaxWidthAndHeight(files, itemsPerRow)
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", maxWidth, maxHeight)

    image := mergeImages(files, maxWidth, maxHeight, itemsPerRow)

    outfilename := "result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
}

func getImages(dir string) []file {
    result := []file{}

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := getImageBounds("images/" + f.Name())

        file := file{ name: f.Name(), w: w, h: h }
        result = append(result, file)
    }

    return result
}

func displayImageBounds(files []file) {
    for _, f := range files {
        fmt.Printf("f=%s, w=%d, h=%d\n", f.name, f.w, f.h)
    }
}

func mergeImages(files []file, maxWidth, maxHeight, itemsPerRow int) image.Image {
    var x, y int
    var curr int

    img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxWidth, maxHeight}})
    maxHeight = 0

    for _, f := range files {
        if strings.HasPrefix(f.name, ".") {
            continue
        }

        src := decodeImage("images/" + f.name)
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

func calculateMaxWidthAndHeight(files []file, itemsPerRow int) (int, int) {
    var maxWidth, maxHeight int
    var curr, currWidth, currHeight int

    for _, f := range files {
        currWidth += f.w

        if currWidth > maxWidth {
            maxWidth = currWidth
        }

        if curr % itemsPerRow == 0 {
            currWidth = 0
            currHeight += f.h
        }

        if currHeight > maxHeight {
            maxHeight = currHeight
        }

        curr++
    }

    return maxWidth, maxHeight
}

func decodeImage(filename string) image.Image {
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

func getImageBounds(filename string) (int, int) {
    src := decodeImage(filename)

    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y

    return w, h
}
