//go:build windows

package heidou

func Umask(mask int) (oldmask int) {
	return 0
}
