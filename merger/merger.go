package merger

import(
    "fmt"
    "image"
    "image/draw"
    "github.com/jvdanker/go-test/util"
    "github.com/jvdanker/go-test/manifest"
)

func MergeImages(manifest manifest.ManifestFile) image.Image {
    canvas := image.NewRGBA(image.Rectangle{
        image.Point{0, 0},
        image.Point{manifest.Layout.TotalWidth, manifest.Layout.TotalHeight}})

    for i, pos := range manifest.Layout.Positions {
        fmt.Println(i, pos)

        src := manifest.Files[i].Processed
        img, _ := util.DecodeImage(manifest.OutputDir + "/" + src.Name)

        draw.Draw(
            canvas,
            image.Rectangle{pos, pos.Add(image.Point{src.W, src.H})},
            img,
            image.ZP,
            draw.Src)
    }

    return canvas
}
