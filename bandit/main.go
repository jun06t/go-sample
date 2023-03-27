package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
	erand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	arm0 := newArm(0.03)
	arm1 := newArm(0.02)
	arm2 := newArm(0.01)

	ThompsonSampling(1000, arm0, arm1, arm2)

	fmt.Printf("コンテンツ0 試行回数: %d, CV回数: %d\n", arm0.total(), arm0.alpha)
	fmt.Printf("コンテンツ1 試行回数: %d, CV回数: %d\n", arm1.total(), arm1.alpha)
	fmt.Printf("コンテンツ2 試行回数: %d, CV回数: %d\n", arm2.total(), arm2.alpha)

	// check by another library
	rewards := []mab.Dist{
		mab.Beta(float64(arm0.alpha+1), float64(arm0.beta+1)),
		mab.Beta(float64(arm1.alpha+1), float64(arm1.beta+1)),
		mab.Beta(float64(arm2.alpha+1), float64(arm2.beta+1)),
	}

	b := mab.Bandit{
		RewardSource: &mab.RewardStub{Rewards: rewards},
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := b.SelectArm(context.Background(), "12345", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Arm)
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
