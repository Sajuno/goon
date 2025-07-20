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
		HistoryFile:     "/tmp/goon_history.tmp", // TODO: save history to postgres somehow
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("failed to start REPL: %v", err)
	}
	defer rl.Close()

	fmt.Println("Goon REPL is ready. Type ':help' or enter a command.")

	h := newCommandHandler(ag)
	inputChan := make(chan string)
	errChan := make(chan error)

	go func() {
		for {
			line, err := rl.Readline()
			if err != nil {
				errChan <- err
				return
			}
			inputChan <- strings.TrimSpace(line)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, exiting..")
			return
		case err := <-errChan:
			if errors.Is(err, readline.ErrInterrupt) || err == io.EOF {
				fmt.Println("Goon REPL interrupted")
				return
			}
			log.Printf("error: %v", err)
			return
		case line := <-inputChan:
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, ":") {
				if handleBuiltin(line) {
					return
				}
				continue
			}
			err = h.handleCommand(ctx, line)
			if err != nil {
				errChan <- err
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
