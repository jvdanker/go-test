package manifest

import (
    "fmt"
    "encoding/json"
    "os"
    "image"
    "github.com/jvdanker/go-test/util"
)

type manifest interface {
    Update()
    Bounds(itemsPerRow int) image.Rectangle
}

type ManifestFile struct {
    InputDir string
    OutputDir string
    Files []File
    TotalWidth int
    TotalHeight int
    ItemsPerRow int
}

type File struct {
    Original Image
    Processed Image
}

type Image struct {
    Name string
    W int
    H int
}

func Create(processedImages []util.ProcessedImage, dir string) ManifestFile {
	fmt.Println("create manifest", dir)

    files := []File{}
    for _, img := range processedImages {
        o := Image{
            Name: img.Original.Name,
            W: img.Original.W,
            H: img.Original.H,
        }
        p := Image{
            Name: img.Processed.Name,
            W: img.Processed.W,
            H: img.Processed.H,
        }
        file := File{
            Original: o,
            Processed: p,
        }
        files = append(files, file)
    }

	manifest := ManifestFile{
	    InputDir: dir,
	    OutputDir: "./output/" + dir,
	    Files: files,
	}

	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	outfile, err := os.Create("./output/" + dir + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(b)

	return manifest
}

func (m ManifestFile) Update() {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	outfile, err := os.Create(m.OutputDir + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(b)
}


func (m ManifestFile) Bounds() image.Rectangle {
    var maxWidth, maxHeight int
    var currWidth, currHeight int

    for i, f := range m.Files {
        currWidth += f.Processed.W

        if currWidth > maxWidth {
            maxWidth = currWidth
        }

        if (i + 1) % m.ItemsPerRow == 0 {
            currWidth = 0
            currHeight += f.Processed.H
        }

        if currHeight > maxHeight {
            maxHeight = currHeight
        }
    }

    return image.Rectangle{Max: image.Point{maxWidth, maxHeight}}
}
