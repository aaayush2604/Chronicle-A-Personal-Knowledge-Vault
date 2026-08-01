package terminal

import (
	"fmt"
	"strings"

	"chronicle/internal/config"
	"chronicle/internal/engine"

	"github.com/ergochat/readline"
)

type REPL struct {
	engine   *engine.Engine
	config   config.Config
	version  string
	terminal *readline.Instance
	prompt   string
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

	fmt.Println("Chronicle -- Personal Knowledge Vault")
	fmt.Println("Type `help` to see available commands")
	fmt.Println()

	completer := NewCompleter()
	newPrompt := fgCyan + "chronicle> " + reset
	r.prompt = newPrompt

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       newPrompt,
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	r.terminal = rl
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("\nExiting Chronicle")
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if r.handle(line) {
			fmt.Println("GoodBye...")
			return
		}
	}

}
