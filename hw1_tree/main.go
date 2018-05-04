package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

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
	return dirTreeAcc(out, path, printFiles, "")
}

func dirTreeAcc(out io.Writer, path string, printFiles bool, ident string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fileInfo, err := file.Readdir(-1)
	if err != nil {
		log.Fatal(err)
		return err
	}
	file.Close()

	sort.Slice(fileInfo, func(i, j int) bool {
		return fileInfo[i].Name() < fileInfo[j].Name()
	})

	if !printFiles {
		fileInfo = filter(fileInfo, func(f os.FileInfo) bool {
			return f.IsDir()
		})
	}

	for i, f := range fileInfo {
		isLast := i+1 == len(fileInfo)

		if f.IsDir() {
			fmt.Fprintln(out, formatStr(f.Name(), ident, isLast))
			if isLast {
				dirTreeAcc(out, path+"/"+f.Name(), printFiles, ident+"\t")
			} else {
				dirTreeAcc(out, path+"/"+f.Name(), printFiles, ident+"│\t")
			}

		} else {
			fmt.Fprintln(out, formatStrFile(f.Name(), f.Size(), ident, isLast))
		}
	}
	return nil
}

func formatStr(name string, ident string, isLast bool) string {
	leaf := "├───"
	if isLast {
		leaf = "└───"
	}
	return ident + leaf + name
}

func formatStrFile(name string, size int64, ident string, isLast bool) string {
	s := "empty"
	if size > 0 {
		s = strconv.FormatInt(size, 10) + "b"
	}
	return formatStr(name, ident, isLast) + " (" + s + ")"
}

func filter(vs []os.FileInfo, f func(os.FileInfo) bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
