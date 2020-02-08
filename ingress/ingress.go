package ingress

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

// indexLastSegment finds the index of the last \n after a complete set of segments
func indexLastSegment(buffer []byte) int {
	offset := len(buffer)
	lastid := ""

	for {
		index := bytes.LastIndex(buffer[:offset], []byte("\n"))
		if index < 0 {
			// If we didn't find a newline which split the segment, we look for the last newline
			// in the entire buffer.
			// If we can't find a newline anywhere in the buffer, we just return the entire buffer.
			findex := bytes.LastIndex(buffer, []byte("\n"))
			if findex < 0 {
				return len(buffer)
			}
			return findex
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
	ch := make(chan []byte, 20)

	go func() {
		offset := int64(0)
		buffer := make([]byte, bufsize)
		fileinfo, _ := f.Stat()
		filesize := fileinfo.Size()
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

			if int64(len(buffer)) < bufsize {
				ch <- buffer
				break
			}

			delim := int64(indexLastSegment(buffer))
			if delim < 0 {
				ch <- buffer
				break
			}

			offset += delim
			ch <- buffer[0:delim]
		}

		close(ch)
	}()

	return ch
}
