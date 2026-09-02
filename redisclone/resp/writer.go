package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// in this file we Encode the value structs into RESP wire format to connect them to a tcp server
type RespWriter struct {
	writer *bufio.Writer
}

func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{writer: bufio.NewWriter(w)}
}
func (rw *RespWriter) Write(value RespValue) error {
	// we Marshal to conver t the struct into JSON bytes
	//from this states what commands the client sent and how we response with RESP and the writer
	_, err := rw.writer.Write(value.Marshal())
	return err
}

// here the Marshal encodes a Value into its RESP byte repersentation
// this will call specific method for each type based on the value type
func (v RespValue) Marshal() []byte {
	switch v.Type {
	case RespString:
		// +OK\r\n
		return []byte("+" + v.String + "\r\n")

	case RespError:
		// -ERR message\r\n
		return []byte("-" + v.String + "\r\n")

	case RespInteger:
		// :42\r\n
		return []byte(":" + strconv.Itoa(v.Intger) + "\r\n")

	case RespBulk:
		if v.Bulk == "" {
			// $-1\r\n — null bulk string (key not found)
			return []byte("$-1\r\n")
		}
		// $5\r\nhello\r\n
		// The length prefix is what makes bulk strings binary-safe.
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Bulk), v.Bulk))

	case RespArray:
		if v.Array == nil {
			// *-1\r\n — null array
			return []byte("*-1\r\n")
		}
		// *3\r\n  then each element encoded recursively
		// Pre-allocate a reasonable buffer to avoid repeated allocations
		out := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
		for _, elem := range v.Array {
			out = append(out, elem.Marshal()...)
		}
		return out

	default:
		//
		return []byte(fmt.Sprintf("-ERR internal: unknown RESP type %q\r\n", v.Type))
	}
}
