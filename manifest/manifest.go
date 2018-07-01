package manifest

import (
    "fmt"
    "encoding/json"
    "os"
    "github.com/jvdanker/go-test/util"
)

type ManifestFile struct {
    InputDir string
    OutputDir string
    Files []File
    TotalWidth int
    TotalHeight int
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

func CreateManifest(processedImages []util.ProcessedImage, id int, dir string) ManifestFile {
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

func Update(m ManifestFile)  {
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