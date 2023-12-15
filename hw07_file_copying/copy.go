package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const BUFFERSIZE = 512

func Copy(fromPath, toPath string, offset, limit int64) error {
	inputFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer inputFile.Close()

	stat, err := inputFile.Stat()
	if err != nil {
		fmt.Println("Can't find the stat")
		return ErrUnsupportedFile
	}
	if offset > stat.Size() {
		fmt.Println("Offset is too big")
		return ErrOffsetExceedsFileSize
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		fmt.Println("Can't open or create output file")
		return ErrUnsupportedFile
	}
	defer outputFile.Close()

	var totalBytesRed int64
	progressBar := pb.StartNew(int(stat.Size() - offset))
	buf := make([]byte, BUFFERSIZE)
	i := offset
	for {
		bytesRed, readErr := readFile(inputFile, buf, i)
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return readErr
		}

		totalBytesRed += int64(bytesRed)
		buf = cutBuffer(buf, totalBytesRed, bytesRed, limit)

		writeErr := writeFile(outputFile, buf)

		if writeErr != nil {
			return writeErr
		}

		if bytesRed == 0 || errors.Is(readErr, io.EOF) {
			break
		}
		if limit > 0 && totalBytesRed >= limit {
			break
		}
		i += BUFFERSIZE
		progressBar.Add(BUFFERSIZE)
	}
	return nil
}

func readFile(inputFile *os.File, buffer []byte, offset int64) (int, error) {
	n, err := inputFile.ReadAt(buffer, offset)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("Error reading the file")
		return 0, ErrUnsupportedFile
	}
	return n, nil
}

func writeFile(outputFile *os.File, buffer []byte) error {
	_, err := outputFile.Write(buffer)
	if err != nil {
		fmt.Println("Error writing the file")
		return ErrUnsupportedFile
	}
	return nil
}

func cutBuffer(buffer []byte, totalBytesRed int64, bytesRed int, limit int64) []byte {
	if limit > 0 && totalBytesRed > limit {
		bytesToKeep := BUFFERSIZE - (totalBytesRed - limit)
		return buffer[:bytesToKeep]
	}
	if bytesRed < BUFFERSIZE {
		return buffer[:bytesRed]
	}
	return buffer
}
