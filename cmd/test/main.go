package main

import (
	"fmt"
	"github.com/jvdanker/go-test/util"
)


func main() {
    ch := util.WalkDirectories()
    ch2 := util.CreateDirectories(ch)
    ch3 := util.WalkFiles(ch2)
    ch4 := util.ResizeFiles(ch3)

    for file := range ch4 {
        fmt.Println(file.Dir + " " + file.Name)
    }
}

