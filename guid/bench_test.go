package main

import "testing"

func BenchmarkGenXID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenXID()
	}
}

func BenchmarkGenSonyflake(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenSonyflake()
	}
}

func BenchmarkGenSonyflakeStr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenSonyflakeStr()
	}
}

func BenchmarkGenUUID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenUUID()
	}
}
