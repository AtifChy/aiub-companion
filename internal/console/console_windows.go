// Package console provides functionality to attach the application to the parent console on Windows systems.
package console

import (
	"os"

	"golang.org/x/sys/windows"
)

const ATTACH_PARENT_PROCESS = ^uint32(0) // (DWORD)-1

var (
	kernel32DLL       = windows.NewLazySystemDLL("kernel32.dll")
	attachConsoleProc = kernel32DLL.NewProc("AttachConsole")
)

func attachConsole() {
	r1, _, _ := attachConsoleProc.Call(uintptr(ATTACH_PARENT_PROCESS))
	if r1 == 0 {
		// No parent console or failed to attach
		return
	}

	if f, err := os.OpenFile("CONOUT$", os.O_RDWR, 0); err == nil {
		os.Stdout = f
		os.Stderr = f
	}
}

func init() {
	attachConsole()
}
