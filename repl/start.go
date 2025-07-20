package repl

import (
	"context"
	"errors"
	"fmt"
	"github.com/sajuno/goon/agent"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

func Start(ctx context.Context, ag *agent.Agent) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m> \033[0m",
		HistoryFile:     "/tmp/goon_history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("failed to start REPL: %v", err)
	}
	defer rl.Close()

	fmt.Println("Goon REPL is ready. Type ':help' or enter a command.")

	h := newCommandHandler(ag)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, exiting..")
			return
		default:
			line, err := rl.Readline()
			if err != nil {
				if errors.Is(err, readline.ErrInterrupt) {
					// exit directly if ctrl+c on empty line
					if len(rl.Line().Line) == 0 {
						fmt.Println("Goon REPL interrupted")
						return
					}
					continue
				} else if err == io.EOF {
					fmt.Println("Exiting...")
					return
				}
				log.Printf("read error: %v", err)
				continue
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, ":") {
				if handleBuiltin(line) {
					return
				}
				continue
			}

			if err := h.handleCommand(ctx, line); err != nil {
				log.Printf("command error: %v", err)
			}
		}
	}
}

func handleBuiltin(cmd string) bool {
	switch cmd {
	case ":quit", ":exit":
		fmt.Println("Exiting...")
		return true
	case ":help":
		fmt.Println(`
Available commands:
  explain <query> -> Ask Goon for an explanation about something in the repository
`)
	default:
		fmt.Println("Unknown command. Try :help")
	}
	return false
}
