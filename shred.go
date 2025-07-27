package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: path not provided\n")
		os.Exit(1)
	}

	err := Shred(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Shred(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	err = overwriteRandom(file)
	if err != nil {
		file.Close()
		return nil
	}

	file.Close()

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// overwrite the given file with random data
func overwriteRandom(file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var randomBuffer = make([]byte, 1024)

	for i := 0; i < 3; i++ {
		file.Seek(0, 0)
		bytesToWrite := fileSize

		for bytesToWrite > 0 {
			bytesWritingThisIteration := randomBuffer[:min(bytesToWrite, int64(len(randomBuffer)))]

			rand.Read(bytesWritingThisIteration)

			bytesWritten, err := file.Write(bytesWritingThisIteration)
			if err != nil {
				return err
			}

			bytesToWrite -= int64(bytesWritten)
		}

		file.Sync()
	}

	return nil
}
