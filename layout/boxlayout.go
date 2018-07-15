package layout

import (
	"image"
	"math"
)

func CreateBoxLayout() LayoutManager {
	return LayoutManager{}
}

func (l *LayoutManager) Layout(bounds []image.Point) {
	var x, y int
	var rowMaxHeight int

	l.ItemsPerRow = int(math.Ceil(math.Sqrt(float64(len(bounds)))))

	for i, p := range bounds {
		l.Positions = append(l.Positions, image.Point{X: x, Y: y})

		x += p.X

		if uint32(x) > l.TotalWidth {
			l.TotalWidth = uint32(x)
		}

		if p.Y > rowMaxHeight {
			rowMaxHeight = p.Y
		}

		if (i+1)%l.ItemsPerRow == 0 {
			x = 0
			y += rowMaxHeight
			rowMaxHeight = 0
		}

		l.TotalHeight = uint32(y + rowMaxHeight)
	}

	//fmt.Printf("numberOfItems=%d, itemsPerRow=%d\n", len(bounds), l.ItemsPerRow)
	//fmt.Printf("maxWidth=%d, maxHeight=%d\n", l.TotalWidth, l.TotalHeight)
}
