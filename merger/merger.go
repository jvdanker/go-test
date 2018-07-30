package merger

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/util"
	"image"
	"image/draw"
)

func MergeImages(manifest manifest.ManifestFile) image.Image {
	canvas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{int(manifest.Layout.TotalWidth), int(manifest.Layout.TotalHeight)}})

	for i, pos := range manifest.Layout.Positions {
		//fmt.Println(i, pos)

		src := manifest.Files[i].Processed
		img, err := util.DecodeImage(manifest.ImagesDir + "/" + src.Name)
		if err != nil {
			fmt.Println(manifest.Files[i])
			panic(err)
		}

		draw.Draw(
			canvas,
			image.Rectangle{pos, pos.Add(image.Point{src.W, src.H})},
			img,
			image.ZP,
			draw.Src)
	}

	return canvas
}
