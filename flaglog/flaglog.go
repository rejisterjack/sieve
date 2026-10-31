package flaglog

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/rejisterjack/sieve"
)

// Entry is one flagged or blocked scan written to the structured log.
type Entry struct {
	Timestamp time.Time     `json:"timestamp"`
	Risk      sieve.Risk    `json:"risk"`
	Matches   []sieve.Match `json:"matches"`
	Source    string        `json:"source,omitempty"`
}

// Logger appends JSON lines for flagged inputs.
type Logger struct {
	mu   sync.Mutex
	path string
}

// New creates a logger that writes to path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Write appends an entry when the result is not safe.
func (l *Logger) Write(result sieve.Result, source string) error {
	if result.Risk == sieve.RiskSafe {
		return nil
	}
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Risk:      result.Risk,
		Matches:   result.Matches,
		Source:    source,
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}
