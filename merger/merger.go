package merger

import(
    "fmt"
    "image"
    "image/draw"
    "github.com/jvdanker/go-test/util"
    "github.com/jvdanker/go-test/layout"
)

func MergeImages(lm layout.LayoutManager) image.Image {
    result := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{lm.TotalWidth, lm.TotalHeight}})

    for i, p := range lm.Positions {
        fmt.Println(i, p)

        src := lm.Files[i].Processed
        img := util.DecodeImage(lm.OutputDir + "/" + src.Name)

        var pos = image.Point{p.X, p.Y}

        draw.Draw(
            result,
            image.Rectangle{pos, pos.Add(image.Point{src.W, src.H})},
            img,
            image.ZP,
            draw.Src)
    }

    return result
}
