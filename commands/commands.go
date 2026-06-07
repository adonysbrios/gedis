package commands

import (
	"fmt"
	"runtime"
	"time"

	"adbr.xx/gedis/database"
)

var startedAt = time.Now()

// SetValue stores a key-value pair in the database.
func SetValue(key, value []byte) {
	database.SetKey(key, value)
}

// GetValue retrieves the value for the given key. Returns false if not found.
func GetValue(key []byte) ([]byte, bool) {
	return database.GetKey(key)
}

// DeleteValue removes the given key from the database. Returns true if the key existed.
func DeleteValue(key []byte) bool {
	return database.DeleteKey(key)
}

// InfoValue returns a formatted string with server info (version, uptime, keys, memory).
func InfoValue() string {
	uptime := int(time.Since(startedAt).Seconds())
	keys := database.KeyCount()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return fmt.Sprintf(
		"# Server\ngedis_version: 0.1.0\nuptime_seconds: %d\n\n# Data\nkeys: %d\nshards: %d\naof_enabled: 1\n\n# Memory\nalloc_mb: %.2f\ntotal_alloc_mb: %.2f\n",
		uptime, keys, 1024, float64(m.Alloc)/1024/1024, float64(m.TotalAlloc)/1024/1024,
	)
}
