package layout

import (
	"image"
)

type Layout interface {
	Layout([]image.Point)
}

type LayoutManager struct {
	ItemsPerRow int
	TotalWidth  uint32
	TotalHeight uint32
	Positions   []image.Point
}
