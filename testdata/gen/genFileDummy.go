/*
This package generates 4.8 MB sized dummy data.

To re-generate the "../dummy.bin" file, run the below from the root of the repo:
    go generate ./...
*/
package main

import (
	"bufio"
	"log"
	"os"
)

const sizeByte = int64(1e6 * 5)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("file path is missing")
	}

	pathFile := os.Args[1]

	f, err := os.Create(pathFile)
	if err != nil {
		log.Fatal("failed to create file:", err)
	}
	defer f.Close()

	fw := bufio.NewWriter(f)

	for i := int64(0); i < sizeByte; i++ {
		if _, err = fw.Write([]byte{byte(i % 254)}); err != nil {
			log.Fatal("failed to write")
		}
	}

	if err = fw.Flush(); err != nil {
		log.Fatal(err)
	}
}
