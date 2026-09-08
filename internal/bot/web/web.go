// Package web serves a read-only local page listing the stored challenges and
// the completions recorded against each one.
package web

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
)

//go:embed static/htmx.min.js
var htmxJS []byte

// Handler builds the challenge viewer's HTTP handler, reading challenges from
// store.
func Handler(store port.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleChallenges(store))
	mux.HandleFunc("GET /challenges/{id}/completions", handleCompletions(store))
	mux.HandleFunc("GET /feedback", handleFeedback(store))
	mux.HandleFunc("GET /htmx.min.js", handleHTMX)
	return mux
}

// Server is the local HTTP server for browsing stored challenges.
type Server struct {
	http *http.Server
}

// New builds a Server listening on addr and reading challenges from store.
func New(addr string, store port.Store) *Server {
	return &Server{
		http: &http.Server{
			Addr:              addr,
			Handler:           Handler(store),
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

// Start serves requests until Shutdown is called. A clean shutdown returns nil.
func (s *Server) Start() error {
	slog.Info("challenge viewer listening", "addr", s.http.Addr)

	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shutdown stops the server, waiting for in-flight requests to finish.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func handleChallenges(store port.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		challenges, err := store.ListChallenges(r.Context())
		if err != nil {
			slog.Error("listing challenges for viewer", "err", err)
			http.Error(w, "could not load challenges", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := challengesPage(challengeViews(challenges)).Render(r.Context(), w); err != nil {
			slog.Error("rendering challenge viewer", "err", err)
		}
	}
}

func handleCompletions(store port.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		challengeID := r.PathValue("id")

		c, err := store.GetChallenge(r.Context(), challengeID)
		if err != nil {
			http.Error(w, "challenge not found", http.StatusNotFound)
			return
		}

		query := r.URL.Query()
		view := newCompletionsView(challengeID, c.Completions(), query.Get("sort"), query.Get("dir"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := completionsPanel(view).Render(r.Context(), w); err != nil {
			slog.Error("rendering completions panel", "challenge_id", challengeID, "err", err)
		}
	}
}

func handleFeedback(store port.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snapshot, err := store.Popularity(r.Context())
		if err != nil {
			slog.Error("loading popularity for viewer", "err", err)
			http.Error(w, "could not load feedback", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := feedbackPage(feedbackTables(snapshot)).Render(r.Context(), w); err != nil {
			slog.Error("rendering feedback viewer", "err", err)
		}
	}
}

func handleHTMX(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if _, err := w.Write(htmxJS); err != nil {
		slog.Debug("writing htmx.min.js", "err", err)
	}
}
