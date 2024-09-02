package main

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// content holds our static web server content.
//
//go:embed ebpf.tar.gz
var content embed.FS

func dumpFiles() error {
	f, err := content.Open("ebpf.tar.gz")
	if err != nil {
		panic(err)
	}
	return ExtractTarGz(f)
}

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)
	env := os.Getenv("SUDO_UID")

	i, err := strconv.Atoi(env)
	if err != nil {
		i = 0
	}

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
			err = os.Chown(header.Name, i, i) // TODO: quick patch to make files usable to $USER
			if err != nil {
				panic(err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				fmt.Errorf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

			err = os.Chown(header.Name, i, i) // TODO: quick patch to make files usable to $USER
			if err != nil {
				panic(err)
			}

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}

	}
	return nil
}
