package main

import (
    "fmt"
    "sync"
    "log"
    "io/ioutil"
    "github.com/jvdanker/go-test/util"
)

func WalkFiles(dir string) <- chan util.File {
    out := make(chan util.File)

    go func() {
        // fmt.Printf("Walkfiles, dir = %v\n", dir)

        files, err := ioutil.ReadDir(dir)
        if err != nil {
            log.Fatal(err)
        }

        for _, f := range files {
            if f.IsDir() {
                continue
            }

            file := util.File{
                Dir: dir,
                Name: f.Name(),
            }

            // fmt.Printf("Walk file: %v\n", file)
            out <- file
        }

        close(out)
    }()

    return out
}

func main() {
    fmt.Println("test")

    var wg sync.WaitGroup
    var i int
    var mutex = &sync.Mutex{}

    ch := util.WalkDirectories()
    for dir := range ch {
        mutex.Lock()
        add := i < 2
        mutex.Unlock()

        if add {
            wg.Add(1)

            mutex.Lock()
            i++
            mutex.Unlock()

            ch := WalkFiles(dir)

            go func(ch <- chan util.File, i *int) {
                for file := range ch {
                    fmt.Println(ch, file)
                }
                wg.Done()

                mutex.Lock()
                *i--
                mutex.Unlock()
            }(ch, &i)
        }
    }

    wg.Wait()
}