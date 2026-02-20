package termutil

import (
	"math"
	"os"

	"golang.org/x/term"
)

func IsTerminal(file *os.File) bool {
	if file == nil {
		return false
	}

	fd := file.Fd()
	if fd > uintptr(math.MaxInt) {
		return false
	}

	return term.IsTerminal(int(fd))
}
