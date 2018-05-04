package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
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
	return dirTreeAcc(out, path, printFiles, 0)
}

func dirTreeAcc(out io.Writer, path string, printFiles bool, ident int) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	fileInfo, err := file.Readdir(-1)
	if err != nil {
		log.Fatal(err)
		return err
	}

	sort.Slice(fileInfo, func(i, j int) bool {
		return fileInfo[i].Name() < fileInfo[j].Name()
	})

	for i, f := range fileInfo {
		if f.IsDir() {
			fmt.Fprintln(out, formatStr(f.Name(), ident, i+1 == len(fileInfo)))
			dirTreeAcc(out, path+"/"+f.Name(), printFiles, ident+1)
		} else if printFiles {
			fmt.Fprintln(out, formatStr(f.Name(), ident, i+1 == len(fileInfo)))
		}
	}
	return nil
}

func formatStr(name string, ident int, isLast bool) string {
	leaf := "├───"
	if isLast {
		leaf = "└───"
	}
	return strings.Repeat("│\t", ident) + leaf + name
}
