package selfupdate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTranslocatedPath(t *testing.T) {
	assert.True(t, isTranslocatedPath("/private/var/folders/pc/1mt578kn6615qkxsygpwrw7m0000gn/T/AppTranslocation/CC9EEF5E-BE8A-47B5-9BA2-2A9793595220/d/Example.app/Contents/MacOS/example"))
	assert.False(t, isTranslocatedPath("/Applications/Example.app/Contents/MacOS/example"))
	assert.False(t, isTranslocatedPath("/Users/me/Applications/Example.app/Contents/MacOS/example"))
	assert.False(t, isTranslocatedPath("/home/me/AppTranslocation"))
	assert.False(t, isTranslocatedPath(""))
}

func TestApplyTranslocatedReadOnly(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("relies on POSIX directory permissions being enforced")
	}

	dir := filepath.Join(t.TempDir(), "AppTranslocation", "CC9EEF5E", "d", "Example.app", "Contents", "MacOS")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "example")
	if err := os.WriteFile(target, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0755) })

	err := Apply(strings.NewReader("new"), Options{TargetPath: target})
	assert.ErrorIs(t, err, ErrTranslocated)
}

func TestApplyReadOnlyNotTranslocated(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("relies on POSIX directory permissions being enforced")
	}

	dir := filepath.Join(t.TempDir(), "Example.app", "Contents", "MacOS")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "example")
	if err := os.WriteFile(target, []byte("old"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0755) })

	err := Apply(strings.NewReader("new"), Options{TargetPath: target})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrTranslocated)
}
