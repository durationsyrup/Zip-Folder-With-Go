// Simple Little Folder Zipper
// Feel free to use my code :)
// Could use some cleaning up, sorry

// Usage "go run application.go -f ./ExampleFolder -o zipped" >> Output zipped.zip
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	f := flag.String("f", "", "-f (directory to zip)")
	o := flag.String("o", "", "-o (output file name)")
	flag.Parse()
	folder := string(*f)
	outputName := string(*o)

	usage := func() {
		fmt.Println("Version 0.1")

		fmt.Println("-f (directory to zip)")
		fmt.Println("-o (output file name)")
	}
	// To prevent some errors
	if folder != "" {
		if folder[len(folder)-1:] != "/" {
			folder = folder + "/"
		}
	} else {
		usage()
		return
	}
	if outputName != "" {
		if len(outputName) <= 4 {
			outputName = outputName + ".zip"
		} else {
			if outputName[len(outputName)-4:] != ".zip" {
				outputName = outputName + ".zip"
			}
		}
	} else {
		usage()
		return
	}

	zipDir(folder, outputName)
}
func zipDir(directory string, outcome string) {
	readFile := func(filename string) []byte {
		// I tested a lot of different ways to read files
		// and this is by far the best I've found to be good no matter the file size

		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		st, _ := f.Stat()
		s := st.Size()
		b := make([]byte, s)
		io.ReadFull(f, b)
		defer f.Close()
		return b
	}

	// simple but lengthy, size formatting func()

	size := func(n int64) string {
		num := int64(1024 - 32)

		if n < (num) {
			ans := int(n / num)
			Type := " byte"
			if ans > 1 {
				Type = Type + "s"
			}
			return strconv.Itoa(ans) + Type
		}
		if n < (num*num) && n > (num) {
			ans := int(n / (num))
			Type := " kilobyte"
			if ans > 1 {
				Type = Type + "s"
			}
			return strconv.Itoa(ans) + Type
		}
		if n > (num*num) && n < (num*num*num) {
			ans := int(n / (num * num))
			Type := " megabyte"
			if ans > 1 {
				Type = Type + "s"
			}
			return strconv.Itoa(ans) + Type
		}
		if n > (num * num * num) {
			ans := int(n / num)
			Type := " gigabyte"
			if ans > 1 {
				Type = Type + "s"
			}
			return strconv.Itoa(int(n/(num*num*num))) + Type
		}
		return ""
	}

	buff := new(bytes.Buffer)

	w := zip.NewWriter(buff)
	for filename, fileinfo := range crawlDir(directory) {
		f, err := w.Create(filename)
		if err != nil {
			log.Fatal(err)
		}

		data := readFile(directory + filename)

		f.Write(data)

		fmt.Println("File:", filename, "\nSize:", size(fileinfo.Size()))
	}

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Output file
	file, err := os.Create(outcome)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created", outcome, "/ Size:", size(int64(buff.Len())))
	file.Write(buff.Bytes())
	defer file.Close()
}
func crawlDir(directory string) map[string]os.FileInfo {
	// crawlDir() I made a live folder over tcp program with this function
	// This finds every file in a directory!
	// Probably a better way but idk
	dir, _ := os.Open(directory)
	defer dir.Close()
	files, _ := dir.Readdir(0)
	filenames := make(map[string]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			directory = directory + file.Name() + "/"
			morefiles := crawlDir(directory)
			directory = directory[:len(directory)-(len(file.Name())+1)]
			for d, found := range morefiles {
				filenames[file.Name()+"/"+d] = found
			}

		} else {
			filenames[file.Name()] = file
		}
	}
	return filenames
}
