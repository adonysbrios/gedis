package protocol

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCmd   string
		wantKey   string
		wantValue string
	}{
		{
			name:      "SET command",
			input:     "3,SET,3,key,5,value\n",
			wantCmd:   "SET",
			wantKey:   "key",
			wantValue: "value",
		},
		{
			name:      "GET command",
			input:     "3,GET,3,key\n",
			wantCmd:   "GET",
			wantKey:   "key",
			wantValue: "",
		},
		{
			name:      "DEL command",
			input:     "3,DEL,3,key\n",
			wantCmd:   "DEL",
			wantKey:   "key",
			wantValue: "",
		},
		{
			name:      "PING command",
			input:     "4,PING\n",
			wantCmd:   "PING",
			wantKey:   "",
			wantValue: "",
		},
		{
			name:      "empty line",
			input:     "\n",
			wantCmd:   "",
			wantKey:   "",
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, key, value := ParseArgs([]byte(tt.input))
			if string(cmd) != tt.wantCmd {
				t.Errorf("cmd = %q, want %q", string(cmd), tt.wantCmd)
			}
			if string(key) != tt.wantKey {
				t.Errorf("key = %q, want %q", string(key), tt.wantKey)
			}
			if string(value) != tt.wantValue {
				t.Errorf("value = %q, want %q", string(value), tt.wantValue)
			}
		})
	}
}

func TestParseArgsMalformed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCmd  string
		wantKey  string
		wantVal  string
	}{
		{"missing length", "SET,3,key\n", "", "", ""},
		{"invalid length", "abc,SET\n", "", "", ""},
		{"truncated value", "3,SET,3,key,5,va", "SET", "key", ""},
		{"negative length", "-1,SET\n", "", "", ""},
		{"no commas", "SET", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, key, value := ParseArgs([]byte(tt.input))
			if string(cmd) != tt.wantCmd {
				t.Errorf("cmd = %q, want %q", string(cmd), tt.wantCmd)
			}
			if string(key) != tt.wantKey {
				t.Errorf("key = %q, want %q", string(key), tt.wantKey)
			}
			if string(value) != tt.wantVal {
				t.Errorf("value = %q, want %q", string(value), tt.wantVal)
			}
		})
	}
}

func TestReadRESPCommand(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCmd   string
		wantKey   string
		wantValue string
		wantOK    bool
	}{
		{
			name:      "SET command",
			input:     "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			wantCmd:   "SET",
			wantKey:   "key",
			wantValue: "value",
			wantOK:    true,
		},
		{
			name:      "GET command",
			input:     "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			wantCmd:   "GET",
			wantKey:   "key",
			wantValue: "",
			wantOK:    true,
		},
		{
			name:      "DEL command",
			input:     "*2\r\n$3\r\nDEL\r\n$3\r\nkey\r\n",
			wantCmd:   "DEL",
			wantKey:   "key",
			wantValue: "",
			wantOK:    true,
		},
		{
			name:      "PING command",
			input:     "*1\r\n$4\r\nPING\r\n",
			wantCmd:   "PING",
			wantKey:   "",
			wantValue: "",
			wantOK:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			cmd, key, value, ok := ReadRESPCommand(reader)
			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}
			if string(cmd) != tt.wantCmd {
				t.Errorf("cmd = %q, want %q", string(cmd), tt.wantCmd)
			}
			if string(key) != tt.wantKey {
				t.Errorf("key = %q, want %q", string(key), tt.wantKey)
			}
			if string(value) != tt.wantValue {
				t.Errorf("value = %q, want %q", string(value), tt.wantValue)
			}
		})
	}
}

func TestReadRESPCommandMalformed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"wrong prefix", "+OK\r\n"},
		{"not an array", "3,SET,3,key\n"},
		{"truncated", "*3\r\n$3\r\nSET\r\n"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			_, _, _, ok := ReadRESPCommand(reader)
			if ok {
				t.Errorf("expected ok=false for malformed input")
			}
		})
	}
}

func TestPeekProtocol(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantRESP bool
	}{
		{"RESP array", "*3\r\n", true},
		{"RESP simple string", "+OK\r\n", true},
		{"RESP error", "-ERR\r\n", true},
		{"RESP integer", ":1\r\n", true},
		{"RESP bulk", "$5\r\n", true},
		{"custom protocol", "3,SET,3,key\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			got, err := PeekProtocol(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantRESP {
				t.Errorf("PeekProtocol = %v, want %v", got, tt.wantRESP)
			}
		})
	}
}
