package main

import (
    "os"
    "fmt"
    "image/png"
    "encoding/json"

    "github.com/jvdanker/go-test/util"
    "github.com/nfnt/resize"
)

type ProcessedImage struct {
    Original util.File
    Processed util.File
}

func main() {
    files := util.GetImages("./images")

    fmt.Printf("Images to process : %d\n", len(files))

    result := []ProcessedImage{}

    for _, f := range files {
        image := util.DecodeImage("./images/" + f.Name);
        image2 := resize.Thumbnail(400, 300, image, resize.NearestNeighbor)

        newName := "./output/" + f.Name + "_400x300.png"
        outfile, err := os.Create(newName)
        if err != nil {
            panic(err)
        }
        defer outfile.Close()
        png.Encode(outfile, image2)

        w, h := util.GetImageBounds(newName)
        newFile := util.File{ Name: newName, W: w, H: h }
        pi := ProcessedImage{ f, newFile }

        result = append(result, pi)
    }

    b, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }

    outfile, err := os.Create("./output/manifest.json")
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    outfile.Write(b)
}