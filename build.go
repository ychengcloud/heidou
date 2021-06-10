//+build linux darwin

package heidou

import "syscall"

func Umask(mask int) (oldmask int) {
	return syscall.Umask(mask)
}
