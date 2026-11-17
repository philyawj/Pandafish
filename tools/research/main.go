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
		return fmt.Errorf("usage: research audit|replay|decisions|mates|games|policy|wording [flags]")
	}
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	switch args[0] {
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
