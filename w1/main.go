package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type sortDir []os.FileInfo

func (a sortDir) Len() int {
	return len(a)
}
func (a sortDir) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a sortDir) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}

func onlyDirs(files []os.FileInfo) []os.FileInfo {
	var toPrint []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			toPrint = append(toPrint, file)
		}
	}
	return toPrint
}

func printableLine(file os.FileInfo, prefix string, last bool) (line string, newPrefix string) {
	var prefixStep string

	if last {
		prefixStep = "├───"
		newPrefix = prefix + "│\t"
	} else {
		prefixStep = "└───"
		newPrefix = prefix + "\t"
	}

	line = prefix + prefixStep + file.Name()

	if !file.IsDir() {
		fileSize := file.Size()
		if fileSize == 0 {
			line += fmt.Sprintf(" (empty)")
		} else {
			line += fmt.Sprintf(" (%db)", fileSize)
		}
	}
	line += "\n"

	return
}

func printDirTree(out io.Writer, path string, printFiles bool, prefix string) error {
	f, err := os.Open(path)
	if err != nil {
		closeError := f.Close()
		if closeError != nil {
			panic(closeError)
		}
		return err
	}

	files, err := f.Readdir(-1)
	closeError := f.Close()
	if closeError != nil {
		panic(closeError)
	}
	if err != nil {
		return err
	}

	if !printFiles {
		files = onlyDirs(files)
	}

	sort.Sort(sortDir(files))

	for key, file := range files {
		line, newPrefix := printableLine(file, prefix, key < len(files)-1)
		fmt.Fprint(out, line)

		if file.IsDir() {
			newPath := path + string(os.PathSeparator) + file.Name()
			err = printDirTree(out, newPath, printFiles, newPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := printDirTree(out, path, printFiles, "")
	if err != nil {
		return err
	}
	return nil
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
