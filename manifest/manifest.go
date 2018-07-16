package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/jvdanker/go-test/layout"
	"github.com/jvdanker/go-test/util"
	"image"
	"io/ioutil"
	"os"
)

type Manifest interface {
	Update()
	Bounds(itemsPerRow int) []image.Point
}

type ManifestFile struct {
	InputDir  string
	OutputDir string
	Files     []File
	Layout    layout.LayoutManager
}

type File struct {
	Original  Image
	Processed Image
}

type Image struct {
	Name string
	W    int
	H    int
}

func Create(processedImages []util.ProcessedImage, inputdir, outputdir string) ManifestFile {
	files := []File{}
	for _, img := range processedImages {
		o := Image{
			Name: img.Original.Name,
			W:    img.Original.W,
			H:    img.Original.H,
		}
		p := Image{
			Name: img.Processed.Name,
			W:    img.Processed.W,
			H:    img.Processed.H,
		}
		file := File{
			Original:  o,
			Processed: p,
		}
		files = append(files, file)
	}

	manifest := ManifestFile{
		InputDir:  inputdir,
		OutputDir: outputdir + "/" + inputdir,
		Files:     files,
	}

	b, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	outfile, err := os.Create(outputdir + "/" + inputdir + "/manifest.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(b)

	return manifest
}

func Read(name string) (ManifestFile, error) {
	var m ManifestFile

	f, err := os.Open(name)
	if err != nil {
		return m, err
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	json.Unmarshal(byteValue, &m)

	return m, nil
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

func (m ManifestFile) Bounds() []image.Point {
	var result = []image.Point{}

	for _, f := range m.Files {
		result = append(result, image.Point{X: f.Processed.W, Y: f.Processed.H})
	}

	return result
}
