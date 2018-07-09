package main

import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"os"
	_"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"github.com/jvdanker/go-test/manifest"
)

func main() {
	fmt.Println("test")

	//tasks.ResizeImages()
	//tasks.SliceImages()

    var total int
	dirs := walker.WalkDirectories("output/images/")
    for dir := range dirs {
        fmt.Println(dir)

        _, err := manifest.Read(dir + "/manifest.json")
        if (err == nil) {
            total++
        }
    }
    fmt.Printf("Total %d\n", total)

    itemsPerRow := int(math.Ceil(math.Sqrt(float64(total))))
    fmt.Println("itemsPerRow", itemsPerRow)
    fmt.Println()

    var i, z, x, y, maxY int
    dirs = walker.WalkDirectories("output/images/")
    for dir := range dirs {
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
            fmt.Printf("output/parts/0/%d/%d.png\n", nx, ny)

            oldName := fmt.Sprintf("../../../%s%s", dir, file.Name)
            newName := fmt.Sprintf("output/parts/%d/%d.png", nx, ny)

            fmt.Println("symlink", oldName, newName)

            if _, err := os.Stat(newName); err == nil {
                os.Remove(newName)
            }

            if _, err := os.Stat(fmt.Sprintf("output/parts/%d", nx)); os.IsNotExist(err) {
                os.MkdirAll(fmt.Sprintf("output/parts/%d", nx), os.ModePerm)
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

        if (i+1) % itemsPerRow == 0 {
            //fmt.Println("reset", maxY)
            x = 0
            y += maxY
            maxY = 0
        }

        i++

        fmt.Println()
    }
}

// maxzoom = Math.ceil( Math.log( (cw/iw > ch/ih ? iw/cw : ih/ch) ) / Math.log(2) );

/*

0 - 0 1 2 3
1 - 0 1 2 3
2 - 0 1 2 3
3 - 0 1 2 3

*/