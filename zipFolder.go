package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func zipFolder(source, target string, includePathInZipFn func(string) bool) error {
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
		sourcePath := strings.Replace(source, "\\", "/", -1)
		relPath := strings.TrimPrefix(path, sourcePath)

		if relPath == "" {
			return nil
		}

		relPath = strings.TrimLeft(relPath, "/")
		if info.IsDir() {
			relPath += "/"
		}

		if !includePathInZipFn(relPath) {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

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
