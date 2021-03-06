package util

import (
	"context"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"
)

type File struct {
	Dir  string
	Name string
	W    int
	H    int
}

type ProcessedDirectory struct {
	InputDir        string
	BaseOutputDir   string
	OutputDir       string
	ManifestDir     string
	ProcessedImages []ProcessedImage
}

type ProcessedImage struct {
	Original  File
	Processed File
	Existing  bool
}

var TMin int64 = math.MaxInt64
var TMax int64

func Create(inputDir, baseOutputDir, imagesOutputDir, manifestOutputDir string) ProcessedDirectory {
	os.MkdirAll(manifestOutputDir, os.ModePerm)

	return ProcessedDirectory{
		InputDir:      inputDir,
		BaseOutputDir: baseOutputDir,
		OutputDir:     imagesOutputDir,
		ManifestDir:   manifestOutputDir,
	}
}

func SetupExitChannel() (context.Context, context.CancelFunc) {
	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// stop after pressing ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		fmt.Println("Press ctrl+c to interrupt...")
		<-c
		fmt.Println("Shutting down...")
		cancel()
	}()

	return ctx, cancel
}

func Timings(f string, start int64) {
	end := time.Now().UnixNano()
	if start < TMin {
		TMin = start
	}
	if end > TMax {
		TMax = end
	}

	fmt.Printf("%v, %v, %v\n", f, start, end)
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
	var result []File

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		w, h := GetImageBounds("images/" + f.Name())

		file := File{Name: f.Name(), W: w, H: h}
		result = append(result, file)
	}

	return result
}

func DisplayImageBounds(files []File) {
	for _, f := range files {
		fmt.Printf("f=%s, w=%d, h=%d\n", f.Name, f.W, f.H)
	}
}

func DecodeImage(filename string) (image.Image, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		fmt.Println("ERROR: ", filename, err)
		return nil, err
	}

	return src, nil
}

func GetImageBounds(filename string) (int, int) {
	src, _ := DecodeImage(filename)

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	return w, h
}

func ResizeFiles(in <-chan File, output string) <-chan ProcessedImage {
	out := make(chan ProcessedImage)

	go func() {
		fmt.Println("Start resize files")
		for file := range in {
			newName := output + "/" + file.Dir + "/" + file.Name + ".png"
			if _, err := os.Stat(newName); err == nil {
				continue
			}

			pi, err := ResizeFile(file, output)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Resized file %v\n", pi)
			out <- pi
		}

		close(out)
		fmt.Println("End resize files")
	}()

	return out
}

func ResizeFile(file File, output string) (ProcessedImage, error) {
	img, err := DecodeImage(file.Dir + "/" + file.Name)
	if err != nil {
		fmt.Println(file)
		return ProcessedImage{}, err
	}

	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	var w2, h2 int
	var existing bool

	newName := output + "/" + file.Name + ".png"
	if _, err := os.Stat(newName); err != nil {
		image2 := resize.Thumbnail(400, 300, img, resize.NearestNeighbor)
		bounds2 := image2.Bounds()
		w2, h2 = bounds2.Max.X, bounds2.Max.Y

		if w2 < 400 {
			image2 = resize.Resize(400, 0, image2, resize.NearestNeighbor)

			r := image.Rect(0, 0, 400, 300)

			if img2, ok := image2.(*image.NRGBA); ok {
				image2 = img2.SubImage(r)
			}
			if img2, ok := image2.(*image.RGBA); ok {
				image2 = img2.SubImage(r)
			}
			if img2, ok := image2.(*image.YCbCr); ok {
				image2 = img2.SubImage(r)
			}

			bounds2 = image2.Bounds()
			w2, h2 = bounds2.Max.X, bounds2.Max.Y
		} else if h2 < 300 {
			image2 = resize.Resize(0, 300, image2, resize.NearestNeighbor)

			r := image.Rect(0, 0, 400, 300)

			if img2, ok := image2.(*image.NRGBA); ok {
				image2 = img2.SubImage(r)
			}
			if img2, ok := image2.(*image.RGBA); ok {
				image2 = img2.SubImage(r)
			}
			if img2, ok := image2.(*image.YCbCr); ok {
				image2 = img2.SubImage(r)
			}

			bounds2 = image2.Bounds()
			w2, h2 = bounds2.Max.X, bounds2.Max.Y
		}

		outfile, err := os.Create(newName)
		if err != nil {
			return ProcessedImage{}, err
		}
		defer outfile.Close()

		if err := png.Encode(outfile, image2); err != nil {
			return ProcessedImage{}, err
		}
		outfile.Sync()

	} else {
		image2, err := DecodeImage(newName)
		if err != nil {
			fmt.Println(newName)
			return ProcessedImage{}, err
		}

		bounds2 := image2.Bounds()
		w2, h2 = bounds2.Max.X, bounds2.Max.Y
		existing = true
	}

	file.W = w
	file.H = h

	processed := File{
		Dir:  output,
		Name: file.Name + ".png",
		W:    w2,
		H:    h2,
	}

	pi := ProcessedImage{
		Original:  file,
		Processed: processed,
		Existing:  existing,
	}

	return pi, nil
}

func CreatePI(file File, output string) (ProcessedImage, error) {
	img, err := DecodeImage(file.Dir + "/" + file.Name)
	if err != nil {
		fmt.Println(file)
		return ProcessedImage{}, err
	}

	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	file.W = w
	file.H = h

	var w2, h2 int

	newName := output + "/" + file.Name + ".png"
	if _, err := os.Stat(newName); err == nil {
		image2, err := DecodeImage(newName)
		if err != nil {
			fmt.Println(newName)
			return ProcessedImage{}, err
		}

		bounds2 := image2.Bounds()
		w2, h2 = bounds2.Max.X, bounds2.Max.Y
	}

	processed := File{
		Dir:  output,
		Name: file.Name + ".png",
		W:    w2,
		H:    h2,
	}

	pi := ProcessedImage{
		Original:  file,
		Processed: processed,
	}

	return pi, nil
}

func Max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}
