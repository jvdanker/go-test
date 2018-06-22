package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
    dirsToProcess := make(map[string]int)

    currDir := ""
	err := filepath.Walk("images/", func (path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
            return err
        }

	    if info.IsDir() {
	        currDir = path
	        // fmt.Printf("visited dir: %q\n", path)
	    } else {
	        _, ok := dirsToProcess[currDir]
	        if ! ok {
	            fmt.Printf("added dir to process: %q\n", currDir)
	            dirsToProcess[currDir] = 1
	        }
	        // fmt.Printf("visited file: %q %q\n", currDir, info.Name())
	    }

        if len(dirsToProcess) > 10 {
            return filepath.SkipDir
        }

	    return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", "images/", err)
	}
}

