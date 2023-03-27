package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/stitchfix/mab"
	"github.com/stitchfix/mab/numint"
	erand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	arms := Arms{
		newArm(0.6),
		newArm(0.4),
		newArm(0.1),
	}

	ThompsonSampling(1000, arms...)

	for i := range arms {
		fmt.Printf("コンテンツ%d 試行回数: %d, CV回数: %d\n", i, arms[i].total(), arms[i].alpha)
	}
	err := arms.plot()
	if err != nil {
		log.Fatal(err)
	}

	// check by another library
	rewards := []mab.Dist{}
	for i := range arms {
		rewards = append(rewards,
			mab.Beta(float64(arms[i].alpha+1), float64(arms[i].beta+1)),
		)
	}

	b := mab.Bandit{
		RewardSource: &mab.RewardStub{Rewards: rewards},
		Strategy:     mab.NewThompson(numint.NewQuadrature()),
		Sampler:      mab.NewSha1Sampler(),
	}

	result, err := b.SelectArm(context.Background(), "12345", nil)
	if err != nil {
		log.Fatal(err)
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

func (a Arms) plot() error {
	// Create a new plot and set its dimensions.
	p := plot.New()
	p.X.Label.Text = "Reward Probability"
	p.Y.Label.Text = "Density"
	p.Y.Min = 0

	for _, arm := range a {
		// Create the data for the plot.
		pts := make(plotter.XYs, 100)
		for i := range pts {
			x := float64(i) / float64(len(pts)-1)
			pts[i].X = x
			pts[i].Y = distuv.Beta{Alpha: float64(arm.alpha + 1), Beta: float64(arm.beta + 1)}.Prob(x)
		}

		// Create a line plotter and add it to the plot.
		lp, err := plotter.NewLine(pts)
		if err != nil {
			return err
		}
		lp.LineStyle.Width = vg.Points(1)
		lp.LineStyle.Color = color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 255,
		}
		p.Add(lp)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "arm_distribution.png"); err != nil {
		return err
	}
	return nil
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
