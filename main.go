package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adbr.xx/gedis/commands"
	"adbr.xx/gedis/database"
	"adbr.xx/gedis/protocol"
)

const maxConnections = 1000
const serverPassword = ""
const connectionTimeout = 5 * time.Minute

var connSem = make(chan struct{}, maxConnections)

type clientHandler struct {
	conn          net.Conn
	reader        *bufio.Reader
	resp          bool
	authenticated bool
}

func (h *clientHandler) writeResponse(data interface{}) {
	if h.resp {
		switch v := data.(type) {
		case string:
			h.conn.Write([]byte("+" + v + "\r\n"))
		case []byte:
			h.conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)))
		case nil:
			h.conn.Write([]byte("$-1\r\n"))
		}
	} else {
		switch v := data.(type) {
		case string:
			h.conn.Write([]byte(v + "\n"))
		case []byte:
			h.conn.Write(append(v, '\n'))
		case nil:
			h.conn.Write([]byte("KEY_NOT_FOUND\n"))
		}
	}
}

func (h *clientHandler) writeError(msg string) {
	if h.resp {
		h.conn.Write([]byte("-ERR " + msg + "\r\n"))
	} else {
		h.conn.Write([]byte("ERR " + msg + "\n"))
	}
}

func (h *clientHandler) handle() {
	defer h.conn.Close()
	h.conn.SetDeadline(time.Now().Add(connectionTimeout))

	for {
		resp, err := protocol.PeekProtocol(h.reader)
		if err != nil {
			return
		}
		h.resp = resp

		var cmd, key, value []byte
		if resp {
			cmd, key, value, _ = protocol.ReadRESPCommand(h.reader)
		} else {
			line, err := h.reader.ReadBytes('\n')
			if err != nil {
				return
			}
			cmd, key, value = protocol.ParseArgs(line)
		}
		h.conn.SetDeadline(time.Now().Add(connectionTimeout))

		if len(cmd) == 0 {
			continue
		}

		if serverPassword != "" && !h.authenticated && string(cmd) != "AUTH" {
			h.writeError("authentication required")
			continue
		}

		cmdStr := string(cmd)
		switch cmdStr {
		case "PING":
			h.writeResponse("PONG")
		case "SET":
			if len(key) == 0 || len(value) == 0 {
				h.writeError("wrong number of arguments for SET")
				continue
			}
			commands.SetValue(key, value)
			h.writeResponse("OK")
		case "GET":
			if len(key) == 0 {
				h.writeError("wrong number of arguments for GET")
				continue
			}
			val, ok := commands.GetValue(key)
			if !ok {
				h.writeResponse(nil)
			} else {
				h.writeResponse(val)
			}
		case "INFO":
			info := commands.InfoValue()
			if h.resp {
				h.conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(info), info)))
			} else {
				h.conn.Write([]byte(info + "\n"))
			}
		case "DEL":
			if len(key) == 0 {
				h.writeError("wrong number of arguments for DEL")
				continue
			}
			commands.DeleteValue(key)
			h.writeResponse("OK")
		case "AUTH":
			if serverPassword == "" {
				h.writeError("AUTH not configured")
				continue
			}
			if string(value) == serverPassword {
				h.authenticated = true
				h.writeResponse("OK")
			} else {
				h.writeError("invalid password")
			}
		default:
			h.writeError("unknown command '" + cmdStr + "'")
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":64666")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("GEDIS LISTENING ON PORT 64666")

	database.InitializeDatabase()
	database.ReadDatabase()
	database.OpenAOF()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		listener.Close()
		database.CloseDatabase()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		connSem <- struct{}{}
		go func(c net.Conn) {
			defer func() { <-connSem }()
			h := &clientHandler{
				conn:   c,
				reader: bufio.NewReader(c),
			}
			h.handle()
		}(conn)
	}
}
