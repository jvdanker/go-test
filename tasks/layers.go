package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func CreateZoomLayers(rootDir string) {
	fmt.Println("CreateZoomLayers, rootDir=", rootDir)

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		os.MkdirAll(rootDir, os.ModePerm)
	}

	for {
		files, err := ioutil.ReadDir(rootDir)
		if err != nil {
			log.Fatal(err)
		}

		zDir := files[0].Name()
		z, _ := strconv.Atoi(zDir)

		// stop when the top layer is reached
		if z == 0 {
			break
		}

		fmt.Println(rootDir, zDir)

		layerDir := fmt.Sprintf("%v/%v", rootDir, zDir)
		files2, err := ioutil.ReadDir(layerDir)
		for _, f := range files2 {
			x, _ := strconv.Atoi(f.Name())

			if x%2 != 0 {
				continue
			}

			fmt.Println("x=", x)

			max := walker.GetDirMax(layerDir + "/" + f.Name())
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

	fmt.Println()
}

func CreateBottomLayer() {
	fmt.Println("CreateBottomLayer")

	items, w, h := getMaxBounds()
	maxzoom := int(math.Ceil(math.Log(float64(util.Max(w, h)/256)) / math.Log(2)))
	fmt.Printf("max zoom = %v\n", maxzoom)

	itemsPerRow := int(math.Ceil(math.Sqrt(float64(items))))
	fmt.Println("itemsPerRow", itemsPerRow)
	fmt.Println()

	var i, z, x, y, maxY int
	dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}

		fmt.Println("dir", dir, z, x, y)

		var imgX int

		files := walker.WalkSlicedFiles(dir)
		for file := range files {
			s := strings.TrimPrefix(file.Name, "sub-")
			s = strings.TrimSuffix(s, ".png")
			parts := strings.Split(s, "-")

			fx, _ := strconv.Atoi(parts[0])
			fy, _ := strconv.Atoi(parts[1])

			nx := x + fx
			ny := y + fy

			//fmt.Printf("file=%v loc=%d, %d newloc=%d, %d\n", file, fx, fy, nx, ny)
			//fmt.Printf("output/parts/%d/%d/%d.png\n", maxzoom, nx, ny)

			oldName := fmt.Sprintf("../../../../%s%s", dir, file.Name)
			newName := fmt.Sprintf("output/parts/%d/%d/%d.png", maxzoom, nx, ny)

			//fmt.Println("symlink", oldName, newName)

			if _, err := os.Stat(newName); err == nil {
				//fmt.Println(" exists ", newName)
				os.Remove(newName)
			}

			if _, err := os.Stat(fmt.Sprintf("output/parts/%d/%d", maxzoom, nx)); os.IsNotExist(err) {
				os.MkdirAll(fmt.Sprintf("output/parts/%d/%d", maxzoom, nx), os.ModePerm)
			}

			err := os.Symlink(oldName, newName)
			if err != nil {
				panic(err)
			}

			fx++
			fy++

			if fx > imgX {
				imgX = fx
			}
			if fy > maxY {
				maxY = fy
			}
		}

		x += imgX

		if (i+1)%itemsPerRow == 0 {
			//fmt.Println("reset", maxY)
			x = 0
			y += maxY
			maxY = 0
		}

		i++

		//fmt.Println()
	}

	fmt.Println()
}

func getMaxBounds() (int, uint32, uint32) {
	var total int
	var tx, ty uint32

	dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
		fmt.Println(dir)

		m, err := manifest.Read(dir + "/manifest.json")
		if err == nil {
			total++

			tx += m.Layout.TotalWidth
			ty += m.Layout.TotalHeight
		}
	}
	fmt.Printf("Total %d, w=%d, h=%d\n", total, tx, ty)

	return total, tx, ty
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
