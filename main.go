// Nāda Guru — multi-agent Carnatic music learning system.
//
// Usage:
//
//	go run main.go --mode=cli      # interactive CLI
//	go run main.go --mode=server   # HTTP REST server on PORT (default 8080)
//	go run ./cmd/launcher          # web UI for agent debugging (ADK launcher)
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/vpondala/nada-guru/cmd/cli"
	"github.com/vpondala/nada-guru/cmd/server"
)

const banner = `
 ███╗   ██╗ █████╗ ██████╗  █████╗      ██████╗ ██╗   ██╗██████╗ ██╗   ██╗
 ████╗  ██║██╔══██╗██╔══██╗██╔══██╗    ██╔════╝ ██║   ██║██╔══██╗██║   ██║
 ██╔██╗ ██║███████║██║  ██║███████║    ██║  ███╗██║   ██║██████╔╝██║   ██║
 ██║╚██╗██║██╔══██║██║  ██║██╔══██║    ██║   ██║██║   ██║██╔══██╗██║   ██║
 ██║ ╚████║██║  ██║██████╔╝██║  ██║    ╚██████╔╝╚██████╔╝██║  ██║╚██████╔╝
 ╚═╝  ╚═══╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝    ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝
                  Nāda Guru · Carnatic Music Learning AI · v0.1.0
`

func main() {
	mode := flag.String("mode", "cli", "run mode: cli | server")
	flag.Parse()

	// Configure structured logging based on LOG_FORMAT env var.
	logFormat := os.Getenv("LOG_FORMAT")
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "ts", Value: a.Value}
			}
			return a
		},
	}
	var handler slog.Handler
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))

	fmt.Print(banner)

	switch *mode {
	case "cli":
		if err := cli.Run(); err != nil {
			slog.Error("CLI exited with error", "error", err)
			os.Exit(1)
		}
	case "server":
		if err := server.Start(); err != nil {
			slog.Error("Server exited with error", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q — use --mode=cli or --mode=server\n", *mode)
		os.Exit(2)
	}
}
