package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func zipFolder(source, target string, includePathInZipFn func(string, bool) bool) error {
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
		isDir := info.IsDir()

		if isDir {
			relPath += "/"
		}

		if !includePathInZipFn(relPath, isDir) {
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

func extractZip(src, dest string) error {
	path, err := exec.LookPath("unzip")

	if err == nil {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}

		unzipCmd := exec.Command(path, src)
		unzipCmd.Dir = dest

		return unzipCmd.Run()
	} else {
		files, err := zip.OpenReader(src)
		if err != nil {
			return err
		}

		defer files.Close()

		for _, file := range files.File {
			err = func() error {
				readCloser, err := file.Open()
				if err != nil {
					return err
				}
				defer readCloser.Close()

				return extractZipArchiveFile(file, dest, readCloser)
			}()

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func extractZipArchiveFile(file *zip.File, dest string, input io.Reader) error {
	filePath := filepath.Join(dest, file.Name)
	fileInfo := file.FileInfo()

	if fileInfo.IsDir() {
		err := os.MkdirAll(filePath, fileInfo.Mode())
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			linkName, err := ioutil.ReadAll(input)
			if err != nil {
				return err
			}
			return os.Symlink(string(linkName), filePath)
		}

		fileCopy, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		_, err = io.Copy(fileCopy, input)
		if err != nil {
			return err
		}
	}

	return nil
}
