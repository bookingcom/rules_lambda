package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
)

var (
	input      = flag.String("input", "", "tar file input")
	output     = flag.String("output", "/dev/stdout", "tar file input")
	compress   = flag.Bool("enable-compression", false, "enable zip file compression")
	stderr     = log.New(os.Stderr, "", 0)
)

func main() {
	flag.Parse()

	if len(*input) == 0 {
		stderr.Println("Input is missing")
		os.Exit(1)
	}

	reader, err := os.Open(*input)

	if err != nil {
		stderr.Printf("Failed to open input file: %s", err.Error())
	}

	defer reader.Close()

	oFile, err := os.Create(*output)
	if err != nil {
		stderr.Fatalf("Failed to create output file: %s", err)
		os.Exit(1)
	}

	defer oFile.Close()

	outBuf := new(bytes.Buffer)

	writer := zip.NewWriter(outBuf)

	ibuf := make([]byte, 4096)

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		fih, err := zip.FileInfoHeader(hdr.FileInfo())
		if err != nil {
			stderr.Fatalf("Failed to create zip file info header from tar file info: %s", err.Error())
			os.Exit(2)
		}

		fih.Name = hdr.Name
		if *compress {
			fih.Method = zip.Deflate
		}

		f, err := writer.CreateHeader(fih)
		if err != nil {
			stderr.Fatalf("Failed while writing %s: %s", hdr.Name, err.Error())
			os.Exit(2)
		}

		for {
			read, err := tr.Read(ibuf)
			if read > 0 {
				_, err = f.Write(ibuf[:read])
				if err != nil {
					stderr.Fatalf("Failed to write to zip buffer: %s", err)
					os.Exit(2)
				}
				continue
			}
			if err == io.EOF {
				if read != 0 {
					stderr.Printf("EOF but read %d\n", read)
				}
				break
			}
			stderr.Fatalf("Something went wrong: %s", err.Error())
			os.Exit(2)
		}
	}

	err = writer.Close()
	if err != nil {
		stderr.Fatalf("Failed to close zip file: %s", err)
		os.Exit(1)
	}

	err = reader.Close()
	if err != nil {
		stderr.Fatalf("Failed to close reader: %s", err)
		os.Exit(1)
	}

	_, err = oFile.Write(outBuf.Bytes())
	if err != nil {
		stderr.Fatalf("Failed to write to output: %s", err)
		os.Exit(1)
	}
}
