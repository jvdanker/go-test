package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/walker"
	"math"
	"os"
	"strconv"
	"strings"
)

func CreateBottomLayer() {
	items, w, h := getMaxBounds()
	maxzoom := int(math.Ceil(math.Log(math.Max(w, h)/256) / math.Log(2)))
	fmt.Printf("max zoom = %v\n", maxzoom)

	itemsPerRow := int(math.Ceil(math.Sqrt(float64(items))))
	fmt.Println("itemsPerRow", itemsPerRow)
	fmt.Println()

	var i, z, x, y, maxY int
	dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}

		fmt.Println("dir", dir, z, x, y)

		var imgX int

		files := walker.WalkSlicedFiles(dir)
		for file := range files {
			s := strings.TrimPrefix(file.Name, "sub-")
			s = strings.TrimSuffix(s, ".png")
			parts := strings.Split(s, "-")

			fx, _ := strconv.Atoi(parts[0])
			fy, _ := strconv.Atoi(parts[1])

			nx := x + fx
			ny := y + fy

			fmt.Printf("file=%v loc=%d, %d newloc=%d, %d\n", file, fx, fy, nx, ny)
			fmt.Printf("output/parts/%d/%d/%d.png\n", maxzoom, nx, ny)

			oldName := fmt.Sprintf("../../../../%s%s", dir, file.Name)
			newName := fmt.Sprintf("output/parts/%d/%d/%d.png", maxzoom, nx, ny)

			fmt.Println("symlink", oldName, newName)

			if _, err := os.Stat(newName); err == nil {
				fmt.Println(" exists ", newName)
				os.Remove(newName)
			}

			if _, err := os.Stat(fmt.Sprintf("output/parts/%d/%d", maxzoom, nx)); os.IsNotExist(err) {
				os.MkdirAll(fmt.Sprintf("output/parts/%d/%d", maxzoom, nx), os.ModePerm)
			}

			err := os.Symlink(oldName, newName)
			if err != nil {
				panic(err)
			}

			fx++
			fy++

			if fx > imgX {
				imgX = fx
			}
			if fy > maxY {
				maxY = fy
			}
		}

		x += imgX

		if (i+1)%itemsPerRow == 0 {
			//fmt.Println("reset", maxY)
			x = 0
			y += maxY
			maxY = 0
		}

		i++

		fmt.Println()
	}
}

func getMaxBounds() (int, float64, float64) {
	var total int
	var tx, ty float64

	dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
		fmt.Println(dir)

		m, err := manifest.Read(dir + "/manifest.json")
		if err == nil {
			total++

			tx += float64(m.Layout.TotalWidth)
			ty += float64(m.Layout.TotalHeight)
		}
	}
	fmt.Printf("Total %d, w=%f, h=%f\n", total, tx, ty)

	return total, tx, ty
}
