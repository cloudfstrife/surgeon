package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	in  = flag.String("i", "", "input file")
	out = flag.String("o", ".", "output folder")
	n   = flag.Int("n", 1024*1024, "splited size")
)

func main() {
	flag.Parse()

	var err error

	if *in == "" {
		log.Fatal("please assign input file by \"-i FILENAME \" ")
	}
	var inStat os.FileInfo
	if inStat, err = os.Stat(*in); err != nil {
		log.Fatalf("can not open input file : %v", err)
	}
	if inStat.IsDir() {
		log.Fatal("please assign input FILE , not FOLDER")
	}

	var outStat os.FileInfo
	if outStat, err = os.Stat(*out); err != nil {
		log.Fatalf("can not open output folder : %v", err)
	}
	if !outStat.IsDir() {
		log.Fatal("please assign output Folder , not FILE")
	}

	var inFile *os.File
	if inFile, err = os.Open(*in); err != nil {
		log.Fatalf("can not open input file : %v", err)
	}
	defer func(f *os.File) {
		if f.Close() != nil {
			log.Fatalf("close input file failed : %v", err)
		}
	}(inFile)

	reader := bufio.NewReader(inFile)
	i := 0
	for {
		i++
		var (
			ct = make([]byte, *n)
			c  int
		)
		if c, err = reader.Read(ct); err != nil {
			if err == io.EOF {
				log.Println("EOF")
				break
			}
			log.Printf("read input file Error : %v", err)
		}
		outFilePath := filepath.Join(*out, inStat.Name()+"."+fmt.Sprintf("%05d", i))
		var outFile *os.File
		if outFile, err = os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
			log.Fatalf("write file : %s failed : %v", outFilePath, err)
		}
		if wc, err := outFile.Write(ct[0:c]); err != nil || wc != c {
			if err := outFile.Close(); err != nil {
				log.Printf("close output file failed : %v", err)
			}
			log.Fatalf("write file failed : %v", err)
		}
	}
}
