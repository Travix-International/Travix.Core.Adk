package appix

import (
	archiveZip "archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type filePickerFunc func(path string) bool

func zipFolder(source, target string, includePathInZipFn filePickerFunc) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := archiveZip.NewWriter(zipfile)
	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		isDir := info.IsDir()

		if relPath == "" {
			return nil
		}

		relPath = strings.TrimLeft(relPath, "/")

		if !includePathInZipFn(relPath) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}

		if isDir {
			return nil
		}

		header, err := archiveZip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.Replace(relPath, string(os.PathSeparator), "/", -1)
		header.Method = archiveZip.Deflate
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)

		return err
	})
}

func extractZip(src, dest string) error {
	reader, err := archiveZip.OpenReader(src)

	if err != nil {
		return err
	}

	defer reader.Close()

	for _, f := range reader.Reader.File {

		zipped, err := f.Open()
		if err != nil {
			return err
		}

		defer zipped.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0777)
		} else {
			dirPath := filepath.Dir(path)
			os.MkdirAll(dirPath, 0777)

			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

			if err != nil {
				return err
			}

			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {
				return err
			}
		}
	}

	return nil
}
