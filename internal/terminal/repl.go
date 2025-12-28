package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chronicle/internal/config"
	"chronicle/internal/engine"
)

type REPL struct {
	engine  *engine.Engine
	config  config.Config
	version string
}

func New(engine *engine.Engine, cfg config.Config, version string) *REPL {
	return &REPL{
		engine:  engine,
		config:  cfg,
		version: version,
	}
}

func (r *REPL) Start() {
	clearScreen()
	if r.config.ShowBanner {
		printBanner(r.version)
	}

	pageSize = r.config.PageSize

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
			fmt.Println("GooBye...")
			return
		}
	}
}
