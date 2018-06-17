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
    count := displayImageBounds()

    var itemsPerRow = int(math.Ceil(math.Sqrt(float64(count))))
    fmt.Println(itemsPerRow)

    maxWidth, maxHeight := calculateMaxWidthAndHeight(itemsPerRow)
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", maxWidth, maxHeight)

    image := mergeImages(maxWidth, maxHeight, itemsPerRow)

    outfilename := "result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
}

func displayImageBounds() int {
    files, err := ioutil.ReadDir("./images")
    if err != nil {
        log.Fatal(err)
    }

    var count int

    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := getImageBounds("images/" + f.Name())
        fmt.Printf("f=%s, w=%d, h=%d\n", f.Name(), w, h)

        count++
    }

    return count
}

func mergeImages(maxWidth, maxHeight, itemsPerRow int) image.Image {
    var x, y int
    var curr int

    img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxWidth, maxHeight}})
    maxHeight = 0

    files, err := ioutil.ReadDir("./images")
    if err != nil {
        log.Fatal(err)
    }

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

        if curr % itemsPerRow == 0 {
            x = 0
            y += maxHeight
            maxHeight = 0
        }
    }

    return img
}

func calculateMaxWidthAndHeight(itemsPerRow int) (int, int) {
    var maxWidth, maxHeight int
    var curr, currWidth, currHeight int

    files, err := ioutil.ReadDir("./images")
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if strings.HasPrefix(f.Name(), ".") {
            continue
        }

        w, h := getImageBounds("images/" + f.Name())
        fmt.Printf("f=%s, w=%d, h=%d\n", f.Name(), w, h)

        currWidth += w

        if currWidth > maxWidth {
            maxWidth = currWidth
        }

        if curr % itemsPerRow == 0 {
            currWidth = 0
            currHeight += h
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
