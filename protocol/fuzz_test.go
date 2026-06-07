package protocol

import (
	"bufio"
	"strings"
	"testing"
)

func FuzzParseArgs(f *testing.F) {
	seeds := []string{
		"3,SET,3,key,5,value\n",
		"3,GET,3,key\n",
		"3,DEL,3,key\n",
		"4,PING\n",
		"\n",
		"3,SET,1,a,1,b\n",
		"abc,SET\n",
		"-1,SET\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, input string) {
		cmd, key, value := ParseArgs([]byte(input))
		// Just make sure it doesn't panic or return garbage
		_ = cmd
		_ = key
		_ = value
	})
}

func FuzzParseRESP(f *testing.F) {
	seeds := []string{
		"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
		"*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
		"*1\r\n$4\r\nPING\r\n",
		"",
		"+OK\r\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, input string) {
		reader := bufio.NewReader(strings.NewReader(input))
		cmd, key, value, ok := ReadRESPCommand(reader)
		_ = cmd
		_ = key
		_ = value
		_ = ok
	})
}
