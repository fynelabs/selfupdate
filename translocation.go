package selfupdate

import (
	"errors"
	"runtime"
	"strings"
)

// ErrTranslocated is returned when the application is running from a macOS App Translocation
// mount. Gatekeeper runs a quarantined application from a randomised read-only location until the
// user moves it with Finder, so no update can be written next to the executable. The user has to
// quit the application, move it with Finder (for example into their Applications folder) and open
// it from there before updates can be applied.
var ErrTranslocated = errors.New("the application is running from a read-only macOS App Translocation location, quit it and move it with Finder (for example into your Applications folder) before it can update itself")

// isTranslocatedPath reports whether path points inside a macOS App Translocation mount.
func isTranslocatedPath(path string) bool {
	return strings.Contains(path, "/AppTranslocation/")
}

// checkTranslocated returns ErrTranslocated when the current executable runs from an App
// Translocation mount. It is a no-op on other operating systems.
func checkTranslocated() error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	exe, err := ExecutableRealPath()
	if err != nil {
		return err
	}
	if isTranslocatedPath(exe) {
		return ErrTranslocated
	}
	return nil
}
