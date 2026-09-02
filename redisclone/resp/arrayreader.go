package resp

import (
	"fmt"
	"strconv"
)

// where going to call the respReader. Note to export the struct
// note the parameter will take a arugment of the struct\
func (r *RespReader) readArray() (RespValue, error) {
	// the var will look at the struct and the elements
	// here we will look at the two componets: the readintger, it will take the len of the
	read, err := r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	// *-1\r\n -> null array
	readarray, err := strconv.Atoi(string(read))
	if err != nil {
		return RespValue{}, fmt.Errorf("invalid array %w ", err)
	}
	if readarray == -1 {
		return NewarrayNull(), nil
	}
	if readarray < 0 {
		return RespValue{}, fmt.Errorf("invalid array length %d", read)
	}
	// for each line, we will parse adn read the value
	values := make([]RespValue, readarray)
	for i := 0; i <= readarray; i++ {
		val, err := r.Read()
		if err != nil {
			return RespValue{}, fmt.Errorf("resp: element %d: %w", i, err)
		}
		// i orginally did values[i]
		values = append(values, val)

	}
	// for each call: we will use a recursion for it
	return NewArray(values), nil

}
