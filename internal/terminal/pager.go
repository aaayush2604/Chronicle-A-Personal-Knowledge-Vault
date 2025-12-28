package terminal

import (
	"bufio"
	"os"
)

var pageSize int

func pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
