package main

import (
	"encoding/json"
	"fmt"
	"github.com/jvdanker/go-test/tasks"
	"github.com/jvdanker/go-test/walker"
	"net/http"
	"os"
)

func main() {
	//input := "/Volumes/Juan/Photos/Diversen/"
	//output := "/Volumes/App/output/"

	input := "test"
	output := "output"

	fmt.Printf("input=%v, output=%v\n", input, output)

	os.RemoveAll(output)
	os.MkdirAll(output, os.ModePerm)
	os.MkdirAll(output+"/images", os.ModePerm)

	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(output, dirs)
	pi := tasks.ResizeImages(dirs, output)
	manifests := tasks.CreateManifest(pi)
	mergedImages := tasks.MergeImages(manifests)
	tasks.SliceImages(mergedImages)

	result := tasks.CreateBottomLayer(output+"/images/", output+"/slices/", output+"/layers/")
	tasks.CreateZoomLayers(output + "/layers/")

	fmt.Println(result)
	fmt.Println("Listening :8080...")

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		j, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, fmt.Sprintf("%s", j))
	})

	fs := http.FileServer(http.Dir("./html/"))
	http.Handle("/", fs)

	fs2 := http.StripPrefix("/output/", http.FileServer(http.Dir("./output/")))
	http.Handle("/output/", fs2)

	http.ListenAndServe(":8080", nil)
}
