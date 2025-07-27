package main

import (
	"crypto/rand"
	"os"
	"path"
	"slices"
	"testing"
)

func TestOverwriteRandom(t *testing.T) {
	const TestFileSize = 1025
	testFilePath := path.Join(t.TempDir(), "testfile")

	testData := make([]byte, TestFileSize)
	rand.Read(testData)

	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer testFile.Close()

	_, err = testFile.Write(testData)
	if err != nil {
		t.Fatal(err)
	}

	err = overwriteRandom(testFile)
	if err != nil {
		t.Error(err)
	}

	fileInfo, err := testFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if fileInfo.Size() != TestFileSize {
		t.Error("file size changed")
	}

	testFile.Seek(0, 0)
	// because Shred randomizes the file there is a chance that a chunk is
	// unchanged but with a chunk size of 32 the probability of that happening
	// is 8.6e-78, although the probability of a false positive is higher because
	// any chunk could match
	const ChunkSize = 32
	testFileChunk := make([]byte, ChunkSize)
	for chunk := range slices.Chunk(testData, ChunkSize) {
		if len(chunk) < ChunkSize {
			// checking this chunk would make the probability of a false
			// positive too high so just ignore it, close enough for this
			continue
		}

		testFile.Read(testFileChunk[:len(chunk)])

		allMatching := true
		for i := 0; i < len(chunk); i++ {
			if chunk[i] != testFileChunk[i] {
				allMatching = false
				break
			}
		}

		if allMatching {
			t.Error("part of the file was unchanged")
		}
	}
}
