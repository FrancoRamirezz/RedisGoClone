package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// these are five resp identifiers. Each of these appear in the front of the wires encoding
// these are data types that appear in the front of the resp command line
const (
	RespString  byte = '+'
	RespError   byte = '-'
	RespInteger byte = ':'
	RespBulk    byte = '$'
	RespArray   byte = '*'
	RespNull    byte = '-'
)

// here we will define the struct to use in the serialization and the deserilaztion proccess
// which holds all commands and arguments we recieve from the client
// note the struct must be upper case to export
type RespValue struct {
	Type   byte // orginally i did Type string, but string takes more memory for the protocol
	String string
	Intger int
	Bulk   string
	Array  []RespValue
}

// this will take the serialized values we found
// curr points to the current byte  stream
// buffer from the byte stream
type Resp struct {
	Vaues []RespValue
	curr  int
	buff  []byte
}

//Note: As noted above, RESP is very simple. A server knows the type of the RESP command by the first byte of data.
// The server knows when a new line starts or the command is finished via the \r\n separator (CLRF).

// factory methods for constructing RESP values
// we will instante
func NewRespString(s string) RespValue {
	return RespValue{Type: RespString, String: s}
}

// to make a bulk string
func NewBulkString(bulk string) RespValue {
	return RespValue{Type: RespBulk, Bulk: bulk}
}
func NewIntger(i int) RespValue {
	return RespValue{Type: RespInteger, Intger: i}
}

// here we will reffer to the arr method
func NewArray(arr []RespValue) RespValue {
	return RespValue{Type: RespArray, Array: arr}
}

func NewNull() RespValue {
	return RespValue{Type: RespBulk}
}

func NewarrayNull() RespValue {
	return RespValue{Type: RespArray}
}

// any error will return from our set/get
func NewError(s string) RespValue {
	return RespValue{Type: RespError, String: s}
}

// Step two: the reader of all the information getting parsed
// bufio handles input and output
type RespReader struct {
	reader *bufio.Reader
}

// here we just insitnate the resp reader struct
// when we call the struct we reffer using the & and pass the struct
// the rd
func NewReader(rd io.Reader) *RespReader {
	return &RespReader{reader: bufio.NewReader(rd)}
}

// here we will create two functions that are essential for the parsing process
// readLine reads the line from the buffer
// readIntger read the intger from the buffer
func (r *RespReader) readLine() (line []byte, err error) {
	// this will keep runnning as we parse the reader
	for {
		read, err := r.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = append(read)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break // use a break statment since were running a infinite for loop
		}

	}
	return line[:len(line)-2], nil // use this to strip \r\n
}

// the function read intger
// atoi' means 'ASCII to integer' And 'itoa' is integer to ASCII
func (r *RespReader) readInteger() (RespValue, error) {
	line, err := r.readLine()
	if err != nil {
		return RespValue{}, err
	}
	n, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, fmt.Errorf("error resp %w", line)
	}
	return RespValue{Type: RespInteger, Intger: n}, nil
}

func (r *RespReader) readSimpleString() (RespValue, error) {
	// reffer to the readLine function
	line, err := r.readLine()
	if err != nil {
		return RespValue{}, nil
	}
	return RespValue{Type: RespString, String: string(line)}, nil
}

func (r *RespReader) readSimpleError() (RespValue, error) {
	line, err := r.readLine()
	if err != nil {
		return RespValue{}, nil
	}
	return RespValue{Type: RespError, String: string(line)}, nil
}

// consumCRLF reads two bytes and verifies they are '\r' and '\n'
func (r *RespReader) CommanLineCRLF() error {
	// the r point to the respReader
	start, err := r.reader.ReadByte()
	if err != nil {
		return err
	}
	// the r point to the RespReader and reffer inside the struct using the dot notation
	end, err := r.reader.ReadByte()
	if err != nil {
		return err
	}
	if start == '\r' || end == '\n' {
		return fmt.Errorf("resp: expected the bulk string, %q%q", start, end)
	}
	return nil
}

// parsing or deserilzation proccess
// we refer to the struct, name of the function
// this function will take the reader from the
func (r *RespReader) Read() (RespValue, error) {
	read, err := r.reader.ReadByte()
	if err != nil {
		return RespValue{}, err
	}

	switch read {
	// note each r.name of function
	case RespArray:
		return r.readArray() // these come from the readLine and readInteger
	case RespBulk:
		return r.readBulk()
	case RespString:
		return r.readSimpleString()
	case RespError:
		return r.readSimpleError()
	case RespInteger:
		return r.readInteger()
	default:
		return RespValue{}, nil

	}

}
