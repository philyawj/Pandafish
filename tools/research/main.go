// Command research verifies preserved experiments and replays archived Go engines.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: research workflow|validation-plan|native-plan|warm-plan|warm-report|depth-review|confirm-mates|snapshot|inventory|map-docs|audit|replay|decisions|mates|games|policy|wording|choices|generate|context-fit|context-targets [flags]")
	}
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	switch args[0] {
	case "validation-plan":
		candidate := f.String("candidate", "tools/benchmark/bin/pandafish-candidate", "development executable destination")
		testBinary := f.String("candidate-test", "/tmp/pandafish-candidate.test", "warm test executable destination")
		out := f.String("out", "", "new validation plan")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return createValidationPlan(*candidate, *testBinary, *out)
	case "confirm-mates":
		dir := f.String("dir", "", "completed choices directory with saved new mate flags")
		teacher := f.String("teacher", "/opt/homebrew/bin/stockfish", "same pinned teacher")
		out := f.String("out", "", "new confirmation directory")
		nodes := f.Int("nodes", 1000000, "fixed larger reference node budget")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return confirmChoiceMates(*dir, *teacher, *out, *nodes)
	case "depth-review":
		flagsOnly := f.Bool("flags-only", false, "review only new mate flags from each native file's sibling quality.json")
		base := f.String("baseline", "", "tagged baseline executable")
		candidate := f.String("candidate", "", "frozen candidate executable")
		input := f.String("input", "", "completed native.json glob")
		out := f.String("out", "", "new diagnostic directory")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return reviewDepths(*base, *candidate, *input, *out, *flagsOnly)
	case "context-targets":
		dir := f.String("dir", "tools/benchmark/research/v0.19-sharp-branches", "preserved balanced corpus")
		teacher := f.String("teacher", "/opt/homebrew/bin/stockfish", "pinned reference executable")
		out := f.String("out", "", "staged new label directory")
		stage := f.String("stage", "diagnose", "diagnose or complete after passed alignment")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return collectContextTargets(*dir, *teacher, *out, *stage)
	case "map-docs":
		record := f.String("record", "", "new before/after transition JSON")
		var files pathFlags
		f.Var(&files, "file", "explicit authorized historical Markdown document; repeatable")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return mapDocuments(*record, files)
	case "context-fit":
		dir := f.String("dir", "tools/benchmark/research/v0.19-sharp-branches", "preserved balanced decision corpus")
		out := f.String("out", "", "new contextual fit directory")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return fitKingContext(*dir, *out)
	case "warm-plan":
		var o warmPlanOptions
		f.StringVar(&o.Preflight, "preflight-state", "", "passed validation workflow state with matching core sources")
		f.BoolVar(&o.Values, "values", false, "require baseline values and first choices; deterministic nodes may differ")
		f.BoolVar(&o.Parallel, "parallel", false, "parallel-value profile with separate serial co-best verification")
		f.StringVar(&o.BaselineTest, "baseline-test", "", "compiled tagged baseline test executable")
		f.StringVar(&o.CandidateTest, "candidate-test", "", "compiled candidate test executable")
		f.StringVar(&o.Input, "input", "", "frozen trajectory directory")
		f.StringVar(&o.Evidence, "evidence", "", "existing warm evidence directory")
		f.StringVar(&o.Out, "out", "", "new JSON plan")
		f.StringVar(&o.Reference, "reference", "", "optional initial candidate row file for exact compaction")
		f.IntVar(&o.Depth, "depth", 5, "frozen complete depth")
		f.IntVar(&o.Roots, "roots", 1145, "expected roots")
		f.Float64Var(&o.MaxTimeRatio, "max-time-ratio", 0, "frozen median elapsed ratio limit; zero skips")
		f.Float64Var(&o.MaxNodeRatio, "max-node-ratio", 0, "frozen actual node ratio limit; zero skips")
		f.BoolVar(&o.Exact, "exact", false, "require identical baseline moves/scores/nodes")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if o.Preflight == "" {
			return fmt.Errorf("warm-plan requires --preflight-state from passed source-matching validation")
		}
		return createWarmPlan(o)
	case "workflow":
		plan := f.String("plan", "", "frozen JSON step plan")
		out := f.String("out", "", "run directory")
		resume := f.Bool("resume", false, "verify saved identities and continue only after passed steps")
		status := f.Bool("status", false, "read compact state without running commands")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if *status {
			return workflowStatus(*out)
		}
		err := runWorkflow(*plan, *out, *resume)
		statusErr := workflowStatus(*out)
		if err != nil {
			return err
		}
		return statusErr
	case "native-plan":
		base := f.String("baseline", "", "released engine executable")
		candidate := f.String("candidate", "", "frozen candidate executable")
		teacher := f.String("teacher", "/opt/homebrew/bin/stockfish", "reference executable")
		cache := f.String("cache", "", "completed matching reference suite directory")
		evidence := f.String("evidence", "", "existing candidate evidence directory")
		profile := f.String("profile", "existing", "existing or fresh")
		input := f.String("input", "", "fixed fresh control input")
		out := f.String("out", "", "new plan JSON file")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return nativePlan(*base, *candidate, *teacher, *cache, *evidence, *profile, *input, *out)
	case "snapshot":
		base := f.String("base", "", "released base tag")
		out := f.String("out", "", "new candidate-source.json file")
		var include pathFlags
		f.Var(&include, "include", "additional repository Go/JSON file glob, e.g. tools/research/*.go; repeatable")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return createSnapshot(*base, *out, include...)
	case "inventory":
		dir := f.String("dir", "", "completed evidence directory")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return createInventory(*dir)
	case "warm-report":
		var o warmReportOptions
		f.BoolVar(&o.Values, "values", false, "require baseline values and first choices; deterministic nodes may differ")
		f.BoolVar(&o.Parallel, "parallel", false, "require release values/legality while reporting scheduling node variation")
		f.StringVar(&o.Baseline, "baseline", "", "baseline row-file glob")
		f.StringVar(&o.Candidate, "candidate", "", "candidate row-file glob")
		f.StringVar(&o.Reference, "reference", "", "optional prior candidate for exact compaction comparison")
		f.StringVar(&o.Out, "out", "", "new report JSON file")
		f.IntVar(&o.Pairs, "pairs", 3, "expected pair count")
		f.IntVar(&o.Roots, "roots", 1145, "expected roots in every file")
		f.IntVar(&o.Depth, "depth", 5, "required complete depth")
		f.Float64Var(&o.MaxTimeRatio, "max-time-ratio", 0, "frozen maximum median elapsed ratio; zero skips timing gate")
		f.Float64Var(&o.MaxNodeRatio, "max-node-ratio", 0, "frozen maximum actual node ratio in every pair; zero skips")
		f.Float64Var(&o.MinOpportunity, "min-opportunity", 0, "frozen minimum outermost diagnostic-node fraction; zero skips")
		f.BoolVar(&o.Exact, "exact", false, "also require baseline moves/scores/nodes identical")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return warmReport(o)
	case "generate":
		base := f.String("baseline", "", "released student executable")
		teacher := f.String("teacher", "/opt/homebrew/bin/stockfish", "reference executable")
		out := f.String("out", "", "new generation directory")
		seed := f.Uint64("seed", 2026111728, "frozen collection seed")
		stream := f.Uint64("stream", 20, "frozen collection stream")
		var excludeInputs, excludeAnalysis pathFlags
		f.Var(&excludeInputs, "exclude-input", "additional consumed case-file glob; repeatable")
		f.Var(&excludeAnalysis, "exclude-analysis", "saved game-analysis JSON glob; repeatable")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return generateProbeCasesWithOptions(*base, *teacher, *out, *seed, *stream, excludeInputs, excludeAnalysis)
	case "choices":
		base := f.String("baseline", "", "released engine executable")
		candidate := f.String("candidate", "", "candidate engine executable")
		teacher := f.String("teacher", "/opt/homebrew/bin/stockfish", "reference engine")
		input := f.String("input", "tools/benchmark/research/v0.17-search-moves/inputs.json", "position cases")
		out := f.String("out", "", "new output directory")
		budget := f.Int("milliseconds", 100, "equal-time allowance")
		trials := f.Int("trials", 3, "trials per actor and position")
		workers := f.Int("workers", 3, "concurrent native searches")
		cache := f.String("cache", "", "optional previously saved matching 100000-node reference queries")
		byTurn := f.Bool("by-turn", false, "also require each root turn's finite-loss tolerance")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		return probeChoices(*base, *candidate, *teacher, *input, *out, *budget, *trials, *workers, *cache, *byTurn)
	case "games":
		dir := f.String("dir", "", "directory containing benchmark summary.json and PGNs")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if *dir == "" {
			return fmt.Errorf("--dir is required")
		}
		v, err := checkGames(*dir)
		if err != nil {
			return err
		}
		return printJSON(v)
	case "policy":
		base := f.String("base", "v0.18", "released base tag")
		version := f.String("version", "v0.19", "current UCI version text")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		v, err := checkPolicy(*base, *version)
		if err != nil {
			return err
		}
		return printJSON(v)
	case "wording":
		file := f.String("file", "", "PR description or release notes to check")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		b, err := os.ReadFile(*file)
		if err != nil {
			return err
		}
		if regexp.MustCompile(`(?i)\bmaintainers?\b`).Match(b) {
			return fmt.Errorf("publication text contains a prohibited role word")
		}
		if regexp.MustCompile(`(?i)\b(i|we|my|our|me|ours)\b`).Match(b) {
			return fmt.Errorf("publication text must use neutral, factual prose without first-person voice")
		}
		return printJSON(map[string]bool{"wording_passed": true})
	case "audit":
		root := f.String("root", ".", "repository root")
		legacy := f.String("legacy", "", "optional original release evidence tar.gz to verify")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		v, err := audit(*root, *legacy)
		if err != nil {
			return err
		}
		return printJSON(v)
	case "replay":
		source := f.String("source", "", "archived candidate-source.json; omitted for unchanged base")
		base := f.String("base", "v0.18", "base tag when source is omitted")
		out := f.String("out", "", "new output directory, required")
		test := f.String("test", "", "optional archived Go test file")
		check := f.Bool("check", false, "run Go tests and build in replay directory, retaining logs")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		if *out == "" {
			return fmt.Errorf("--out is required")
		}
		return replay(*source, *base, *out, *test, *check)
	case "decisions":
		dir := f.String("dir", "tools/benchmark/research/v0.19-sharp-branches", "sharp-branch archive")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		v, err := decisionReport(*dir)
		if err != nil {
			return err
		}
		return printJSON(v)
	case "mates":
		input := f.String("input", "tools/benchmark/research/v0.19-sharp-branches/training-quiet-threat-scan.json", "saved training proof records")
		if err := f.Parse(args[1:]); err != nil {
			return err
		}
		v, err := mateProofs(*input)
		if err != nil {
			return err
		}
		return printJSON(v)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printJSON(v any) error {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	return e.Encode(v)
}
func load(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
