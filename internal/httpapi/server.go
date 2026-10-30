package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rejisterjack/sieve"
	"github.com/rejisterjack/sieve/internal/flaglog"
)

// Server is the HTTP sidecar that screens agent inputs.
type Server struct {
	Sieve     *sieve.Sieve
	Log       *flaglog.Logger
	WebhookURL string
	Mux       *http.ServeMux
}

// New wires routes onto a ServeMux.
func New(s *sieve.Sieve, log *flaglog.Logger, webhookURL string) *Server {
	srv := &Server{Sieve: s, Log: log, WebhookURL: webhookURL, Mux: http.NewServeMux()}
	srv.Mux.HandleFunc("GET /healthz", srv.healthz)
	srv.Mux.HandleFunc("POST /v1/scan", srv.scan)
	return srv
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
}

type scanRequest struct {
	Input  string `json:"input"`
	Source string `json:"source,omitempty"`
}

type scanResponse struct {
	OK      bool          `json:"ok"`
	Risk    sieve.Risk    `json:"risk"`
	Matches []sieve.Match `json:"matches,omitempty"`
	Error   string        `json:"error,omitempty"`
}

func (s *Server) scan(w http.ResponseWriter, r *http.Request) {
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, scanResponse{OK: false, Error: "invalid json"})
		return
	}
	if req.Input == "" {
		writeJSON(w, http.StatusBadRequest, scanResponse{OK: false, Error: "input is required"})
		return
	}

	result := s.Sieve.Scan(req.Input)
	_ = s.Log.Write(result, req.Source)

	if s.Sieve.ShouldBlock(result) {
		if s.WebhookURL != "" {
			go notifyWebhook(s.WebhookURL, result, req.Source)
		}
		writeJSON(w, http.StatusForbidden, scanResponse{
			OK:      false,
			Risk:    result.Risk,
			Matches: result.Matches,
			Error:   "input blocked by sieve",
		})
		return
	}

	writeJSON(w, http.StatusOK, scanResponse{
		OK:      true,
		Risk:    result.Risk,
		Matches: result.Matches,
	})
}

func writeJSON(w http.ResponseWriter, status int, body scanResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func notifyWebhook(url string, result sieve.Result, source string) {
	payload, err := json.Marshal(map[string]any{
		"risk":    result.Risk,
		"matches": result.Matches,
		"source":  source,
		"at":      time.Now().UTC(),
	})
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
}
