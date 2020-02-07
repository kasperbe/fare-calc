package ingress

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func findDelimiter(buffer []byte) int {
	offset := len(buffer)
	lastid := ""

	for {
		index := bytes.LastIndex(buffer[:offset], []byte("\n"))
		if index < 0 {
			return len(buffer)
		}

		line := string(buffer[index:offset])
		row := strings.Split(line, ",")
		id := row[0]

		if lastid != "" && lastid != id {
			return offset
		}

		lastid = id
		offset -= len(line)
	}
}

// Read a file in chunks, when chunk is read find the last fare id
// traverse backwards until we find a different id and break the chunk at this point
// This is done to ensure we have chunks that doesn't break in the middle of a segment.
func Read(f *os.File, bufsize int64) chan []byte {
	offset := int64(0)
	ch := make(chan []byte, 10)
	buffer := make([]byte, bufsize)
	fileinfo, _ := f.Stat()
	filesize := fileinfo.Size()

	go func() {
		for {
			if bufsize > filesize-offset {
				// Make sure we don't try to read more bytes from the file than what's left.
				buffer = make([]byte, filesize-offset)
			}

			if filesize-offset == 0 {
				// Done
				break
			}

			_, err := f.ReadAt(buffer, offset)
			if err != nil {
				log.Fatal(fmt.Errorf("%w trying to read file", err))
			}

			delim := int64(findDelimiter(buffer))
			if delim < 0 {
				ch <- buffer
				break
			}

			offset += delim
			ch <- buffer[0:delim]
			buffer = make([]byte, bufsize)
		}

		close(ch)
	}()

	return ch
}
