package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chronicle/internal/engine"
)

type REPL struct {
	engine *engine.Engine
}

func New(engine *engine.Engine) *REPL {
	return &REPL{engine: engine}
}

func (r *REPL) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Chronicle -- Personal Knowledge Vault")
	fmt.Println("Type `help` to see available commands")
	fmt.Println()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			fmt.Println("\nExiting Chronicle")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if r.handle(line) {
			return
		}
	}
}
