package protocol

import (
	"bufio"
	"fmt"
	"io"
)

func parseInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("empty input")
	}
	n := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid byte: %c", c)
		}
		n = n*10 + int(c-'0')
		if n < 0 {
			return 0, fmt.Errorf("overflow")
		}
	}
	return n, nil
}

// ParseArgs parses the custom wire protocol.
// Format: <len>,<data>,<len>,<data>,...
// Example: "3,SET,3,key,5,value" -> cmd="SET", key="key", value="value".
// The trailing newline is optional and ignored.
func ParseArgs(b []byte) (cmd, key, value []byte) {
	start := 0
	actualArg := 0

	for start < len(b) {
		end := start
		for end < len(b) && b[end] != ',' {
			end++
		}
		if end >= len(b) {
			break
		}
		argSize, err := parseInt(b[start:end])
		if err != nil || argSize < 0 {
			break
		}
		if end+1+argSize > len(b) {
			break
		}
		val := b[end+1 : end+1+argSize]
		start = end + 2 + argSize

		switch actualArg {
		case 0:
			cmd = val
		case 1:
			key = val
		case 2:
			value = val
		}
		actualArg++
	}
	return cmd, key, value
}

// PeekProtocol checks whether the next data in the reader follows the RESP protocol
// by looking at the first byte. RESP messages start with '*', '+', '-', ':', or '$'.
func PeekProtocol(reader *bufio.Reader) (bool, error) {
	b, err := reader.Peek(1)
	if err != nil {
		return false, err
	}
	switch b[0] {
	case '*', '+', '-', ':', '$':
		return true, nil
	}
	return false, nil
}

// ReadRESPCommand reads a RESP array command from the reader and extracts up to three arguments
// (command name, key, value). Returns ok=false on any parse error.
func ReadRESPCommand(reader *bufio.Reader) (cmd, key, value []byte, ok bool) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, nil, nil, false
	}
	if b != '*' {
		return nil, nil, nil, false
	}

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, nil, nil, false
	}
	l := len(line)
	if l < 2 || line[l-2] != '\r' {
		return nil, nil, nil, false
	}
	arrLen, err := parseInt(line[:l-2])
	if err != nil || arrLen < 1 || arrLen > 3 {
		return nil, nil, nil, false
	}

	args := make([][]byte, 0, arrLen)
	for i := 0; i < arrLen; i++ {
		b, err := reader.ReadByte()
		if err != nil || b != '$' {
			return nil, nil, nil, false
		}
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, nil, nil, false
		}
		l := len(line)
		if l < 2 || line[l-2] != '\r' {
			return nil, nil, nil, false
		}
		strLen, err := parseInt(line[:l-2])
		if err != nil || strLen < 0 {
			return nil, nil, nil, false
		}

		buf := make([]byte, strLen+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, nil, nil, false
		}
		if buf[strLen] != '\r' || buf[strLen+1] != '\n' {
			return nil, nil, nil, false
		}
		args = append(args, buf[:strLen])
	}

	if len(args) >= 1 {
		cmd = args[0]
	}
	if len(args) >= 2 {
		key = args[1]
	}
	if len(args) >= 3 {
		value = args[2]
	}
	return cmd, key, value, true
}
