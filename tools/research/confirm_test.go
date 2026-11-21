package main

import "testing"

func TestConfirmationKeepsBothActorsAndMateLabelsSeparate(t *testing.T) {
	negative, cp := -7, -800
	row := probeNative{Trials: []probeTrial{{Actor: "baseline", Query: probeQuery{Move: "a1a2"}}, {Actor: "candidate", Query: probeQuery{Move: "b1b2"}}}}
	scores := map[string]probeScore{"a1a2": {CP: &cp}, "b1b2": {Mate: &negative}}
	flags := confirmedMateFlags(row, scores)
	if flags["baseline"] || !flags["candidate"] {
		t.Fatal(flags)
	}
	scores["a1a2"] = probeScore{Mate: &negative}
	flags = confirmedMateFlags(row, scores)
	if !flags["baseline"] || !flags["candidate"] {
		t.Fatal("both losing reports hidden", flags)
	}
}
