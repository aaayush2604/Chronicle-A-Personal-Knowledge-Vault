package terminal

import (
	"bufio"
	"os"
)

const pageSize = 10

func pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
