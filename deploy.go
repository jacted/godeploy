package main

import (
	"io"
	"os"
	"path"
)

func walkFiles(srcPath string, info os.FileInfo, err error) error {
	if info.IsDir() {
		sftpClient.Mkdir(path.Join(tmpPath, srcPath))
	} else {
		destPath := path.Join(tmpPath, srcPath)
		err := uploadFile(srcPath, destPath)
		if err != nil {
			return nil
		}
	}

	return nil
}

func uploadFile(srcPath string, destPath string) error {
	b, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer b.Close()
	f, err := sftpClient.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, b); err != nil {
		return err
	}
	return nil
}
