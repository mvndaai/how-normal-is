package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mvndaai/ctxerr"
	"github.com/mvndaai/how-normal-is/internal/domain"
)

type Server struct {
	http *http.Server
}

func NewServer() *Server {
	mux := http.NewServeMux()
	h := &handlers{}
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /api/v1/traits/suggest", h.suggestTraits)
	mux.HandleFunc("POST /api/v1/moderation/ban-window", h.banWindow)
	mux.HandleFunc("POST /api/v1/questions/archive-decision", h.archiveDecision)

	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}
	return &Server{http: &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second}}
}

func (s *Server) ListenAndServe() error {
	err := s.http.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

type handlers struct{}

func (h *handlers) health(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

var knownTraits = []domain.Trait{
	{Name: "Brown eyes", Group: "eye_color"},
	{Name: "Blue eyes", Group: "eye_color"},
	{Name: "Green eyes", Group: "eye_color"},
	{Name: "Man", Group: "gender"},
	{Name: "Woman", Group: "gender"},
	{Name: "Non-binary", Group: "gender"},
	{Name: "ADHD"},
	{Name: "Autism"},
	{Name: "Left handed"},
}

func (h *handlers) suggestTraits(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result := domain.SuggestTraits(q, knownTraits, 8)
	respondJSON(w, http.StatusOK, result)
}

type banWindowRequest struct {
	StrongRejectCount int `json:"strongRejectCount"`
}

func (h *handlers) banWindow(w http.ResponseWriter, r *http.Request) {
	var req banWindowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, ctxerr.Wrap(r.Context(), err, "decode request"), http.StatusBadRequest)
		return
	}
	respondJSON(w, http.StatusOK, domain.BanWindowForStrongRejects(time.Now().UTC(), req.StrongRejectCount))
}

type archiveDecisionRequest struct {
	YesAnswers       int `json:"yesAnswers"`
	SometimesAnswers int `json:"sometimesAnswers"`
	NeverAnswers     int `json:"neverAnswers"`
	ConfusingVotes   int `json:"confusingVotes"`
	MinVotes         int `json:"minVotes"`
}

func (h *handlers) archiveDecision(w http.ResponseWriter, r *http.Request) {
	var req archiveDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, ctxerr.Wrap(r.Context(), err, "decode request"), http.StatusBadRequest)
		return
	}
	stats := domain.QuestionStats{YesAnswers: req.YesAnswers, SometimesAnswers: req.SometimesAnswers, NeverAnswers: req.NeverAnswers, ConfusingVotes: req.ConfusingVotes}
	respondJSON(w, http.StatusOK, map[string]any{
		"shouldArchive": stats.ShouldArchiveForConfusing(req.MinVotes),
		"percentNormal": stats.PercentNormal(),
	})
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, err error, status int) {
	respondJSON(w, status, map[string]string{"error": err.Error()})
}
