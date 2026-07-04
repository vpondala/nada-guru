// Package cli implements the interactive CLI read-eval-print loop for Nāda Guru.
package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vpondala/nada-guru/agents"
	"github.com/vpondala/nada-guru/knowledge"
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
	fmt.Println(banner)
	fmt.Println("Type your question, or 'quit' to exit.")
	fmt.Println(strings.Repeat("─", 60))

	store, err := knowledge.New()
	if err != nil {
		return fmt.Errorf("failed to load knowledge base: %w", err)
	}

	agent, err := agents.New(store)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
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
		// In a real implementation, this would invoke the agent.
		// For now, just acknowledge input.
		_ = agent
		_ = ctx
		fmt.Printf("Guru > [agent not yet wired in CLI runner]\n")
	}

	return nil
}
