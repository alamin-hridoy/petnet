//go:build linux
// +build linux

package mainpkg

import (
	"os"
	"os/signal"
	"syscall"
)

func signals() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
