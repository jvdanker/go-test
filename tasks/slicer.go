package tasks

import (
	"context"
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"image"
	"image/draw"
	"os"
)

//in <-chan manifest.ManifestFile
func SliceImages(ctx context.Context, dirs <-chan string, output string) {
	for dir := range dirs {
		select {
		case <-ctx.Done():
			fmt.Println("Aborting SliceImages...")
			return
		default:
			// do nothing
		}

		mf := fmt.Sprintf("%v/manifest.json", dir)
		m, err := manifest.Read(mf)
		if err != nil {
			fmt.Println("ERROR reading manifest, ", mf)
			continue
		}

		fmt.Println("Slice Images, input=", m.InputDir, ", output=", m.ImagesDir)

		img, err := util.DecodeImage(m.ImagesDir + "/result.png")
		if err != nil {
			fmt.Println("ERROR reading image, ", m.ImagesDir+"/result.png")
			continue
		}

		m.SlicedDir = output + "/" + m.InputDir
		if _, err := os.Stat(m.SlicedDir); os.IsNotExist(err) {
			os.MkdirAll(m.SlicedDir, os.ModePerm)
		}
		m.Update()

		//var x,y int
		bounds := img.Bounds()
		w, h := bounds.Max.X, bounds.Max.Y
		fmt.Println("dir=", m.ImagesDir, ", w, h=", w, h)

		var x, y, i, j int
		i = 0
		j = 0

		for y < h {
			for x = 0; x < w; x += 256 {
				r := image.Rect(x, y, x+256, y+256)
				//fmt.Println("Slicer = ", r)

				var sub image.Image
				if img2, ok := img.(*image.NRGBA); ok {
					sub = img2.SubImage(r)
				}
				if img2, ok := img.(*image.RGBA); ok {
					sub = img2.SubImage(r)
				}

				b2 := sub.Bounds()

				canvas := image.NewRGBA(image.Rectangle{
					image.Point{0, 0},
					image.Point{256, 256}})

				draw.Draw(
					canvas,
					image.Rectangle{image.Point{0, 0}, image.Point{b2.Max.X, b2.Max.Y}},
					sub,
					image.Point{x, y},
					draw.Src)

				filename := fmt.Sprintf("%s/%d-%d.png", m.SlicedDir, i, j)

				util.CreateImage(filename, canvas)

				i++
			}

			y += 256
			j++
			i = 0
		}
	}
}
