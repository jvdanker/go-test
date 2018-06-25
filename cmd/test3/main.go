package main

import (
    "fmt"
    "sync"
)

func main() {
    c := make(chan int)

    var wg sync.WaitGroup

    wg.Add(1)
    go func(c <-chan int) {
        for i := range c {
            fmt.Println(i)
        }
        wg.Done()
    }(c)

    c <- 1
    c <- 1
    close(c)

    wg.Wait()
}