package main

import (
    "fmt"
    "os"
    "math"
    "image/png"
    "github.com/jvdanker/go-test/util"
)

func main() {
    files := util.GetImages("./images")

    util.DisplayImageBounds(files)

    var itemsPerRow = int(math.Ceil(math.Sqrt(float64(len(files)))))
    fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(files), itemsPerRow)

    maxWidth, maxHeight := util.CalculateMaxWidthAndHeight(files, itemsPerRow)
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", maxWidth, maxHeight)

    image := util.MergeImages(files, maxWidth, maxHeight, itemsPerRow)

    outfilename := "result.png"
    outfile, err := os.Create(outfilename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
}


