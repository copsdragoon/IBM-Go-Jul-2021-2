package main

import (
	"fmt"
	"io"
	"os"
)

type AlphaReader struct {
	str   string
	count int64
}

const (
	SeekStart   = 0 // seek relative to the origin of the file
	SeekCurrent = 1 // seek relative to the current offset
	SeekEnd     = 2 // seek relative to the end
)

func (alphaReader *AlphaReader) Read(p []byte) (n int, err error) {
	count := 0
	l, _ := alphaReader.Seek(0, SeekCurrent)

	for idx := int(l); idx < len(alphaReader.str); idx++ {
		if (alphaReader.str[idx] >= 'A' && alphaReader.str[idx] <= 'Z') ||
			(alphaReader.str[idx] >= 'a' && alphaReader.str[idx] <= 'z') {
			p[count] = alphaReader.str[idx]
			count++
			if count >= len(p) {
				alphaReader.Seek(int64(idx), SeekStart)
				return count, nil
			}
		}
	}
	alphaReader.Seek(0, SeekEnd)
	return count, io.EOF
}

func (alphaReader *AlphaReader) Seek(offset int64, whence int) (int64, error) {
	length := int64(0)
	maxLen := int64(len(alphaReader.str))
	switch whence {
	case SeekStart:
		length = offset
	case SeekCurrent:
		if offset == 0 {
			return alphaReader.count, io.EOF
		} else {
			length = alphaReader.count + offset
		}
	case SeekEnd:
		length = maxLen - offset
	}
	if length < maxLen {
		alphaReader.count = length
	} else {
		alphaReader.count, length = maxLen, maxLen
	}
	return length, io.EOF
}

func copy(writer io.Writer, reader io.ReadSeeker) {
	totalBytesRead := 0
	for {
		buffer := make([]byte, 10)
		bytesRead, err := reader.Read(buffer)

		totalBytesRead += bytesRead
		// fmt.Println("\nTotal bytes read ", totalBytesRead)
		if bytesRead > 0 {
			writer.Write(buffer[:bytesRead])
		}
		if err == io.EOF {
			fmt.Println()
			break
		}
	}
}

func main() {
	alphaReader := AlphaReader{"this is a sample string with 1234 numbers and *&^%^$ special characters", 0}
	copy(os.Stdout, &alphaReader)
}
