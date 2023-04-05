package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
)

var (
	input    = flag.String("input", "", "tar file input")
	output   = flag.String("output", "/dev/stdout", "tar file input")
	compress = flag.Bool("enable-compression", false, "enable zip file compression")
	stderr   = log.New(os.Stderr, "", 0)
)

func main() {
	flag.Parse()

	if len(*input) == 0 {
		stderr.Println("Input is missing")
		os.Exit(1)
	}

	file, err := os.Open(*input)

	if err != nil {
		stderr.Printf("Failed to open input file: %s", err.Error())
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	uncompressed, err := gzip.NewReader(reader)
	if err != nil {
		stderr.Fatalf("Failed to start decompressor for input file: %s", err.Error())
		os.Exit(1)
	}

	defer uncompressed.Close()

	oFile, err := os.Create(*output)
	if err != nil {
		stderr.Fatalf("Failed to create output file: %s", err.Error())
		os.Exit(1)
	}

	defer oFile.Close()
	ibuf := make([]byte, 4096)
	for {
		count, err := uncompressed.Read(ibuf)
		if count > 0 {
			_, err := oFile.Write(ibuf[:count])
			if err != nil {
				stderr.Fatalf("Failed to write batch: %s", err.Error())
				os.Exit(1)
			}
			continue
		}

		if err == io.EOF {
			break
		}
		stderr.Fatalf("Something went wrong: %s", err.Error())
		os.Exit(1)
	}
}
