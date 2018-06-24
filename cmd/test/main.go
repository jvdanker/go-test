package main

import (
	"fmt"
	"sync"
	"github.com/jvdanker/go-test/util"
)

func merge(cs ...chan util.ProcessedImage) <-chan util.ProcessedImage {
    var wg sync.WaitGroup
    out := make(chan util.ProcessedImage)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan util.ProcessedImage) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }

    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

func createManifest(cs ...chan util.ProcessedImage) []string {
    var wg sync.WaitGroup

    fmt.Println("create manifest")
    var result []string

    output := func(c <-chan util.ProcessedImage) {
        for n := range c {
            result = append(result, n.Original.Name)
        }
        wg.Done()
    }

    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    wg.Wait()
    return result
}

func spawnDigesters(files <-chan util.File) []chan util.ProcessedImage {
    fmt.Println("spawn")

    const numDigesters = 1

    c := make([]chan util.ProcessedImage, numDigesters)
    for i:=0; i<numDigesters; i++ {
        c[i] = make(chan util.ProcessedImage)
    }

    var wg sync.WaitGroup
    wg.Add(numDigesters)

    emitter := func(i int) {
        d := util.ResizeFiles(files)
        for file := range d {
            c[i] <- file
        }
        wg.Done()
    }

    for i := 0; i<numDigesters; i++ {
        go emitter(i)
    }

    go func() {
        wg.Wait()
        for i:=0; i<numDigesters; i++ {
            close(c[i])
        }
    }()

    return c
}

func main() {
    fmt.Println("Test")

    ch := util.WalkDirectories()
    ch2 := util.CreateDirectories(ch)
    ch3 := util.WalkFiles(ch2)

    fmt.Println(ch3)

    // channels := util.ResizeFiles(ch3)

    // channels := spawnDigesters(ch3)

    // result := createManifest(channels...)
    // fmt.Println(result)

//    for file := range merge(channels...) {
//        fmt.Printf("%v\n", file)
//    }
}

