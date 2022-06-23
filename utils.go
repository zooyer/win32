package win32

import (
	"syscall"
)

func toBool(i int) bool {
	return i != 0
}

func fromBool(b bool) int {
	if b {
		return 1
	}

	return 0
}

func fromError(err error) error {
	if err == syscall.Errno(0) {
		return nil
	}

	return err
}
