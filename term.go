package main

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func getTerminalWidth() int {
	term, err := os.Open("/dev/tty")
	if err != nil {
		term = os.Stdin
	}

	window := struct {
		Rows    uint16
		Columns uint16
		X       uint16
		Y       uint16
	}{}

	result, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		term.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&window)),
	)
	if int(result) == -1 || err != nil {
		return 0
	}

	return int(window.Columns)
}

func getLineSuffix(width, length int) string {
	suffix := width - length
	if suffix > 0 {
		return strings.Repeat(" ", suffix)
	}

	return ""
}
