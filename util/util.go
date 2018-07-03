package util

import (
    "fmt"
    "io/ioutil"
     "log"
     "os"
     "image"
     "strings"
     _ "image/jpeg"
     "image/png"
     "github.com/nfnt/resize"
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

func CreateImage(filename string, image image.Image) {
    outfile, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    png.Encode(outfile, image)
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
    image := DecodeImage(file.Dir + "/" + file.Name)
    bounds := image.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y

    var w2, h2 int
    newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"
    if _, err := os.Stat(newName); err != nil {
        image2 := resize.Thumbnail(400, 300, image, resize.NearestNeighbor)
        bounds2 := image2.Bounds()
        w2, h2 = bounds2.Max.X, bounds2.Max.Y

        outfile, err := os.Create(newName)
        if err != nil {
            panic(err)
        }
        defer outfile.Close()
        png.Encode(outfile, image2)
    } else {
        image2 := DecodeImage(newName)
        bounds2 := image2.Bounds()
        w2, h2 = bounds2.Max.X, bounds2.Max.Y
    }

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