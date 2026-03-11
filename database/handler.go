package database

import (
	"sync"
)

var SHARDS_COUNT int = 128

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

var store []*SafeMap

func InitializeDatabase() {
	store = make([]*SafeMap, SHARDS_COUNT)
	for i := 0; i < SHARDS_COUNT; i++ {
		store[i] = &SafeMap{
			data: make(map[string]string),
		}
	}
}

func GetShard(key string) int {
	// Simple hash function to determine the shard index
	hash := 0
	for _, char := range key {
		hash += int(char)
	}
	return hash % SHARDS_COUNT
}

func CreateDatabase() {

}

func ReadDatabase() {
	// If exist the log file, import the map and then execute
	// the commands in the log file to update the map
}

func SaveDatabase() {
	// Delete the log of commands and save the map to a file
}

func CloseDatabase() {
	SaveDatabase()
}

func LogCommand(cmd string) {
	// Log command to a file
}

func GetKey(key string) (string, bool) {
	var shardIndex int = GetShard(key)
	store[shardIndex].mu.RLock()
	defer store[shardIndex].mu.RUnlock()
	value, ok := store[shardIndex].data[key]
	return value, ok
}

func SetKey(key string, value string) {
	var shardIndex int = GetShard(key)
	store[shardIndex].mu.Lock()
	defer store[shardIndex].mu.Unlock()
	store[shardIndex].data[key] = value
	go LogCommand("SET" + key + " " + value)
}

func DeleteKey(key string) {
	var shardIndex int = GetShard(key)
	store[shardIndex].mu.Lock()
	defer store[shardIndex].mu.Unlock()
	delete(store[shardIndex].data, key)
	go LogCommand("DEL" + key)
}
