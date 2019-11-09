# Introduction

imgdupfinder is a tool to find duplicate images. Its written in golang.
Currently it support following image formats 'png' and 'jpg/jpeg' .

# Features!

  - Recursive searching.

# Installation
## Usage
```sh
$ ./dupimgfinder -rootpath string
        Root folder full path.
    -recursive
    	Recursive search in subfolders. Default is false.
```
**Note:**  If '-rootpath' is not provided it will take the current directory as rootfolder.

## Eg
Finds the duplicate images in the current folder without looking into subfolders
```sh
$ ./dupimgfinder
```

Finds the duplicate images from root folder recursively checking the subfolders
```sh
$ ./dupimgfinder --rootpath=<rootpath> --recursive=true
```

# TODO
- Conurrent support
- Copy the duplicate files as softlinks to a duplicate folder.

