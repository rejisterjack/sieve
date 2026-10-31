package sieve

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// MiddlewareOptions configures how Middleware extracts text to scan.
type MiddlewareOptions struct {
	// SourceLabel is passed to FlagLogger / ScanWithSource.
	SourceLabel string
	// Header, if set, is scanned instead of the request body (e.g. "X-User-Message").
	Header string
	// Query, if set, is scanned from the URL query parameter of that name.
	Query string
	// JSONField, if set with a JSON body, extracts that top-level string field to scan.
	// When empty and body is used, the raw body bytes are scanned as text.
	JSONField string
}

// Middleware returns net/http middleware that scans inbound input and rejects
// blocked requests with HTTP 403 and a JSON body matching the sidecar shape.
func (s *Sieve) Middleware(opts MiddlewareOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			input, body, err := extractInput(r, opts)
			if err != nil {
				writeMWJSON(w, http.StatusBadRequest, map[string]any{
					"ok":    false,
					"error": err.Error(),
				})
				return
			}
			if body != nil {
				r.Body = io.NopCloser(bytes.NewReader(body))
			}
			if input == "" {
				next.ServeHTTP(w, r)
				return
			}

			result := s.ScanWithSource(input, opts.SourceLabel)
			if s.ShouldBlock(result) {
				writeMWJSON(w, http.StatusForbidden, map[string]any{
					"ok":      false,
					"risk":    result.Risk,
					"matches": result.Matches,
					"error":   "input blocked by sieve",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractInput(r *http.Request, opts MiddlewareOptions) (string, []byte, error) {
	if opts.Header != "" {
		return r.Header.Get(opts.Header), nil, nil
	}
	if opts.Query != "" {
		return r.URL.Query().Get(opts.Query), nil, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	_ = r.Body.Close()

	if opts.JSONField != "" {
		var m map[string]any
		if err := json.Unmarshal(body, &m); err != nil {
			return "", body, err
		}
		v, _ := m[opts.JSONField].(string)
		return v, body, nil
	}
	return string(body), body, nil
}

func writeMWJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
