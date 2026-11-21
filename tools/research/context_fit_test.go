package main

import (
	"math"
	"reflect"
	"testing"
)

func TestContextFeatureOrientation(t *testing.T) {
	a, err := contextFeatures("4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", "w")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a.Features, []contextFeature{{151, 1}, {369, 1}}) || a.Phase != 1 {
		t.Fatal("independent knight geometry", a)
	}
	b, err := contextFeatures("4k3/8/8/3n4/8/8/8/4K3 b - - 0 1", "b")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("color mirror differs", a, b)
	}
	c, err := contextFeatures("4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", "b")
	if err != nil {
		t.Fatal(err)
	}
	for i := range a.Features {
		if a.Features[i].Index != c.Features[i].Index || a.Features[i].Count != -c.Features[i].Count {
			t.Fatal("root sign")
		}
	}
	start, err := contextFeatures("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "w")
	if err != nil || len(start.Features) != 0 || start.Phase != 24 {
		t.Fatal("starting symmetry", start, err)
	}
}
func TestContextSolverScaleOrdinalAndBounds(t *testing.T) {
	for _, c := range []struct {
		name     string
		target   float64
		ordinal  bool
		positive bool
	}{{"zero", 0, false, false}, {"residual", 100, false, true}, {"ordinal", 25, true, true}, {"satisfied ordinal", -1000, true, false}} {
		t.Run(c.name, func(t *testing.T) {
			pairs := []contextPair{{X: []contextTerm{{Index: 0, Value: 1}}, Target: c.target, Weight: 1, Ordinal: c.ordinal}}
			w, solver := solveContext(pairs)
			if solver["converged"] != true {
				t.Fatal(solver)
			}
			if c.positive && w[0] <= 0 || !c.positive && math.Abs(w[0]) > 1e-8 {
				t.Fatal("target direction", w[0])
			}
			if c.name == "residual" {
				want := 75 / (1 + contextRidge*44.0/43)
				if math.Abs(w[0]-want) > 1e-5 {
					t.Fatal("known constrained optimum", w[0], want)
				}
			}
			q := quantizeContext(w)
			for group := 0; group < 10; group++ {
				sum := 0
				for i := 0; i < contextBins; i++ {
					value := q[group*contextBins+i]
					if value < -32 || value > 32 {
						t.Fatal("coefficient cap")
					}
					sum += value
				}
				if sum != 0 || q[group*contextBins+4] != 0 {
					t.Fatal("material/intersection constraint")
				}
			}
		})
	}
	root := rootCase{ID: "control", Split: "train", Alternatives: []alternative{{Base: 1000, Teacher: label{CP: contextInt(1000)}}, {Base: 0, Teacher: label{CP: contextInt(0)}}}}
	pairs, _ := contextTrainingPairs([]rootCase{root}, map[string][]contextBranch{"control": {{Phase: 24, Features: []contextFeature{{0, 1}}}, {}}})
	if len(pairs) != 1 || pairs[0].Target != 0 {
		t.Fatal("already-correct material target", pairs)
	}
}
func contextInt(v int) *int { return &v }

func TestContextGradientAndProjection(t *testing.T) {
	pairs := []contextPair{
		{X: []contextTerm{{0, 2}, {46, -.5}}, Target: 100, Weight: .5},
		{X: []contextTerm{{90, -3}, {135, 1}}, Target: 500, Weight: .25},
		{X: []contextTerm{{180, .5}}, Target: 25, Weight: 1, Ordinal: true},
	}
	w := make([]float64, contextDimensions)
	for i := range w {
		w[i] = float64(i%7) - 3
	}
	_, gradient := contextObjective(pairs, w, true)
	for i := range w {
		old := w[i]
		w[i] = old + 1e-4
		up, _ := contextObjective(pairs, w, false)
		w[i] = old - 1e-4
		down, _ := contextObjective(pairs, w, false)
		w[i] = old
		numeric := (up - down) / 2e-4
		if math.Abs(numeric-gradient[i]) > 1e-5*math.Max(1, math.Abs(gradient[i])) {
			t.Fatal("gradient", i, numeric, gradient[i])
		}
	}
	for i := range w {
		w[i] *= 1000
	}
	projected := projectContext(w)
	for group := 0; group < 10; group++ {
		sum := 0.
		for i := 0; i < contextBins; i++ {
			value := projected[group*contextBins+i]
			sum += value
			if math.Abs(value) > 32 {
				t.Fatal("projected bound")
			}
		}
		if math.Abs(sum) > 1e-6 || projected[group*contextBins+4] != 0 {
			t.Fatal("projection constraints", sum)
		}
	}
	_ = quantizeContext(projected)
}

func TestContextLargeLossKnownOptimum(t *testing.T) {
	pairs := []contextPair{{X: []contextTerm{{0, 1}}, Target: 100, Weight: 1}}
	// Equal opposing saturated Huber losses add a large constant without
	// changing the optimum of the first constraint near its solution.
	for i := 0; i < 100; i++ {
		pairs = append(pairs, contextPair{X: []contextTerm{{0, 1}}, Target: 400, Weight: 1}, contextPair{X: []contextTerm{{0, 1}}, Target: -400, Weight: 1})
	}
	w, solver := solveContext(pairs)
	want := 75 / (1 + contextRidge*44.0/43)
	if solver["converged"] != true || math.Abs(w[0]-want) > 1e-5 {
		t.Fatal("large-loss optimum", w[0], want, solver)
	}
}
