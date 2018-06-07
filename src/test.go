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

func main() {
    files, err := ioutil.ReadDir("./images")
    if err != nil {
        log.Fatal(err)
    }

    var width, height, count int

    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := displayImageBounds("images/" + f.Name());
        w++
        h++
        count++
    }

    width++
    height++

    var items = int(math.Ceil(math.Sqrt(float64(count))))
    fmt.Println(items)

    var maxWidth, maxHeight, curr, currWidth, currHeight int
    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := displayImageBounds("images/" + f.Name());
        currWidth += w

        if currWidth > maxWidth {
            maxWidth = currWidth
        }

        if curr % items == 0 {
            currWidth = 0
            currHeight += h
        }

        if currHeight > maxHeight {
            maxHeight = currHeight
        }

        curr++
    }

    fmt.Printf("maxWidth=%d, maxHeight=%d\n", maxWidth, maxHeight)

    var x, y int
    curr = 0
    img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxWidth, maxHeight}})
    maxHeight = 0
    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        src := decodeImage("images/" + f.Name())
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

        if curr % items == 0 {
            x = 0
            y += maxHeight
            maxHeight = 0
        }
    }

    // Encode the grayscale image to the output file
    outfilename := "result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, img)
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

func displayImageBounds(filename string) (int, int) {
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

    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y

    fmt.Printf("f=%s, w=%d, h=%d\n", filename, w, h);

    return w, h
}