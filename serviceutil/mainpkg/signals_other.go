//go:build !linux
// +build !linux

package mainpkg

import (
	"os"
	"os/signal"
)

func signals() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}
