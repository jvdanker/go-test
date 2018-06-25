package main

import (
    "fmt"
    "sync"
    "log"
    "os"
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


func filesWorker(id int, files <-chan util.File, wg *sync.WaitGroup, out chan <- util.File) {
    //fmt.Println("file worker", id, "started job")

    for file := range files {
        // fmt.Println(id, files, file)
        newName := "./output/" + file.Dir + "/" + file.Name + "_400x300.png"
        if _, err := os.Stat(newName); err == nil {
            out <- file
            continue
        }

        util.ResizeFile(file)
        out <- file
    }

    //fmt.Println("file worker", id, "finished job")

    wg.Done()
}

func dirWorker(id int, dirs <-chan string, wg *sync.WaitGroup) {
    for dir := range dirs {
        //fmt.Println("dir worker", id, "started job", dir)

        files := WalkFiles(dir)

        out := make(chan util.File)
        wg2 := sync.WaitGroup{}
        for w := 1; w <= 1; w++ {
            wg2.Add(1)
            go filesWorker(w, files, &wg2, out)
        }

        go func(id int, dir string) {
            for file := range out {
                fmt.Println(id, dir, file)
            }
        }(id, dir)

        wg2.Wait()
        close(out)

        //fmt.Println("dir worker", id, "finished job")
    }

    wg.Done()
}

func main() {
    fmt.Println("test")

    dirs := util.WalkDirectories()
    dirs = util.CreateDirectories(dirs)

    wg := sync.WaitGroup{}

    for w := 1; w <= 1; w++ {
        wg.Add(1)
        go dirWorker(w, dirs, &wg)
    }

    wg.Wait()
}