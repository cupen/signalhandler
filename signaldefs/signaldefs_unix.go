package signaldefs

import "syscall"

const (
	SIGQUIT = syscall.SIGQUIT
	SIGINT  = syscall.SIGINT
	SIGTERM = syscall.SIGTERM

	SIGUSR1 = syscall.SIGUSR1
	SIGUSR2 = syscall.SIGUSR2
)
