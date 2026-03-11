package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"adbr.xx/gedis/commands"
	"adbr.xx/gedis/database"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := bufio.NewReader(conn)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		var args []string = strings.Fields(line)
		var cmd string = strings.ToUpper(args[0])
		var key string = args[1]
		var values []string = args[2:]

		switch cmd {
		case "SET":
			commands.SetValue(conn, key, values)
			continue
		case "DEL":
			commands.DeleteValue(conn, key)
			continue
		case "GET":
			commands.GetValue(conn, key)
			continue
		case "PING":
			conn.Write([]byte("PONG\n"))
			continue
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":64666")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	database.InitializeDatabase()
	fmt.Println("GEDIS LISTENING ON PORT 64666")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
