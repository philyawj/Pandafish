package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type warmReportRow struct {
	Game, Ply, Depth, Score                                 int
	Best, FEN                                               string
	Nodes, DeltaNodes, DeltaCandidates, DeltaContradictions uint64
	SearchMS                                                float64
}
type warmReportOptions struct {
	Values                                     bool
	Parallel                                   bool
	Baseline, Candidate, Out, Reference        string
	Pairs, Roots, Depth                        int
	MaxTimeRatio, MaxNodeRatio, MinOpportunity float64
	Exact                                      bool
}

func readWarmRows(path string) ([]warmReportRow, error) {
	var data struct{ Rows []warmReportRow }
	err := load(path, &data)
	return data.Rows, err
}
func sameWarmResult(a, b warmReportRow) bool {
	return a.Game == b.Game && a.Ply == b.Ply && a.FEN == b.FEN && a.Best == b.Best && a.Depth == b.Depth && a.Score == b.Score && a.Nodes == b.Nodes
}
func warmReport(o warmReportOptions) error {
	if o.Parallel && (o.Exact || o.Reference != "") {
		return fmt.Errorf("parallel and exact-tree profiles are distinct")
	}
	if o.Out == "" || o.Pairs < 1 || o.Roots < 1 || o.Depth < 1 {
		return fmt.Errorf("output and positive pairs/roots/depth required")
	}
	for _, v := range []float64{o.MaxTimeRatio, o.MaxNodeRatio, o.MinOpportunity} {
		if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
			return fmt.Errorf("invalid frozen gate")
		}
	}
	if _, err := os.Stat(o.Out); !os.IsNotExist(err) {
		return fmt.Errorf("refuse existing report: %s", o.Out)
	}
	base, err := filepath.Glob(o.Baseline)
	if err != nil {
		return err
	}
	candidate, err := filepath.Glob(o.Candidate)
	if err != nil {
		return err
	}
	if len(base) != o.Pairs || len(candidate) != o.Pairs {
		return fmt.Errorf("expected %d alternating pairs, found %d/%d", o.Pairs, len(base), len(candidate))
	}
	slices.Sort(base)
	slices.Sort(candidate)
	issues := []string{}
	pairs := []map[string]any{}
	ratios, nodeRatios := []float64{}, []float64{}
	identities := map[string]fileEntry{}
	var originalBase, originalCandidate, reference []warmReportRow
	if o.Reference != "" {
		reference, err = readWarmRows(o.Reference)
		if err != nil {
			return err
		}
		if len(reference) != o.Roots {
			return fmt.Errorf("reference coverage")
		}
		b, e := os.ReadFile(o.Reference)
		if e != nil {
			return e
		}
		identities[o.Reference] = hash(b)
	}
	for pair := 0; pair < o.Pairs; pair++ {
		a, err := readWarmRows(base[pair])
		if err != nil {
			return err
		}
		b, err := readWarmRows(candidate[pair])
		if err != nil {
			return err
		}
		if len(a) != o.Roots || len(b) != o.Roots {
			return fmt.Errorf("pair %d coverage %d/%d, expected %d", pair+1, len(a), len(b), o.Roots)
		}
		for _, path := range []string{base[pair], candidate[pair]} {
			data, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			identities[path] = hash(data)
		}
		if pair == 0 {
			originalBase, originalCandidate = a, b
		}
		am, bm := 0., 0.
		var an, bn, dn, dc, contradictions uint64
		changed := 0
		for i, x := range a {
			y := b[i]
			label := fmt.Sprintf("pair%d/root%d", pair+1, i)
			if x.Game != y.Game || x.Ply != y.Ply || x.FEN != y.FEN || x.Depth != o.Depth || y.Depth != o.Depth || x.Nodes == 0 || y.Nodes == 0 || x.SearchMS <= 0 || y.SearchMS <= 0 || y.DeltaNodes > y.Nodes {
				issues = append(issues, label+": coverage/depth/measurement")
			}
			candidateStable := sameWarmResult(y, originalCandidate[i])
			if o.Parallel {
				candidateStable = y.Score == originalCandidate[i].Score
				if y.Score != x.Score {
					issues = append(issues, label+": parallel value differs from release")
				}
			}
			if !sameWarmResult(x, originalBase[i]) || !candidateStable {
				issues = append(issues, label+": nondeterministic result")
			}
			if !sameWarmResult(x, y) {
				changed++
				if o.Exact {
					issues = append(issues, label+": exact baseline differs")
				}
			}
			if o.Values && (x.Score != y.Score || x.Best != y.Best) {
				issues = append(issues, label+": exact value/first choice differs")
			}
			if reference != nil && !sameWarmResult(y, reference[i]) {
				issues = append(issues, label+": exact prior candidate differs")
			}
			if pair == 0 || o.Parallel {
				g, e := probeGame(probeCase{FEN: y.FEN})
				if e != nil {
					return e
				}
				legal := false
				for _, m := range g.Position().ValidMoves() {
					if m.String() == y.Best {
						legal = true
						break
					}
				}
				if !legal {
					issues = append(issues, label+": illegal candidate move")
				}
			}
			am += x.SearchMS
			bm += y.SearchMS
			an += x.Nodes
			bn += y.Nodes
			dn += y.DeltaNodes
			dc += y.DeltaCandidates
			contradictions += y.DeltaContradictions
		}
		ratio, nr := bm/am, float64(bn)/float64(an)
		if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
			return fmt.Errorf("invalid timing")
		}
		ratios = append(ratios, ratio)
		nodeRatios = append(nodeRatios, nr)
		pairs = append(pairs, map[string]any{"pair": pair + 1, "baseline_ms": am, "candidate_ms": bm, "time_ratio": ratio, "baseline_nodes": an, "candidate_nodes": bn, "node_ratio": nr, "changed_roots": changed, "opportunity_nodes": dn, "opportunity_fraction": float64(dn) / float64(bn), "opportunity_captures": dc, "contradictions": contradictions})
	}
	slices.Sort(ratios)
	slices.Sort(nodeRatios)
	median := ratios[len(ratios)/2]
	passed := len(issues) == 0
	if o.MaxTimeRatio > 0 && median > o.MaxTimeRatio {
		passed = false
	}
	if o.MaxNodeRatio > 0 && nodeRatios[len(nodeRatios)-1] > o.MaxNodeRatio {
		passed = false
	}
	for _, p := range pairs {
		if o.MinOpportunity > 0 && p["opportunity_fraction"].(float64) < o.MinOpportunity {
			passed = false
		}
	}
	report := map[string]any{"passed": passed, "options": o, "median_time_ratio": median, "time_ratios": ratios, "node_ratios": nodeRatios, "pairs": pairs, "issues": issues, "input_identities": identities, "created_utc": time.Now().UTC().Format(time.RFC3339Nano)}
	if err = writeNewJSON(o.Out, report); err != nil {
		return err
	}
	return printJSON(map[string]any{"passed": passed, "median_time_ratio": median, "node_ratios": nodeRatios, "issue_count": len(issues), "report": o.Out})
}
