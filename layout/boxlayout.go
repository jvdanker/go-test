package layout

import (
    "fmt"
    "math"
    "image"
    "github.com/jvdanker/go-test/manifest"
)

func CreateBoxLayout(m manifest.ManifestFile) LayoutManager {
    return LayoutManager{
        ManifestFile: m,
    }
}

func (l *LayoutManager) Layout() {
    var x, y int
    var rowMaxHeight int

    l.ItemsPerRow = int(math.Ceil(math.Sqrt(float64(len(l.Files)))))

    for i, f := range l.Files {
        img := f.Processed
        fmt.Println("img", x, y, img)

        p := image.Point{X: x, Y: y}
        l.Positions = append(l.Positions, p)

        x += img.W

        if x > l.TotalWidth {
            l.TotalWidth = x
        }
        if y > rowMaxHeight {
            rowMaxHeight = y
        }

        if (i+1) % l.ItemsPerRow == 0 {
        fmt.Println("reset")
            x = 0
            y += rowMaxHeight
            rowMaxHeight = 0
        }

        fmt.Println("y=", l.TotalHeight)
    }

    fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(l.Files), l.ItemsPerRow)
    fmt.Printf("maxWidth=%d, maxHeight=%d\n", l.TotalWidth, l.TotalHeight)
}
