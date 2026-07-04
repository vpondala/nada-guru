// Package cli implements the interactive CLI read-eval-print loop for Nāda Guru.
package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
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

const banner = `
 ███╗   ██╗ █████╗ ██████╗  █████╗      ██████╗ ██╗   ██╗██████╗ ██╗   ██╗
 ████╗  ██║██╔══██╗██╔══██╗██╔══██╗    ██╔════╝ ██║   ██║██╔══██╗██║   ██║
 ██╔██╗ ██║███████║██║  ██║███████║    ██║  ███╗██║   ██║██████╔╝██║   ██║
 ██║╚██╗██║██╔══██║██║  ██║██╔══██║    ██║   ██║██║   ██║██╔══██╗██║   ██║
 ██║ ╚████║██║  ██║██████╔╝██║  ██║    ╚██████╔╝╚██████╔╝██║  ██╗╚██████╗
 ╚═╝  ╚═══╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝     ╚═════╝  ╚═════╝ ╚═╝╚═════╝ ╚═════╝
                   Nāda Guru · Carnatic Music Learning AI · v0.1.0
`

var port = flag.Int("port", 8080, "HTTP server port")

// Run starts the interactive CLI loop.
func Run() error {
	fmt.Print(banner)
	fmt.Println("Type your question, or 'quit' to exit.")
	fmt.Println(strings.Repeat("─", 60))

	store, err := knowledge.New()
	if err != nil {
		return fmt.Errorf("failed to load knowledge base: %w", err)
	}

	rootAgent, err := agents.New(store)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	r, err := runner.New(runner.Config{
		AppName:           "nada-guru-cli",
		Agent:             rootAgent,
		SessionService:    session.InMemoryService(),
		AutoCreateSession: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create runner: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	done := false
	go func() {
		<-sigCh
		done = true
		fmt.Println("\nBye! Subham astu 🙏")
		os.Exit(0)
	}()

	for !done {
		fmt.Print("You > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "quit" || input == "exit" {
			fmt.Println("Bye! Subham astu 🙏")
			return nil
		}

		ctx := context.Background()
		start := time.Now()
		msg := genai.NewContentFromText(input, genai.RoleUser)

		fmt.Print("Guru > ")
		var responseText strings.Builder
		var agentUsed = "nada_guru"
		var toolsCalled = []string{}
		var runErr error

		for event, err := range r.Run(ctx, "cli-user", "cli-session", msg, agent.RunConfig{}) {
			if err != nil {
				runErr = err
				fmt.Printf("\n[Error: %v]\n", err)
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
						fmt.Print(part.Text)
						responseText.WriteString(part.Text)
					}
					if part.FunctionCall != nil {
						toolsCalled = append(toolsCalled, part.FunctionCall.Name)
					}
				}
			}
		}
		fmt.Println()

		var errStr any = nil
		if runErr != nil {
			errStr = runErr.Error()
		}

		slog.Info("agent invocation completed",
			"agent", agentUsed,
			"query", input,
			"tools_called", toolsCalled,
			"latency_ms", time.Since(start).Milliseconds(),
			"session_id", "cli-session",
			"error", errStr,
		)
	}

	return nil
}
