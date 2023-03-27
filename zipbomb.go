package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	iterations := 2
	var wg sync.WaitGroup

	placePayload()

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go unzipFile("payload.zip", fmt.Sprintf("%s_%d", "output", i), &wg)
	}
	wg.Wait()
}

func unzipFile(src, dst string, wg *sync.WaitGroup) {
	defer wg.Done()
	archive, err := zip.OpenReader(src)

	if err != nil {
		panic(err)
	}

	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

}

func placePayload() {
	payload_500Zip, err := payload_500ZipBytes()
	if err != nil {
		panic(err)
	}

	payloadFile, err := os.Create("payload.zip")
	if err != nil {
		panic(err)
	}

	defer payloadFile.Close()

	_, err = payloadFile.Write(payload_500Zip)
	if err != nil {
		panic(err)
	}
}
