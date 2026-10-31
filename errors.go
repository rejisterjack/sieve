package sieve

import (
	"errors"
	"fmt"
)

// ErrBlocked is returned by Guard when the input must be rejected.
var ErrBlocked = errors.New("sieve: input blocked")

// BlockedError wraps a scan Result for blocked inputs.
type BlockedError struct {
	Result Result
}

func (e *BlockedError) Error() string {
	return fmt.Sprintf("sieve: input blocked (risk=%s)", e.Result.Risk)
}

func (e *BlockedError) Unwrap() error {
	return ErrBlocked
}

// Guard scans input and returns ErrBlocked (via *BlockedError) when ShouldBlock is true.
func (s *Sieve) Guard(input string) (Result, error) {
	r := s.Scan(input)
	if s.ShouldBlock(r) {
		return r, &BlockedError{Result: r}
	}
	return r, nil
}
