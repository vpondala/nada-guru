// Package server implements the HTTP REST server for Nāda Guru.
package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vpondala/nada-guru/agents"
	"github.com/vpondala/nada-guru/knowledge"
)

// Start begins listening on PORT (default 8080) and serves the REST API.
// It exits cleanly on SIGINT / SIGTERM.
func Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store, err := knowledge.New()
	if err != nil {
		return err
	}

	rootAgent, err := agents.New(store)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /chat", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var req struct {
			Message   string `json:"message"`
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("bad request", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request body"}`))
			return
		}
		req.Message = strings.TrimSpace(req.Message)
		if req.Message == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"message must not be empty"}`))
			return
		}
		if req.SessionID == "" {
			req.SessionID = "default"
		}

		ctx := r.Context()
		_ = ctx
		_ = rootAgent

		resp := map[string]any{
			"response":   "[agent invocation not yet wired]",
			"session_id": req.SessionID,
			"agent_used": "nada_guru",
			"latency_ms": time.Since(start).Milliseconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		slog.Info("request completed", "session_id", req.SessionID, "latency_ms", time.Since(start).Milliseconds())
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	slog.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
