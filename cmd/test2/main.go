package main

import (
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	_ "github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func getImage(filename string) image.Image {
	img, err := util.DecodeImage(filename)
	if err != nil {
		canvas := image.NewRGBA(image.Rectangle{
			image.Point{0, 0},
			image.Point{256, 256}})
		return canvas
	}

	return img
}

func combineImages(x, y, z int) {
	canvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{512, 512}})

	img := getImage(fmt.Sprintf("output/parts/%v/%v/%v.png", z, x, y))
	draw.Draw(
		canvas,
		image.Rectangle{image.ZP, image.Point{256, 256}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("output/parts/%v/%v/%v.png", z, x+1, y))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{256, 0}, image.Point{512, 256}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("output/parts/%v/%v/%v.png", z, x, y+1))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{0, 256}, image.Point{256, 512}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("output/parts/%v/%v/%v.png", z, x+1, y+1))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{256, 256}, image.Point{512, 512}},
		img,
		image.ZP,
		draw.Src)

	outputDir := fmt.Sprintf("output/parts/%v/%v", z-1, x/2)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	image2 := resize.Resize(256, 256, canvas, resize.NearestNeighbor)
	util.CreateImage(fmt.Sprintf("%v/%v.png", outputDir, y/2), image2)
}

func main() {
	fmt.Println("test")

	//tasks.ResizeImages()
	//tasks.SliceImages()
	tasks.CreateBottomLayer()

	if _, err := os.Stat("output/parts"); os.IsNotExist(err) {
		os.MkdirAll("output/parts", os.ModePerm)
	}

	for {
		files, err := ioutil.ReadDir("output/parts")
		if err != nil {
			log.Fatal(err)
		}

		zDir := files[0].Name()
		z, _ := strconv.Atoi(zDir)

		if z == 0 {
			break
		}

		fmt.Println("output/parts/", zDir)

		files2, err := ioutil.ReadDir("output/parts/" + zDir)
		for _, f := range files2 {
			x, _ := strconv.Atoi(f.Name())

			if x%2 != 0 {
				continue
			}

			fmt.Println("x=", x)

			max := walker.GetDirMax("output/parts/" + zDir + "/" + f.Name())
			fmt.Println("max=", max)

			for y2 := 0; y2 <= max; y2 += 2 {
				fmt.Println(f.Name(), x, y2)

				fmt.Printf("\t%v, %v\n", x, y2)
				fmt.Printf("\t%v, %v\n", x+1, y2)
				fmt.Printf("\t%v, %v\n", x, y2+1)
				fmt.Printf("\t%v, %v -> as %v, %v, %v\n", x+1, y2+1, z-1, x/2, y2/2)

				combineImages(x, y2, z)
			}
		}
	}
}
