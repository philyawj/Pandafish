package main

import (
	"fmt"
	"go/format"
	"math"
	"math/bits"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/corentings/chess/v2"
)

const contextBins = 45
const contextDimensions = 10 * contextBins
const contextCap = 32.0
const contextRidge = 4.0

type contextFeature struct{ Index, Count int }
type contextTerm struct {
	Index int
	Value float64
}
type contextPair struct {
	X              []contextTerm
	Target, Weight float64
	Ordinal        bool
}
type contextBranch struct {
	Features []contextFeature
	Phase    int
}

// Board interaction features, oriented for each piece's own color. Root turn
// changes only the final value sign; the leaf's side to move does not.
func contextFeatures(fen, rootTurn string) (contextBranch, error) {
	option, err := chess.FEN(fen)
	if err != nil {
		return contextBranch{}, err
	}
	sets := chess.NewGame(option).Position().Board().EvaluationPieces()
	if bits.OnesCount64(sets[0]) != 1 || bits.OnesCount64(sets[6]) != 1 || rootTurn != "w" && rootTurn != "b" {
		return contextBranch{}, fmt.Errorf("invalid context kings/turn")
	}
	kings := [2]int{bits.LeadingZeros64(sets[0]), bits.LeadingZeros64(sets[6])}
	counts := [contextDimensions]int{}
	phase := 0
	for i, set := range sets {
		phase += bits.OnesCount64(set) * [6]int{0, 4, 2, 1, 1, 0}[i%6]
		if i%6 == 0 {
			continue
		}
		side, kind := i/6, i%6-1
		sign := 1
		if side == 1 {
			sign = -1
		}
		if rootTurn == "b" {
			sign = -sign
		}
		for set != 0 {
			sq := bits.LeadingZeros64(set)
			set &^= uint64(1) << uint(63-sq)
			for relation := 0; relation < 2; relation++ {
				king := kings[side]
				if relation == 1 {
					king = kings[1-side]
				}
				from := sq
				if side == 1 {
					from ^= 56
					king ^= 56
				}
				dx := min(4, absContext(from%8-king%8))
				dy := max(-4, min(4, from/8-king/8))
				counts[(relation*5+kind)*contextBins+dx*9+dy+4] += sign
			}
		}
	}
	result := contextBranch{Phase: min(24, phase)}
	for index, count := range counts {
		if count != 0 {
			result.Features = append(result.Features, contextFeature{index, count})
		}
	}
	return result, nil
}
func absContext(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func contextDifference(a, b contextBranch) []contextTerm {
	var values [contextDimensions]float64
	for _, f := range a.Features {
		values[f.Index] += float64(f.Count*a.Phase) / 24
	}
	for _, f := range b.Features {
		values[f.Index] -= float64(f.Count*b.Phase) / 24
	}
	terms := []contextTerm{}
	for index, value := range values {
		if value != 0 {
			terms = append(terms, contextTerm{index, value})
		}
	}
	return terms
}
func contextDot(x []contextTerm, w []float64) float64 {
	s := 0.
	for _, f := range x {
		s += f.Value * w[f.Index]
	}
	return s
}

func projectContext(v []float64) []float64 {
	out := make([]float64, len(v))
	for group := 0; group < 10; group++ {
		start := group * contextBins
		lo, hi := math.Inf(1), math.Inf(-1)
		for i := 0; i < contextBins; i++ {
			if i == 4 {
				continue
			}
			lo = math.Min(lo, v[start+i]-contextCap)
			hi = math.Max(hi, v[start+i]+contextCap)
		}
		for round := 0; round < 50; round++ {
			mid := (lo + hi) / 2
			sum := 0.
			for i := 0; i < contextBins; i++ {
				if i != 4 {
					sum += math.Max(-contextCap, math.Min(contextCap, v[start+i]-mid))
				}
			}
			if sum > 0 {
				lo = mid
			} else {
				hi = mid
			}
		}
		for i := 0; i < contextBins; i++ {
			if i != 4 {
				out[start+i] = math.Max(-contextCap, math.Min(contextCap, v[start+i]-(lo+hi)/2))
			}
		}
	}
	return out
}
func contextObjective(pairs []contextPair, w []float64, gradient bool) (float64, []float64) {
	var grad []float64
	if gradient {
		grad = make([]float64, len(w))
	}
	loss := 0.
	for _, p := range pairs {
		if len(p.X) == 0 {
			continue
		}
		e := contextDot(p.X, w) - p.Target
		if p.Ordinal {
			e = math.Min(0, e)
		} else {
			e = math.Copysign(math.Max(0, math.Abs(e)-25), e)
		}
		d := math.Max(-200, math.Min(200, e))
		value := .5 * e * e
		if math.Abs(e) > 200 {
			value = 200 * (math.Abs(e) - 100)
		}
		loss += p.Weight * value
		if gradient {
			for _, f := range p.X {
				grad[f.Index] += p.Weight * d * f.Value
			}
		}
	}
	for i, value := range w {
		loss += .5 * contextRidge * value * value
		if gradient {
			grad[i] += contextRidge * value
		}
	}
	return loss, grad
}
func solveContext(pairs []contextPair) ([]float64, map[string]any) {
	w, y := make([]float64, contextDimensions), make([]float64, contextDimensions)
	lipschitz := 1.
	converged := false
	iterations := 0
	history := []map[string]any{}
	for iteration := 1; iteration <= 1500; iteration++ {
		fy, gradient := contextObjective(pairs, y, true)
		var next []float64
		var cost float64
		for {
			trial := make([]float64, len(w))
			for i := range trial {
				trial[i] = y[i] - gradient[i]/lipschitz
			}
			next = projectContext(trial)
			cost, _ = contextObjective(pairs, next, false)
			upper := fy
			for i := range next {
				d := next[i] - y[i]
				upper += gradient[i]*d + .5*lipschitz*d*d
			}
			// Summed CP losses can be millions; use a relative round-off
			// allowance rather than letting cancellation inflate curvature.
			if cost <= upper+1e-12*math.Max(1, math.Abs(fy)) {
				break
			}
			lipschitz *= 2
			if lipschitz > 1e12 {
				return w, map[string]any{"converged": false, "error": "backtracking failed", "iterations": iteration}
			}
		}
		distance, norm := 0., 0.
		for i := range w {
			distance += (next[i] - w[i]) * (next[i] - w[i])
			norm += w[i] * w[i]
		}
		iterations = iteration
		if iteration%100 == 0 || iteration == 1 {
			history = append(history, map[string]any{"iteration": iteration, "objective": cost, "lipschitz": lipschitz})
		}
		if iteration > 20 && math.Sqrt(distance) <= 1e-7*(1+math.Sqrt(norm)) {
			_, g := contextObjective(pairs, next, true)
			trial := make([]float64, len(next))
			for i := range trial {
				trial[i] = next[i] - g[i]/lipschitz
			}
			projected := projectContext(trial)
			residual := 0.
			for i := range next {
				residual += (projected[i] - next[i]) * (projected[i] - next[i])
			}
			if lipschitz/contextRidge*math.Sqrt(residual) <= 1e-7*(1+math.Sqrt(norm)) {
				w = next
				converged = true
				break
			}
		}
		// Ridge supplies a known strong-convexity lower bound. Constant
		// curvature-aware momentum avoids persistent late FISTA oscillation.
		momentum := math.Max(0, (math.Sqrt(lipschitz)-math.Sqrt(contextRidge))/(math.Sqrt(lipschitz)+math.Sqrt(contextRidge)))
		for i := range y {
			y[i] = next[i] + momentum*(next[i]-w[i])
		}
		w = next
	}
	objective, _ := contextObjective(pairs, w, false)
	return w, map[string]any{"converged": converged, "iterations": iterations, "objective": objective, "lipschitz": lipschitz, "history": history, "method": "strong-convexity-aware projected acceleration", "line_search_relative_roundoff": 1e-12, "projected_gradient_relative_tolerance": 1e-7, "stopping_metric": "curvature/ridge scaled projected step"}
}
func quantizeContext(w []float64) []int {
	q := make([]int, len(w))
	for group := 0; group < 10; group++ {
		indices := []int{}
		sum := 0
		for i := 0; i < contextBins; i++ {
			if i == 4 {
				continue
			}
			index := group*contextBins + i
			q[index] = int(math.Floor(w[index]))
			sum += q[index]
			indices = append(indices, index)
		}
		slices.SortStableFunc(indices, func(a, b int) int {
			fa, fb := w[a]-float64(q[a]), w[b]-float64(q[b])
			if fa > fb {
				return -1
			}
			if fa < fb {
				return 1
			}
			return 0
		})
		for _, index := range indices {
			if sum == 0 {
				break
			}
			if sum < 0 && q[index] < 32 {
				q[index]++
				sum++
			}
		}
		if sum != 0 {
			panic("context zero-sum quantization failed")
		}
	}
	return q
}
func contextCorrection(a contextBranch, w []int) int {
	score := 0
	for _, f := range a.Features {
		score += f.Count * int(math.Round(float64(w[f.Index]*a.Phase)/24))
	}
	return score
}

func contextTrainingPairs(roots []rootCase, features map[string][]contextBranch) ([]contextPair, map[string]int) {
	pairs := []contextPair{}
	counts := map[string]int{}
	for _, root := range roots {
		if root.Split != "train" {
			continue
		}
		local := []contextPair{}
		for i, a := range root.Alternatives {
			for j := i + 1; j < len(root.Alternatives); j++ {
				b := root.Alternatives[j]
				f := features[root.ID]
				if a.Teacher.CP != nil && b.Teacher.CP != nil {
					target := float64(*a.Teacher.CP - *b.Teacher.CP - (a.Base - b.Base))
					target = math.Max(-400, math.Min(400, target))
					local = append(local, contextPair{X: contextDifference(f[i], f[j]), Target: target})
					counts["finite"]++
				} else if losing(a) != losing(b) {
					safe, bad := i, j
					if losing(a) {
						safe, bad = j, i
					}
					local = append(local, contextPair{X: contextDifference(f[safe], f[bad]), Target: 25 - float64(root.Alternatives[safe].Base-root.Alternatives[bad].Base), Ordinal: true})
					counts["ordinal"]++
				}
			}
		}
		for _, p := range local {
			p.Weight = 1 / float64(len(local))
			if len(p.X) == 0 {
				counts["constant"]++
			}
			pairs = append(pairs, p)
		}
	}
	return pairs, counts
}

func fitKingContext(dir, out string) (resultErr error) {
	if out == "" {
		return fmt.Errorf("new fit output directory required")
	}
	if err := os.Mkdir(out, 0755); err != nil {
		return err
	}
	defer func() {
		if resultErr != nil {
			_ = writeNewJSON(filepath.Join(out, "failure.json"), map[string]string{"error": resultErr.Error()})
		}
	}()
	if _, err := decisionReport(dir); err != nil {
		return fmt.Errorf("preserved decision verification: %w", err)
	}
	var boundary struct{ Passed bool }
	if err := load(filepath.Join(dir, "boundary-audit-after-leaves.json"), &boundary); err != nil {
		return err
	}
	if !boundary.Passed {
		return fmt.Errorf("preserved family boundary audit failed")
	}
	var data struct{ Roots []rootCase }
	var leaves []struct {
		ID, Move string
		Witness  struct {
			FEN, Reason string
			BaseRoot    int
		}
	}
	if err := load(filepath.Join(dir, "fit-data.json"), &data); err != nil {
		return err
	}
	if err := load(filepath.Join(dir, "leaves.json"), &leaves); err != nil {
		return err
	}
	lookup := map[string]struct {
		FEN, Reason string
		BaseRoot    int
	}{}
	for _, row := range leaves {
		key := row.ID + "/" + row.Move
		if _, ok := lookup[key]; ok {
			return fmt.Errorf("duplicate witness %s", key)
		}
		lookup[key] = row.Witness
	}
	features := map[string][]contextBranch{}
	learnable := 0
	for _, root := range data.Roots {
		if _, ok := features[root.ID]; ok {
			return fmt.Errorf("duplicate root %s", root.ID)
		}
		rows := make([]contextBranch, len(root.Alternatives))
		for i, a := range root.Alternatives {
			if !a.Learnable {
				continue
			}
			witness, ok := lookup[root.ID+"/"+a.Move]
			if !ok || witness.Reason != "static" || witness.BaseRoot != a.Base {
				return fmt.Errorf("invalid witness/base link %s/%s", root.ID, a.Move)
			}
			f, err := contextFeatures(witness.FEN, root.Turn)
			if err != nil {
				return err
			}
			if f.Phase != a.Phase {
				return fmt.Errorf("phase differs: %s/%s", root.ID, a.Move)
			}
			rows[i] = f
			learnable++
		}
		features[root.ID] = rows
	}
	pairs, counts := contextTrainingPairs(data.Roots, features)
	if len(pairs) == 0 {
		return fmt.Errorf("no training constraints")
	}
	w, solver := solveContext(pairs)
	q := quantizeContext(w)
	identities, err := workflowIdentities([]string{filepath.Join(dir, "fit-data.json"), filepath.Join(dir, "leaves.json"), filepath.Join(dir, "model-weights.json"), filepath.Join(dir, "fit-results.json"), filepath.Join(dir, "boundary-audit-after-leaves.json")})
	if err != nil {
		return err
	}
	model := map[string]any{"weights": q, "continuous_weights": w, "dimensions": contextDimensions, "cap_cp": contextCap, "ridge": contextRidge, "huber_delta": 200, "finite_band_cp": 25, "ordinal_margin_cp": 25, "input_identities": identities, "solver": solver, "constraints": counts, "learnable_witnesses": learnable, "data_role": "Old validation/holdout are inspected development controls; fresh actual-search confirmation remains necessary."}
	if err = writeNewJSON(filepath.Join(out, "model.json"), model); err != nil {
		return err
	}
	if solver["converged"] != true {
		report := map[string]any{"passed": false, "solver": solver, "development_evaluated": false}
		if err = writeNewJSON(filepath.Join(out, "summary.json"), report); err != nil {
			return err
		}
		return printJSON(report)
	}
	all := map[string]map[string]metrics{}
	brief := map[string]any{}
	passed := solver["converged"] == true
	for _, split := range []string{"train", "validation", "holdout"} {
		base, adjusted := []rootCase{}, []rootCase{}
		for _, root := range data.Roots {
			if root.Split != split {
				continue
			}
			base = append(base, root)
			changed := root
			changed.Alternatives = append([]alternative(nil), root.Alternatives...)
			for i := range changed.Alternatives {
				a := &changed.Alternatives[i]
				if a.Learnable {
					a.Base += contextCorrection(features[root.ID][i], q)
				}
				a.Learnable = false
			}
			adjusted = append(adjusted, changed)
		}
		b, err := scoreRoots(base, make([]int, 24))
		if err != nil {
			return err
		}
		c, err := scoreRoots(adjusted, make([]int, 24))
		if err != nil {
			return err
		}
		all[split] = map[string]metrics{"baseline": b, "candidate": c}
		oldRisk := map[string]bool{}
		for _, r := range b.Decisions {
			oldRisk[r.ID] = r.Risk
		}
		newRisk := []string{}
		for _, r := range c.Decisions {
			if r.Risk && !oldRisk[r.ID] {
				newRisk = append(newRisk, r.ID)
			}
		}
		good := c.Mean <= .95*b.Mean && c.Severe <= b.Severe && len(newRisk) == 0
		if split != "train" && !good {
			passed = false
		}
		brief[split] = map[string]any{"roots": b.Roots, "baseline_mean": b.Mean, "candidate_mean": c.Mean, "baseline_severe": b.Severe, "candidate_severe": c.Severe, "new_avoidable_loss_roots": newRisk, "surrogate_gate_passed": good}
	}
	if err = writeNewJSON(filepath.Join(out, "metrics.json"), all); err != nil {
		return err
	}
	var source strings.Builder
	source.WriteString("// Code generated by research context-fit; DO NOT EDIT.\npackage main\n\nvar kingContextWeights = [10][45]int16{\n")
	for group := 0; group < 10; group++ {
		source.WriteString("{")
		for _, value := range q[group*contextBins : (group+1)*contextBins] {
			fmt.Fprintf(&source, "%d,", value)
		}
		source.WriteString("},\n")
	}
	source.WriteString("}\n")
	formatted, err := format.Source([]byte(source.String()))
	if err != nil {
		return err
	}
	if err = writeNewFile(filepath.Join(out, "king_context_weights.go.txt"), formatted); err != nil {
		return err
	}
	report := map[string]any{"passed": passed, "solver": solver, "constraints": counts, "splits": brief, "data_role": "inspected development surrogate, not strength"}
	if err = writeNewJSON(filepath.Join(out, "summary.json"), report); err != nil {
		return err
	}
	return printJSON(map[string]any{"passed": passed, "converged": solver["converged"], "iterations": solver["iterations"], "constraints": counts, "splits": brief})
}
