//go:build !windows
// +build !windows

package selfupdate

func hideFile(_ string) error {
	return nil
}
