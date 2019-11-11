package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nr90/imgsim"
	"github.com/rif/imgdup2go/hasher"
)

var (
	imageFormat    = map[string]func(io.Reader) (image.Image, error){"jpg": jpeg.Decode, "jpeg": jpeg.Decode, "png": png.Decode, "gif": gif.Decode}
	noOfDuplicates = 0
	store          *hasher.ImgsimStore
	ch             = make(chan string, 10)
)

//findImageFiles function get the files in the folder and add it to the channel
func findImageFiles(rootPath *string, recursive bool) {

	if recursive {
		err := filepath.Walk(*rootPath, func(path string, info os.FileInfo, err error) error {
			ch <- path
			return nil
		})

		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	} else {
		files, err := ioutil.ReadDir(*rootPath)
		if err != nil {
			fmt.Println("Error", err.Error())
		}
		for _, f := range files {
			if !f.IsDir() {
				ch <- (*rootPath) + "/" + f.Name()
			}
		}
	}
	close(ch)

}

//findDupImage reads the file path from channel and query the store to check the match
func findDupImage(imgPath string, fileNameMatch bool) {

	ext := filepath.Ext(imgPath)
	if !strings.HasPrefix(ext, ".") {
		return
	}

	ext = ext[1:]
	if _, ok := imageFormat[ext]; !ok {
		return
	}

	file, err := os.Open(imgPath)
	if err != nil {
		return
	}

	decodeFunc, _ := imageFormat[ext]
	img, err := decodeFunc(file)
	if err != nil {
		return
	}

	hash := imgsim.AverageHash(img)
	matches := store.Query(hash)
	if matches != nil {

		if fileNameMatch {
			_, file := filepath.Split(fmt.Sprintf("%s", matches))
			_, file1 := filepath.Split(imgPath)

			if file != file1 {
				return
			}
		}

		fmt.Println(matches, " matches ", imgPath)
		noOfDuplicates++
	} else {
		store.Add(imgPath, hash)
	}
}

func main() {

	//Argument parsing
	rootFolder := flag.String("rootpath", "", "RootFolder fullpath")
	recursive := flag.Bool("recursive", false, "Recursive search in subfolders.")
	fileNameMatch := flag.Bool("filenamematch", true, "Search result should match file name.")

	flag.Parse()

	if "" == *rootFolder {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		fmt.Println("Setting the current directory", wd, "as search rootpath.")
		*rootFolder = wd
	}

	// Create an empty store.
	store = hasher.NewImgsimStore()

	//find out the duplicate image files
	go findImageFiles(rootFolder, *recursive)
	for path := range ch {
		findDupImage(path, *fileNameMatch)
	}
	fmt.Println("Total Number of duplicate files found ", noOfDuplicates)
}
