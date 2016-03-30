package main

import (
	"archive/zip"
	"fmt"
    "io"
    "os"
	"path/filepath"
	"strings"
)

func zipFolder(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		path = strings.Replace(path, "\\", "/", -1)
		relPath := strings.TrimPrefix(path, source)
		
		if (relPath == "") {
			return nil
		}
		
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		fmt.Printf("\tAdding %s\n", header.Name)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		_, err = io.Copy(writer, file)
		
		return err
	})

	return err
}