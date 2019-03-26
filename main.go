package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type folderInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	sl := make([]string, 1)
	result, err := walk(path, false, printFiles, sl)

	if err != nil {
		return err
	}

	fmt.Fprintf(out, result)
	return nil
}

func walk(path string, last, printFiles bool, preffixSlice []string) (string, error) {
	var level int
	var str string
	filesSorted := make([]folderInfo, 0)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, v := range files {
		t := folderInfo{v.Name(), v.Size(), v.IsDir()}
		if printFiles {
			filesSorted = append(filesSorted, t)
		} else {
			if v.IsDir() {
				filesSorted = append(filesSorted, t)
			}
		}
	}

	sort.Slice(filesSorted, func(i, j int) bool { return filesSorted[i].Name < filesSorted[j].Name })

	if path == "." {
		level = 0
	} else {
		level = len(strings.Split(path, string(os.PathSeparator)))
	}

	for i, f := range filesSorted {
		var tempString string

		if len(filesSorted) == 1 || i == len(filesSorted)-1 {
			last = true
		} else {
			last = false
		}

		if level >= len(preffixSlice) {
			if last {
				preffixSlice = append(preffixSlice, "\t")
			} else {
				preffixSlice = append(preffixSlice, "│\t")
			}
		} else {
			if last {
				preffixSlice[level] = "\t"
			} else {
				preffixSlice[level] = "│\t"
			}
		}

		if level > 0 {
			for i := 0; i < level; i++ {
				tempString += preffixSlice[i]
			}
		}

		if last {
			tempString += "└"
		} else {
			tempString += "├"
		}

		tempString += "───"

		if f.IsDir {
			tempString += f.Name + "\n"
		} else {
			var size string
			if f.Size == 0 {
				size = "empty"
			} else {
				size = strconv.Itoa(int(f.Size))
				size += "b"
			}
			tempString += f.Name + " (" + size + ")\n"
		}

		if f.IsDir {
			partString, err := walk(filepath.Join(path, f.Name), last, printFiles, preffixSlice)
			if err != nil {
				return "", err
			}
			tempString += partString
		}

		str += tempString
	}

	return str, nil
}
