package layout

import (
    "image"
    "github.com/jvdanker/go-test/manifest"
)

type Layout interface {
    Layout()
}

type LayoutManager struct {
    manifest.ManifestFile
    Positions []image.Point
}
