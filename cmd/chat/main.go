package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-chat/internal/app"
	"go-chat/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	versionFlag := flag.Bool("version", false, "print version")
	serveFlag := flag.Bool("serve", false, "run as headless relay (no TUI)")
	flag.Parse()

	if *versionFlag {
		fmt.Println("go-chat v0.1.0")
		return
	}

	if *serveFlag {
		runRelay(*configPath)
		return
	}

	a, err := app.New(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		tui.NewModel(a),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		a.Close()
		os.Exit(1)
	}

	if err := a.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during shutdown: %v\n", err)
		os.Exit(1)
	}
}

func runRelay(cfgPath string) {
	fmt.Println("go-chat relay mode")
	fmt.Println("Starting relay server...")

	a, err := app.New(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	peerID := a.PeerID()
	fmt.Printf("Relay Peer ID: %s\n", peerID)
	fmt.Println("Share this address with peers:")
	for _, addr := range a.AllAddrs() {
		fmt.Printf("  %s\n", addr)
	}
	fmt.Println()
	fmt.Println("Peers connect with: /connect <address>")
	fmt.Println("Or add to their relay_peers config.")
	fmt.Println("Press Ctrl+C to stop.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		fmt.Println("\nShutting down...")
	case <-ctx.Done():
	}

	a.Close()
}
