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

		layerDir := fmt.Sprintf("%v/%v", rootDir, zDir)
		fmt.Printf("Processing images in layer %v\n", layerDir)

		files2, err := ioutil.ReadDir(layerDir)
		for _, f := range files2 {
			x, _ := strconv.Atoi(f.Name())

			if x%2 != 0 {
				continue
			}

			//fmt.Println("x=", x)

			max := walker.GetDirMax(layerDir + "/" + f.Name())
			//fmt.Println("max=", max)

			for y2 := 0; y2 <= max; y2 += 2 {
				//fmt.Println(f.Name(), x, y2)

				//fmt.Printf("\t%v, %v\n", x, y2)
				//fmt.Printf("\t%v, %v\n", x+1, y2)
				//fmt.Printf("\t%v, %v\n", x, y2+1)
				//fmt.Printf("\t%v, %v -> as %v, %v, %v\n", x+1, y2+1, z-1, x/2, y2/2)

				combineImages(rootDir, x, y2, z)
			}
		}
	}

	fmt.Println()
}

func CreateBottomLayer(input, slices, output string) map[string]string {
	fmt.Println("CreateBottomLayer, input=", input, ", slices=", slices, "output=", output)

	items, w, h := getMaxBounds(input)
	maxzoom := int(math.Ceil(math.Log(float64(util.Max(w, h)/256)) / math.Log(2)))
	fmt.Printf("max zoom = %v\n", maxzoom)

	itemsPerRow := int(math.Ceil(math.Sqrt(float64(items))))
	fmt.Println("itemsPerRow", itemsPerRow)
	fmt.Println()

	result := map[string]string{}
	result["items"] = strconv.Itoa(items)
	result["maxzoom"] = strconv.Itoa(maxzoom)
	result["w"] = strconv.FormatUint(uint64(w), 10)
	result["h"] = strconv.FormatUint(uint64(h), 10)

	var i, z, x, y, maxY int

	dirs := walker.WalkDirectories(slices)
	for dir := range dirs {
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}

		fmt.Println("dir", dir, z, x, y)

		var imgX int

		files := walker.WalkSlicedFiles(dir)
		for file := range files {
			s := strings.TrimSuffix(file.Name, ".png")
			parts := strings.Split(s, "-")

			fx, _ := strconv.Atoi(parts[0])
			fy, _ := strconv.Atoi(parts[1])

			nx := x + fx
			ny := y + fy

			//fmt.Printf("file=%v loc=%d, %d newloc=%d, %d\n", file, fx, fy, nx, ny)
			//fmt.Printf("output/parts/%d/%d/%d.png\n", maxzoom, nx, ny)

			oldName := fmt.Sprintf("../../../../%s%s", dir, file.Name)
			targetDir := fmt.Sprintf("%v/%d/%d", output, maxzoom, nx)
			newName := fmt.Sprintf("%s/%d.png", targetDir, ny)

			//fmt.Println("symlink", oldName, newName)

			if _, err := os.Stat(newName); err == nil {
				//fmt.Println(" exists ", newName)
				os.Remove(newName)
			}

			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				os.MkdirAll(targetDir, os.ModePerm)
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

	return result
}

func getMaxBounds(input string) (int, uint32, uint32) {
	var total int
	var tx, ty uint32

	dirs := walker.WalkDirectories(input)
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

func combineImages(output string, x, y, z int) {
	canvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{512, 512}})

	img := getImage(fmt.Sprintf("%v/%v/%v/%v.png", output, z, x, y))
	draw.Draw(
		canvas,
		image.Rectangle{image.ZP, image.Point{X: 256, Y: 256}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("%v/%v/%v/%v.png", output, z, x+1, y))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{X: 256, Y: 0}, image.Point{X: 512, Y: 256}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("%v/%v/%v/%v.png", output, z, x, y+1))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{X: 0, Y: 256}, image.Point{X: 256, Y: 512}},
		img,
		image.ZP,
		draw.Src)

	img = getImage(fmt.Sprintf("%v/%v/%v/%v.png", output, z, x+1, y+1))
	draw.Draw(
		canvas,
		image.Rectangle{image.Point{256, 256}, image.Point{512, 512}},
		img,
		image.ZP,
		draw.Src)

	outputDir := fmt.Sprintf("%v/%v/%v", output, z-1, x/2)

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
