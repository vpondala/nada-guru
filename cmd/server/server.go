// Package server implements the HTTP REST server for Nāda Guru.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vpondala/nada-guru/agents"
	"github.com/vpondala/nada-guru/knowledge"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
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

	runSrv, err := runner.New(runner.Config{
		AppName:           "nada-guru-server",
		Agent:             rootAgent,
		SessionService:    session.InMemoryService(),
		AutoCreateSession: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create runner: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]any{
			"status": "ok",
			"knowledge_base": map[string]int{
				"ragas":     len(store.Ragas),
				"talas":      len(store.Talas),
				"kritis":     len(store.Kritis),
				"composers":  len(store.Composers),
			},
			"version": "0.1.0",
		}
		_ = json.NewEncoder(w).Encode(resp)
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
		msg := genai.NewContentFromText(req.Message, genai.RoleUser)

		var responseText strings.Builder
		var agentUsed = "nada_guru"
		var toolsCalled = []string{}
		var runErr error

		for event, err := range runSrv.Run(ctx, "server-user", req.SessionID, msg, agent.RunConfig{}) {
			if err != nil {
				runErr = err
				break
			}
			if event == nil {
				continue
			}
			if event.Author != "" && event.Author != "user" {
				agentUsed = event.Author
			}
			if event.LLMResponse.Content != nil {
				for _, part := range event.LLMResponse.Content.Parts {
					if part == nil {
						continue
					}
					if part.Text != "" && !part.Thought {
						responseText.WriteString(part.Text)
					}
					if part.FunctionCall != nil {
						toolsCalled = append(toolsCalled, part.FunctionCall.Name)
					}
				}
			}
		}

		var errStr any = nil
		if runErr != nil {
			errStr = runErr.Error()
			slog.Error("agent invocation failed", "session_id", req.SessionID, "error", runErr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			resp := map[string]any{
				"error":      "agent invocation failed",
				"detail":     runErr.Error(),
				"session_id": req.SessionID,
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		resp := map[string]any{
			"response":   responseText.String(),
			"session_id": req.SessionID,
			"agent_used": agentUsed,
			"latency_ms": time.Since(start).Milliseconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)

		slog.Info("agent invocation completed",
			"agent", agentUsed,
			"query", req.Message,
			"tools_called", toolsCalled,
			"latency_ms", time.Since(start).Milliseconds(),
			"session_id", req.SessionID,
			"error", errStr,
		)
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
