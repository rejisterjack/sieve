package demo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Example is one demo / corpus fixture.
type Example struct {
	ID      string `json:"id"`
	Input   string `json:"input"`
	Expect  string `json:"expect"`
	Preview string `json:"preview"`
}

type corpusRow struct {
	ID     string `json:"id"`
	Input  string `json:"input"`
	Expect string `json:"expect"`
}

// LoadExamples parses embedded JSONL files and the scraped HTML page.
func LoadExamples() ([]Example, error) {
	var out []Example
	for _, name := range []string{"safe.jsonl", "blocked.jsonl", "suspicious.jsonl"} {
		data, err := corpusFS.ReadFile("corpus/" + name)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(bytes.NewReader(data))
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			var row corpusRow
			if err := json.Unmarshal([]byte(line), &row); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			out = append(out, Example{
				ID:      row.ID,
				Input:   row.Input,
				Expect:  row.Expect,
				Preview: preview(row.Input),
			})
		}
		if err := sc.Err(); err != nil {
			return nil, err
		}
	}
	html, err := corpusFS.ReadFile("corpus/scraped_page.html")
	if err != nil {
		return nil, err
	}
	out = append(out, Example{
		ID:      "scraped_page",
		Input:   string(html),
		Expect:  "blocked",
		Preview: "Scraped travel page with hidden HTML instructions",
	})
	return out, nil
}

// FindExample returns an example by id.
func FindExample(id string) (Example, bool) {
	list, err := LoadExamples()
	if err != nil {
		return Example{}, false
	}
	for _, e := range list {
		if e.ID == id {
			return e, true
		}
	}
	return Example{}, false
}

func preview(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 96 {
		return s[:93] + "..."
	}
	return s
}
