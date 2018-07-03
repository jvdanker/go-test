package layout

import (
    "image"
)

type Layout interface {
    Layout([]image.Point)
}

type LayoutManager struct {
    ItemsPerRow int
    TotalWidth int
    TotalHeight int
    Positions []image.Point
}
