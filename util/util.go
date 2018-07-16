package util

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type File struct {
	Dir  string
	Name string
	W    int
	H    int
}

type ProcessedImage struct {
	Original  File
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

			pi := ResizeFile(file, output)

			fmt.Printf("Resized file %v\n", pi)
			out <- pi
		}

		close(out)
		fmt.Println("End resize files")
	}()

	return out
}

func ResizeFile(file File, output string) ProcessedImage {
	img, err := DecodeImage(file.Dir + "/" + file.Name)
	if err != nil {
		fmt.Println(file)
		panic(err)
	}
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	var w2, h2 int
	newName := output + "/" + file.Dir + "/" + file.Name + ".png"
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
			panic(err)
		}
		defer outfile.Close()
		png.Encode(outfile, image2)
	} else {
		image2, err := DecodeImage(newName)
		if err != nil {
			fmt.Println(newName)
			panic(err)
		}

		bounds2 := image2.Bounds()
		w2, h2 = bounds2.Max.X, bounds2.Max.Y
	}

	file.W = w
	file.H = h

	processed := File{
		Dir:  "./output/" + file.Dir + "/",
		Name: file.Name + ".png",
		W:    w2,
		H:    h2,
	}

	pi := ProcessedImage{
		Original:  file,
		Processed: processed,
	}

	return pi
}

func Max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}
