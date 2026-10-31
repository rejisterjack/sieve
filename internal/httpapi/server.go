package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"time"

	"github.com/rejisterjack/sieve"
	"github.com/rejisterjack/sieve/flaglog"
	"github.com/rejisterjack/sieve/internal/httpapi/demo"
)

// Server is the HTTP sidecar that screens agent inputs.
type Server struct {
	Sieve      *sieve.Sieve
	Log        *flaglog.Logger
	WebhookURL string
	Mux        *http.ServeMux
}

// New wires routes onto a ServeMux.
func New(s *sieve.Sieve, log *flaglog.Logger, webhookURL string) *Server {
	srv := &Server{Sieve: s, Log: log, WebhookURL: webhookURL, Mux: http.NewServeMux()}
	srv.Mux.HandleFunc("GET /healthz", srv.healthz)
	srv.Mux.HandleFunc("POST /v1/scan", srv.scan)
	srv.Mux.HandleFunc("GET /v1/examples", srv.listExamples)
	srv.Mux.HandleFunc("POST /v1/examples/{id}/scan", srv.scanExample)

	index, err := demo.FS.ReadFile("index.html")
	if err == nil {
		srv.Mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(index)
		})
	}
	static, err := fs.Sub(demo.FS, ".")
	if err == nil {
		srv.Mux.Handle("GET /demo/", http.StripPrefix("/demo/", http.FileServer(http.FS(static))))
	}
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
	ID      string        `json:"id,omitempty"`
	Expect  string        `json:"expect,omitempty"`
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
	s.writeScan(w, req.Input, req.Source, "", "")
}

func (s *Server) listExamples(w http.ResponseWriter, _ *http.Request) {
	examples, err := demo.LoadExamples()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, scanResponse{OK: false, Error: err.Error()})
		return
	}
	type item struct {
		ID      string `json:"id"`
		Expect  string `json:"expect"`
		Preview string `json:"preview"`
	}
	out := make([]item, 0, len(examples))
	for _, e := range examples {
		out = append(out, item{ID: e.ID, Expect: e.Expect, Preview: e.Preview})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"examples": out})
}

func (s *Server) scanExample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex, ok := demo.FindExample(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, scanResponse{OK: false, Error: "example not found"})
		return
	}
	s.writeScan(w, ex.Input, "demo:"+ex.ID, ex.ID, ex.Expect)
}

func (s *Server) writeScan(w http.ResponseWriter, input, source, id, expect string) {
	result := s.Sieve.ScanWithSource(input, source)
	if s.Log != nil {
		_ = s.Log.Write(result, source)
	}

	if s.Sieve.ShouldBlock(result) {
		if s.WebhookURL != "" {
			go notifyWebhook(s.WebhookURL, result, source)
		}
		writeJSON(w, http.StatusForbidden, scanResponse{
			OK:      false,
			Risk:    result.Risk,
			Matches: result.Matches,
			Error:   "input blocked by sieve",
			ID:      id,
			Expect:  expect,
		})
		return
	}

	writeJSON(w, http.StatusOK, scanResponse{
		OK:      true,
		Risk:    result.Risk,
		Matches: result.Matches,
		ID:      id,
		Expect:  expect,
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
