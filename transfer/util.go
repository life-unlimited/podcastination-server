package transfer

import (
	"fmt"
	"io"
	"os"
)

// MoveFile moves a file from source path to destination path.
// Taken from https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b.
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		_ = inputFile.Close()
		return fmt.Errorf("open dest file: %s", err)
	}
	_, err = io.Copy(outputFile, inputFile)
	_ = inputFile.Close()
	if err != nil {
		_ = outputFile.Close()
		return fmt.Errorf("write dest file: %s", err)
	}
	// Delete source file.
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("remove source file: %s", err)
	}
	return nil
}
