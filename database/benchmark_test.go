package database

import (
	"fmt"
	"testing"
)

func BenchmarkSetKey(b *testing.B) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		value := []byte(fmt.Sprintf("value%d", i))
		SetKey(key, value)
	}
}

func BenchmarkGetKey(b *testing.B) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		value := []byte(fmt.Sprintf("value%d", i))
		SetKey(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i%10000))
		GetKey(key)
	}
}

func BenchmarkDeleteKey(b *testing.B) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		SetKey(key, []byte("value"))
		DeleteKey(key)
	}
}

func BenchmarkGetShard(b *testing.B) {
	keys := make([][]byte, b.N)
	for i := range keys {
		keys[i] = []byte(fmt.Sprintf("very-long-key-name-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetShard(keys[i%len(keys)])
	}
}

func BenchmarkConcurrentSetGet(b *testing.B) {
	cleanAOF()
	defer cleanAOF()
	InitializeDatabase()
	OpenAOF()
	defer CloseDatabase()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := []byte(fmt.Sprintf("key%d", i))
			value := []byte(fmt.Sprintf("value%d", i))
			SetKey(key, value)
			GetKey(key)
			i++
		}
	})
}
