package util

import (
    "fmt"
	"os"
	"path/filepath"
)

func Walk() <- chan string {
    out := make(chan string)

    go func() {
        filepath.Walk("images/", func (path string, info os.FileInfo, err error) error {
            if err != nil {
                fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
                return nil
            }

            if info.IsDir() {
                out <- path
            }

            return nil
        })

        close(out)
    }()

	return out
}

