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
)

func findImageFiles(rootPath *string, recursive bool) {

	// Create an empty store.
	store := hasher.NewImgsimStore()

	if recursive {
		err := filepath.Walk(*rootPath, func(path string, info os.FileInfo, err error) error {
			findDupImage(store, path)
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
				findDupImage(store, *rootPath+"/"+f.Name())
			}
		}
	}

}

func findDupImage(store *hasher.ImgsimStore, imgPath string) {

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
	flag.Parse()
	flag.PrintDefaults()

	if "" == *rootFolder {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		fmt.Println("Setting the current directory", wd, "as search rootpath.")
		*rootFolder = wd
	}

	//find out the duplicate image files
	findImageFiles(rootFolder, *recursive)
	fmt.Println("Total Number of duplicate files found ", noOfDuplicates)

}
