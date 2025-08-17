package main

import (
	"fmt"
	"hash/crc32"
	"hash/fnv"
)

type BloomFilter struct {
	bitset []bool
	size   uint32
}

// ハッシュ関数1（FNV）
func hash1(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// ハッシュ関数2（CRC32）
func hash2(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// Bloom Filter作成
func NewBloomFilter(size uint32) *BloomFilter {
	return &BloomFilter{
		bitset: make([]bool, size),
		size:   size,
	}
}

// 要素を追加
func (bf *BloomFilter) Add(s string) {
	h1 := hash1(s) % bf.size
	h2 := hash2(s) % bf.size

	for i := 0; i < 3; i++ { // k=3 の場合
		pos := (h1 + uint32(i)*h2) % bf.size
		bf.bitset[pos] = true
	}
}

// 存在チェック
func (bf *BloomFilter) Exists(s string) bool {
	h1 := hash1(s) % bf.size
	h2 := hash2(s) % bf.size

	for i := 0; i < 3; i++ { // k=3 の場合
		pos := (h1 + uint32(i)*h2) % bf.size
		if !bf.bitset[pos] {
			return false // 1つでも立ってなければ「存在しない」
		}
	}
	return true
}

// 実行例
func main() {
	bf := NewBloomFilter(1000)

	bf.Add("apple")
	bf.Add("banana")

	fmt.Println("apple:", bf.Exists("apple"))   // true
	fmt.Println("banana:", bf.Exists("banana")) // true
	fmt.Println("grape:", bf.Exists("grape"))   // false or true (偽陽性の可能性)
}
