package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"image"
	"image/draw"
)

func SliceImages() {
	fmt.Println("Slice Images")

	dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
		fmt.Println(dir)

		img, err := util.DecodeImage(dir + "/result.png")
		if err != nil {
			fmt.Println(err)
			continue
		}

		//var x,y int
		bounds := img.Bounds()
		w, h := bounds.Max.X, bounds.Max.Y
		fmt.Println("w,h=", w, h)

		var x, y, i, j int
		i = 0
		j = 0

		for y < h {
			for x = 0; x < w; x += 256 {
				r := image.Rect(x, y, x+256, y+256)
				fmt.Println("Slicer = ", r)

				var sub image.Image
				if img2, ok := img.(*image.NRGBA); ok {
					sub = img2.SubImage(r)
				}
				if img2, ok := img.(*image.RGBA); ok {
					sub = img2.SubImage(r)
				}

				b2 := sub.Bounds()
				fmt.Printf("Bounds = %d\n", b2)

				canvas := image.NewRGBA(image.Rectangle{
					image.Point{0, 0},
					image.Point{256, 256}})

				draw.Draw(
					canvas,
					image.Rectangle{image.Point{0, 0}, image.Point{b2.Max.X, b2.Max.Y}},
					sub,
					image.Point{x, y},
					draw.Src)

				util.CreateImage(fmt.Sprintf("%s/sub-%d-%d.png", dir, i, j), canvas)

				i++
			}

			y += 256
			j++
			i = 0
		}

		fmt.Println()
	}
}
