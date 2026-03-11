package commands

import (
	"net"
	"strings"

	"adbr.xx/gedis/database"
)

func SetValue(conn net.Conn, key string, values []string) {
	value := strings.Join(values, " ")
	database.SetKey(key, value)
	response := "OK\n"
	conn.Write([]byte(response))
}

func GetValue(conn net.Conn, key string) {
	value, ok := database.GetKey(key)
	if !ok {
		conn.Write([]byte("KEY_NOT_FOUND"))
		return
	}
	response := value + "\n"
	conn.Write([]byte(response))
}

func DeleteValue(conn net.Conn, key string) {
	database.DeleteKey(key)
	response := "OK\n"
	conn.Write([]byte(response))
}
