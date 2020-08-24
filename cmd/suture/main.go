package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	in   = flag.String("i", ".", "input folder")
	name = flag.String("n", "", "file name")
	out  = flag.String("o", ".", "output folder")
)

func main() {
	flag.Parse()

	var err error
	if *name == "" {
		log.Fatal("please assign file name by \"-n FILENAME \" ")
	}

	outFilePath := *out

	var outStat os.FileInfo
	if outStat, err = os.Stat(*out); err != nil {
		log.Fatalf("can not open output folder : %v", err)
	}
	if outStat.IsDir() {
		outFilePath = *out + string(filepath.Separator) + *name
	}

	fnames, err := listFile(*in, *name)
	if err != nil {
		log.Fatal(err)
	}

	var outFile *os.File
	if outFile, err = os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		log.Fatalf("write file : %s failed : %v", outFilePath, err)
	}
	defer outFile.Close()
	w := bufio.NewWriter(outFile)
	defer w.Flush()
	for _, fname := range fnames {
		inFPath := filepath.Join(*in, fname)
		var fct []byte
		if fct, err = ioutil.ReadFile(inFPath); err != nil {
			log.Fatalf("read input file : %s FAILED : %v", inFPath, err)
		}
		w.Write(fct)
	}
}

func listFile(in, name string) ([]string, error) {
	var err error
	var inStat os.FileInfo
	if inStat, err = os.Stat(in); err != nil {
		return nil, fmt.Errorf("can not open input file : %w", err)
	}
	if !inStat.IsDir() {
		return nil, fmt.Errorf("input is a FILE , not FOLDER")
	}
	var inFile *os.File
	if inFile, err = os.OpenFile(in, os.O_RDONLY, os.ModePerm); err != nil {
		return nil, fmt.Errorf("open file : %s failed : %w", in, err)
	}
	defer inFile.Close()
	var fiList []os.FileInfo
	if fiList, err = inFile.Readdir(0); err != nil {
		return nil, fmt.Errorf("read file info failed : %w", err)
	}

	result := make([]string, 0)
	for _, fi := range fiList {
		if fi.IsDir() {
			continue
		}
		if strings.HasPrefix(fi.Name(), name) {
			result = append(result, fi.Name())
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result, nil
}
