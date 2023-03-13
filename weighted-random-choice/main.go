package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// itemWeights アイテム毎の重み一覧
var itemWeights = []int{
	5,  // 金貨  5%
	25, // 銀貨 25%
	70, // 銅貨 70%
}

func main() {
	sort.Ints(itemWeights) // 重みの昇順でチェックしていくので重み一覧をソートする

	// 抽選結果チェックの基準となる境界値を生成
	boundaries := make([]int, len(itemWeights)+1)
	for i := 1; i < len(itemWeights)+1; i++ {
		boundaries[i] = boundaries[i-1] + itemWeights[i-1]
	}
	boundaries = boundaries[1:len(boundaries)] // 頭の0値は不要なので破棄

	// 1000000回抽選を行う
	rand.Seed(time.Now().UnixNano())
	result := make([]int, len(itemWeights))
	n := 1000000
	for i := 0; i < n; i++ {
		draw := rand.Intn(boundaries[len(boundaries)-1]) + 1
		for i, boundary := range boundaries {
			if draw <= boundary {
				result[i]++
				break
			}
		}
	}

	fmt.Printf("境界値: %v\n", boundaries)
	fmt.Println("重み  想定      結果")
	fmt.Println("-----------------------------")
	for i := 0; i < len(itemWeights); i++ {
		fmt.Printf("%4d  %f  %f\n",
			itemWeights[i],
			float64(itemWeights[i])/float64(boundaries[len(boundaries)-1]),
			float64(result[i])/float64(n),
		)
	}
}
