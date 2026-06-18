//go:build windows

package main

import (
	"os"
	"syscall"
)

var (
	modkernel32       = syscall.NewLazyDLL("kernel32.dll")
	procAttachConsole = modkernel32.NewProc("AttachConsole")
)

func attachConsole() {
	const attachParentProcess = ^uintptr(0)
	r0, _, _ := procAttachConsole.Call(attachParentProcess)
	if r0 != 0 {
		if hOut, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE); err == nil {
			os.Stdout = os.NewFile(uintptr(hOut), "/dev/stdout")
		}
		if hErr, err := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE); err == nil {
			os.Stderr = os.NewFile(uintptr(hErr), "/dev/stderr")
		}
		if hIn, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE); err == nil {
			os.Stdin = os.NewFile(uintptr(hIn), "/dev/stdin")
		}
	}
}
