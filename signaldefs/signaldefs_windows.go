package signaldefs

import "syscall"

const (
	SIGQUIT = syscall.SIGQUIT
	SIGINT  = syscall.SIGINT
	SIGTERM = syscall.SIGTERM
	SIGUSR1 = SIG_NONE // none
	SIGUSR2 = SIG_NONE // none
)
