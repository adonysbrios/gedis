package database

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const shardsCount = 1024
const aofPath = "gedis.aof"

// SafeMap is a concurrency-safe map using a sync.RWMutex.
type SafeMap struct {
	mu   sync.RWMutex
	data map[string][]byte
}

var (
	store            []*SafeMap
	aofFile          *os.File
	aofWriter        *bufio.Writer
	aofMu            sync.Mutex
	replaying        bool
	writeCount       int64
	rewriteThreshold int64 = 10000
	rewriteCh        chan struct{}
	flushTicker      *time.Ticker
	flushDone        chan struct{}
)

// InitializeDatabase creates the shared sharded map and starts the AOF rewrite loop.
func InitializeDatabase() {
	store = make([]*SafeMap, shardsCount)
	for i := 0; i < shardsCount; i++ {
		store[i] = &SafeMap{
			data: make(map[string][]byte),
		}
	}
	if rewriteCh != nil {
		close(rewriteCh)
	}
	ch := make(chan struct{}, 1)
	rewriteCh = ch
	go func() {
		for range ch {
			RewriteAOF()
		}
	}()
}

// GetShard returns the shard index for a given key using the FNV-1a hash algorithm.
func GetShard(key []byte) int {
	h := uint32(2166136261)
	for _, c := range key {
		h ^= uint32(c)
		h *= 16777619
	}
	return int(h) % shardsCount
}

// KeyCount returns the total number of keys across all shards.
func KeyCount() int {
	count := 0
	for _, s := range store {
		s.mu.RLock()
		count += len(s.data)
		s.mu.RUnlock()
	}
	return count
}

// GetKey retrieves the value associated with the given key. Returns false if the key does not exist.
func GetKey(key []byte) ([]byte, bool) {
	keyStr := string(key)
	s := store[GetShard(key)]
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[keyStr]
	return v, ok
}

// SetKey stores the given key-value pair in the database and logs the operation to the AOF.
func SetKey(key, value []byte) {
	keyStr := string(key)
	valCopy := make([]byte, len(value))
	copy(valCopy, value)

	s := store[GetShard(key)]
	s.mu.Lock()
	s.data[keyStr] = valCopy
	s.mu.Unlock()

	if !replaying {
		logCommand("SET", key, valCopy)
	}
}

// DeleteKey removes the given key from the database and logs the operation to the AOF.
// Returns true if the key existed.
func DeleteKey(key []byte) bool {
	keyStr := string(key)
	s := store[GetShard(key)]
	s.mu.Lock()
	_, existed := s.data[keyStr]
	delete(s.data, keyStr)
	s.mu.Unlock()

	if !replaying {
		logCommand("DEL", key, nil)
	}
	return existed
}

func logCommand(cmd string, key, value []byte) {
	aofMu.Lock()
	if aofWriter == nil {
		aofMu.Unlock()
		return
	}
	var line string
	switch cmd {
	case "SET":
		line = fmt.Sprintf("3,SET,%d,%s,%d,%s", len(key), key, len(value), value)
	case "DEL":
		line = fmt.Sprintf("3,DEL,%d,%s", len(key), key)
	}
	aofWriter.WriteString(line)

	writeCount++
	needsRewrite := writeCount >= rewriteThreshold
	aofMu.Unlock()

	if needsRewrite {
		select {
		case rewriteCh <- struct{}{}:
		default:
		}
	}
}

// RewriteAOF compacts the AOF file by writing the current in-memory state to a new file
// and atomically replacing the old one.
func RewriteAOF() {
	aofMu.Lock()
	defer aofMu.Unlock()

	tmpPath := aofPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return
	}
	w := bufio.NewWriter(f)

	for _, s := range store {
		s.mu.RLock()
		for k, v := range s.data {
			line := fmt.Sprintf("3,SET,%d,%s,%d,%s", len(k), k, len(v), v)
			w.WriteString(line)
		}
		s.mu.RUnlock()
	}
	w.Flush()
	f.Close()

	if aofWriter != nil {
		aofWriter.Flush()
	}
	if aofFile != nil {
		aofFile.Close()
	}

	os.Rename(tmpPath, aofPath)

	aofFile, _ = os.OpenFile(aofPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	aofWriter = bufio.NewWriter(aofFile)
	writeCount = 0
}

// ReadDatabase replays the AOF file to restore the in-memory state on startup.
func ReadDatabase() {
	f, err := os.Open(aofPath)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		cmd, key, value, err := readAOFRecord(r)
		if err != nil {
			break
		}
		replaying = true
		switch string(cmd) {
		case "SET":
			SetKey(key, value)
		case "DEL":
			DeleteKey(key)
		}
		replaying = false
	}
}

func readAOFRecord(r *bufio.Reader) (cmd, key, value []byte, err error) {
	cmd, err = readArg(r)
	if err != nil {
		return nil, nil, nil, err
	}
	// skip trailing comma after cmd data
	r.ReadByte()

	switch string(cmd) {
	case "SET":
		key, err = readArg(r)
		if err != nil {
			return nil, nil, nil, err
		}
		// skip trailing comma after key data
		r.ReadByte()

		value, err = readArg(r)
		if err != nil {
			return nil, nil, nil, err
		}
	case "DEL":
		key, err = readArg(r)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return cmd, key, value, nil
}

func readArg(r *bufio.Reader) ([]byte, error) {
	lenBytes, err := r.ReadBytes(',')
	if err != nil {
		return nil, err
	}
	argLen := 0
	for _, c := range lenBytes[:len(lenBytes)-1] {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid character in length")
		}
		argLen = argLen*10 + int(c-'0')
	}
	data := make([]byte, argLen)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// OpenAOF opens the AOF file for appending write operations and starts a periodic flusher.
func OpenAOF() {
	if flushTicker != nil {
		flushTicker.Stop()
	}
	if flushDone != nil {
		close(flushDone)
	}

	f, err := os.OpenFile(aofPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	aofFile = f
	aofWriter = bufio.NewWriter(f)

	ft := time.NewTicker(time.Second)
	fd := make(chan struct{})
	flushTicker = ft
	flushDone = fd
	go func() {
		for {
			select {
			case <-ft.C:
				aofMu.Lock()
				if aofWriter != nil {
					aofWriter.Flush()
				}
				aofMu.Unlock()
			case <-fd:
				return
			}
		}
	}()
}

// CloseDatabase flushes and rewrites the AOF, then closes all file handles.
func CloseDatabase() {
	if flushTicker != nil {
		flushTicker.Stop()
		flushTicker = nil
	}
	if flushDone != nil {
		close(flushDone)
		flushDone = nil
	}

	RewriteAOF()

	aofMu.Lock()
	defer aofMu.Unlock()
	if aofWriter != nil {
		aofWriter.Flush()
	}
	if aofFile != nil {
		aofFile.Close()
		aofFile = nil
		aofWriter = nil
	}
}
