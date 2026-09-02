package resp

import (
	"fmt"
	"io"
	"strconv"
)

func (r *RespReader) readBulk() (RespValue, error) {
	line, err := r.readLine()
	if err != nil {
		return RespValue{}, nil
	}
	// note make sure to compare to it as a string, we will convert from the string to byte
	bulklength, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, fmt.Errorf("invalid bulk length %f", err)
	}
	// null bulk string just return nil
	if bulklength == -1 {
		return NewNull(), nil
	} // can reffer to the trailing bulk length
	bulk := make([]byte, bulklength)
	if _, err := io.ReadFull(r.reader, bulk); err != nil {
		return RespValue{}, fmt.Errorf("resp")
	}
	// here we can refer to thje
	return RespValue{}, nil
}
