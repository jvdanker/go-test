package tasks

import (
	"fmt"
	"github.com/jvdanker/go-test/manifest"
	"github.com/jvdanker/go-test/pipes"
	"github.com/jvdanker/go-test/util"
	"github.com/jvdanker/go-test/walker"
	"os"
)

func ResizeImages(input, output string) {
	dirs := walker.WalkDirectories(input)
	dirs = walker.CreateDirectories(output, dirs)

	containers := dirsToContainer(dirs)
	dirsToProcess := pipes.FanoutAndMerge(containers, 1, func(c *pipes.Container, worker int) {
		inputdir := c.Payload.(string)
		fmt.Printf("dirWorker=%v: dirWorker=%v\n", worker, inputdir)

		if _, err := os.Stat(fmt.Sprintf("%v/%v/manifest.json", output, inputdir)); err == nil {
			return
		}

		files := walker.WalkFiles(inputdir)
		c.Payload = files
	})

	containers = dirsToContainer2(dirsToProcess)
	pipes.FanoutAndWait(containers, 1, func(c *pipes.Container, worker int) {
		files := c.Payload.(<-chan util.File)
		filesContainer := filesToContainer(files)

		// FIXME race condition on output var
		// FIXME second fanout shouldn'be part of this method, this method
		// should return a channel with all files to process
		// FIXME come-up with a function that converts arbitrary data to a container
		pipes.FanoutAndMerge(filesContainer, 1, func(c *pipes.Container, worker int) {
			file := c.Payload.(util.File)

			fmt.Printf("fileWorker=%v: filesWorkers=%v\n", worker, file.Name)
			pi := util.ResizeFile(file, output)
			c.Payload = pi
		})

		//createManifestOfProcessedFiles(processedImages, worker, inputdir, output)
	})

}

func dirsToContainer(dirs <-chan string) <-chan pipes.Container {
	out := make(chan pipes.Container)

	go func() {
		for dir := range dirs {
			out <- pipes.Container{Payload: dir}
		}

		close(out)
	}()

	return out
}

func dirsToContainer2(dirs <-chan interface{}) <-chan pipes.Container {
	out := make(chan pipes.Container)

	go func() {
		for dir := range dirs {
			out <- pipes.Container{Payload: dir}
		}

		close(out)
	}()

	return out
}

func filesToContainer(in <-chan util.File) <-chan pipes.Container {
	out := make(chan pipes.Container)

	go func() {
		for data := range in {
			out <- pipes.Container{Payload: data}
		}

		close(out)
	}()

	return out
}

func createManifestOfProcessedFiles(processedImages <-chan interface{}, worker int, inputdir string, outputdir string) {
	var processedFiles []util.ProcessedImage

	for file := range processedImages {
		f := file.(util.ProcessedImage)
		processedFiles = append(processedFiles, f)
	}

	if len(processedFiles) > 0 {
		// create manifest file
		fmt.Printf("dirWorker=%v: createManifest=%v, count=%v\n", worker, inputdir, len(processedFiles))
		manifest.Create(processedFiles, inputdir, outputdir)
	}
}
