package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MoveFile copy file from source to destination directory. os.Rename gives error "invalid cross-device link"
func MoveFile(sourcePath, destPath string) error {

	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("file with the same name already exist")
	}

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	return nil
}

// Returns different file name if file with the same already exist
func ReturnDifNameIfDublicate(fileName, dir string) string {
	var newFileName = fileName
	fileNameWithoutExt := strings.Split(fileName, filepath.Ext(fileName))[0]
	ext := filepath.Ext(fileName)

	for i := 1; true; i++ {
		if _, err := os.Stat(dir + newFileName); err != nil {
			break
		}
		newFileName = fileNameWithoutExt + "_" + fmt.Sprint(i) + ext
	}

	return newFileName
}
