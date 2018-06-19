package main

import (
    "os"
    "fmt"
    "image/png"
    "encoding/json"
    "math"
    "time"

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

    var count int
    var start = time.Now()

    for _, f := range files {
        //fmt.Print("\u001b[s");
        fmt.Print("\u001b7");
        fmt.Printf("%3.0f%%, %d items, %s",
            math.Round((float64(count) / float64(len(files))) * 100),
            len(files) - count,
            time.Since(start))
        //fmt.Print("\u001b[0K\u001b[u")
        fmt.Print("\u001b[0K\u001b8")

        image := util.DecodeImage("./images/" + f.Name)
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

        count++
        time.Sleep(1500 * time.Millisecond)
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

    fmt.Print("\u001b[0KDone\n")
}