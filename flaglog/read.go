package flaglog

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/rejisterjack/sieve"
)

// Filter selects entries when reading a flag log.
type Filter struct {
	Risk     sieve.Risk
	Category string
	Limit    int
}

// ReadEntries loads all JSONL entries from path.
func ReadEntries(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Entry
	sc := bufio.NewScanner(f)
	// Allow long JSON lines.
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ApplyFilter returns entries matching filter (empty fields match all).
func ApplyFilter(entries []Entry, f Filter) []Entry {
	var out []Entry
	for _, e := range entries {
		if f.Risk != "" && e.Risk != f.Risk {
			continue
		}
		if f.Category != "" && !entryHasCategory(e, f.Category) {
			continue
		}
		out = append(out, e)
		if f.Limit > 0 && len(out) >= f.Limit {
			break
		}
	}
	return out
}

func entryHasCategory(e Entry, category string) bool {
	for _, m := range e.Matches {
		if m.Category == category {
			return true
		}
	}
	return false
}

// Summary counts entries by risk and by match category.
type Summary struct {
	Total      int
	ByRisk     map[sieve.Risk]int
	ByCategory map[string]int
}

// Summarize aggregates entry counts.
func Summarize(entries []Entry) Summary {
	s := Summary{
		Total:      len(entries),
		ByRisk:     map[sieve.Risk]int{},
		ByCategory: map[string]int{},
	}
	for _, e := range entries {
		s.ByRisk[e.Risk]++
		seen := map[string]struct{}{}
		for _, m := range e.Matches {
			if _, ok := seen[m.Category]; ok {
				continue
			}
			seen[m.Category] = struct{}{}
			s.ByCategory[m.Category]++
		}
	}
	return s
}
