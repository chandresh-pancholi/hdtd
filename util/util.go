package util

import (
	"bytes"
	"compress/gzip"
	"log"
	"time"
)

func Compress(name string, data []byte) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	// Setting the Header fields is optional.
	zw.Name = name
	zw.ModTime = time.Now()

	_, err := zw.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
}
