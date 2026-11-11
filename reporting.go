package main

import (
	"encoding/json"
	"fmt"
	"github.com/philyawj/Pandafish/internal/telemetry"
	"io"
)

func reportSearch(out io.Writer, result searchResult) {
	m := result.metrics
	if m.Depth > 0 {
		fmt.Fprintf(out, "info depth %d nodes %d time %d", m.Depth, m.Nodes, int64(m.ElapsedMS))
		if m.ScoreKind == "cp" {
			fmt.Fprintf(out, " score cp %d", m.Score)
		}
		fmt.Fprintln(out, " pv", result.move)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintln(out, telemetry.Prefix+string(data))
}
