package main

import (
	"fmt"
	"os"
	"github.com/jvdanker/go-test/util"
)

func createDirectories(in <- chan string) <- chan string {
    out := make(chan string)

    go func() {
        for dir := range in {
            out <- dir

            if _, err := os.Stat("output/" + dir); os.IsNotExist(err) {
                os.MkdirAll("output/" + dir, os.ModePerm)
            }
        }

        close(out)
    }()

    return out
}

func main() {
    ch := util.Walk()
    ch2 := createDirectories(ch)

    for dir := range ch2 {
        fmt.Println(dir)
    }
}

