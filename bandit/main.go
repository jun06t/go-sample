package main

import (
	"fmt"
	"math/rand"
	"time"

	erand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	arm1 := newArm(0.03)
	arm2 := newArm(0.02)
	arm3 := newArm(0.01)

	ThompsonSampling(1000, arm1, arm2, arm3)

	fmt.Printf("コンテンツA 試行回数: %d, CV回数: %d\n", arm1.total(), arm1.alpha)
	fmt.Printf("コンテンツB 試行回数: %d, CV回数: %d\n", arm2.total(), arm2.alpha)
	fmt.Printf("コンテンツC 試行回数: %d, CV回数: %d\n", arm3.total(), arm3.alpha)
}

func ThompsonSampling(pull int, arms ...*Arm) {
	for i := 0; i < pull; i++ {
		arm := Arms(arms).next()
		arm.reward()
	}
}

type Arms []*Arm

func (a Arms) next() *Arm {
	idx := 0
	var max float64
	for i := range a {
		if prob := a[i].sampling(); prob > max {
			idx = i
			max = prob
		}
	}
	return a[idx]
}

type Arm struct {
	probs float64
	alpha int
	beta  int
}

func newArm(prob float64) *Arm {
	return &Arm{
		probs: prob,
	}
}

func (a *Arm) reward() {
	if rand.Float64() < a.probs {
		a.alpha++
	} else {
		a.beta++
	}
}

func (a *Arm) total() int {
	return a.alpha + a.beta
}

func (a *Arm) sampling() float64 {
	bd := distuv.Beta{
		Alpha: float64(a.alpha) + 1,
		Beta:  float64(a.beta) + 1,
		Src:   erand.NewSource(uint64(time.Now().Nanosecond())),
	}
	return bd.Rand()
}
