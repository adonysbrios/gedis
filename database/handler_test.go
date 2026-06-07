package database

import (
	"os"
	"testing"
)

func cleanAOF() {
	os.Remove("gedis.aof")
	os.Remove("gedis.aof.tmp")
}

func TestSetAndGet(t *testing.T) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	key := []byte("testkey")
	value := []byte("testvalue")

	SetKey(key, value)

	got, ok := GetKey(key)
	if !ok {
		t.Fatal("expected key to exist")
	}
	if string(got) != string(value) {
		t.Errorf("got %q, want %q", string(got), string(value))
	}
}

func TestGetNonExistent(t *testing.T) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	_, ok := GetKey([]byte("nonexistent"))
	if ok {
		t.Error("expected ok=false for non-existent key")
	}
}

func TestDelete(t *testing.T) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	key := []byte("todelete")
	SetKey(key, []byte("value"))

	existed := DeleteKey(key)
	if !existed {
		t.Error("expected DeleteKey to return true")
	}

	_, ok := GetKey(key)
	if ok {
		t.Error("expected key to be deleted")
	}

	existed = DeleteKey(key)
	if existed {
		t.Error("expected DeleteKey to return false for already deleted key")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func(n int) {
			key := []byte{byte(n)}
			SetKey(key, []byte{byte(n)})
			got, ok := GetKey(key)
			if !ok {
				t.Errorf("key %d not found", n)
			} else if len(got) != 1 || got[0] != byte(n) {
				t.Errorf("key %d: unexpected value", n)
			}
			DeleteKey(key)
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestShardDistribution(t *testing.T) {
	cleanAOF()
	defer cleanAOF()
	keys := make([][]byte, 1000)
	for i := range keys {
		keys[i] = []byte{byte(i)}
	}

	shardCounts := make(map[int]int)
	for _, k := range keys {
		s := GetShard(k)
		shardCounts[s]++
	}

	if len(shardCounts) <= 1 {
		t.Errorf("expected keys to be distributed across multiple shards, got %d shards", len(shardCounts))
	}
}

func TestAOFReplay(t *testing.T) {
	cleanAOF()
	defer cleanAOF()

	InitializeDatabase()
	OpenAOF()
	SetKey([]byte("persistkey"), []byte("persistvalue"))
	CloseDatabase()

	InitializeDatabase()
	ReadDatabase()
	OpenAOF()

	val, ok := GetKey([]byte("persistkey"))
	if !ok {
		t.Fatal("expected key to exist after AOF replay")
	}
	if string(val) != "persistvalue" {
		t.Errorf("got %q, want %q", string(val), "persistvalue")
	}
	CloseDatabase()
}

func TestRewriteAOF(t *testing.T) {
	cleanAOF()
	defer cleanAOF()

	InitializeDatabase()
	OpenAOF()
	for i := 0; i < 100; i++ {
		key := []byte{byte(i)}
		value := []byte{byte(i + 1)}
		SetKey(key, value)
	}
	DeleteKey([]byte{0})
	CloseDatabase()

	InitializeDatabase()
	ReadDatabase()
	OpenAOF()

	for i := 1; i < 100; i++ {
		key := []byte{byte(i)}
		val, ok := GetKey(key)
		if !ok {
			t.Errorf("key %d missing after rewrite replay", i)
			continue
		}
		if int(val[0]) != i+1 {
			t.Errorf("key %d: got %d, want %d", i, val[0], i+1)
		}
	}
	_, ok := GetKey([]byte{0})
	if ok {
		t.Error("deleted key should not exist after rewrite replay")
	}
	CloseDatabase()
}
