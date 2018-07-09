package tasks

import (
	"fmt"
	"image"
	"github.com/jvdanker/go-test/walker"
	"github.com/jvdanker/go-test/util"
)

func SliceImages() {
    dirs := walker.WalkDirectories("output/images/")
	for dir := range dirs {
	    fmt.Println(dir)

        img, err := util.DecodeImage(dir + "/result.png")
	    if err != nil {
	        fmt.Println(err)
	        continue
	    }

	    //var x,y int
	    bounds := img.Bounds()
        w, h := bounds.Max.X, bounds.Max.Y
        fmt.Println("w,h=", w, h)

        var x,y,i,j int
        i = 0
        j = 0

        for y<h {
            for x=0;  x<w; x+=256 {
                r := image.Rect(x, y, x + 256, y + 256)
                fmt.Println(r)

                // todo resize image to 256x256 if less than this
                
                if img2, ok := img.(*image.NRGBA); ok {
                    sub := img2.SubImage(r)
                    util.CreateImage(fmt.Sprintf("%s/sub-%d-%d.png", dir, i, j), sub)
                }
                if img2, ok := img.(*image.RGBA); ok {
                    sub := img2.SubImage(r)
                    util.CreateImage(fmt.Sprintf("%s/sub-%d-%d.png", dir, i, j), sub)
                }

                i++
            }

            y += 256
            j++
            i = 0
        }

        fmt.Println()
	}
}
