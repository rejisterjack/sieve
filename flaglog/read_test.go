package flaglog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rejisterjack/sieve"
)

func TestReadAndFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sieve.log")
	l := New(path)
	_ = l.Write(sieve.Result{
		Risk: sieve.RiskBlocked,
		Matches: []sieve.Match{
			{Category: "jailbreak", Pattern: "x", Risk: sieve.RiskBlocked},
		},
	}, "a")
	_ = l.Write(sieve.Result{
		Risk: sieve.RiskSuspicious,
		Matches: []sieve.Match{
			{Category: "delimiter_injection", Pattern: "y", Risk: sieve.RiskSuspicious},
		},
	}, "b")

	entries, err := ReadEntries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %d", len(entries))
	}

	blocked := ApplyFilter(entries, Filter{Risk: sieve.RiskBlocked})
	if len(blocked) != 1 || blocked[0].Risk != sieve.RiskBlocked {
		t.Fatalf("blocked filter = %+v", blocked)
	}

	sum := Summarize(entries)
	if sum.Total != 2 || sum.ByRisk[sieve.RiskBlocked] != 1 {
		t.Fatalf("summary = %+v", sum)
	}
}

func TestReadMissingFile(t *testing.T) {
	_, err := ReadEntries(filepath.Join(t.TempDir(), "missing.log"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("err = %v", err)
	}
}
