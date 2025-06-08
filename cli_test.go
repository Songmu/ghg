package ghg

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var tempdir string

func TestMain(m *testing.M) {
	tempd, err := os.MkdirTemp("", "ghgtest-")
	if err != nil {
		panic("fialed to create tempdir in test")
	}
	os.Setenv(EnvHome, tempd)
	tempdir = tempd
	exit := m.Run()
	os.RemoveAll(tempdir)
	os.Exit(exit)
}

var tests = []struct {
	name     string
	input    string
	out      string
	exitCode int
}{
	{
		name:     "simple",
		input:    "Songmu/ghg@v0.0.1",
		out:      "ghg",
		exitCode: 0,
	},
	{
		name:     "naked binary",
		input:    "bcicen/ctop@v0.4.1",
		out:      "ctop",
		exitCode: 0,
	},
}

func TestGet(t *testing.T) {
	for _, tt := range tests {
		exitCode := (&CLI{
			ErrStream: io.Discard,
			OutStream: io.Discard,
		}).Run([]string{
			"get",
			tt.input,
		})
		if exitCode != tt.exitCode {
			t.Errorf("%s(exitCode): out=%d want=%d", tt.name, exitCode, tt.exitCode)
		}
		expectFile := tt.out
		if runtime.GOOS == "windows" {
			expectFile += ".exe"
		}
		fname := filepath.Join(tempdir, "bin", expectFile)
		fi, err := os.Stat(fname)
		if err != nil {
			t.Errorf("%s(exists): %q shoud be exists, but not found: %s", tt.name, expectFile, err)
		}
		if (fi.Mode() & 0111) == 0 {
			t.Errorf("%s(executable): %q shoud be executable, but not.", tt.name, expectFile)
		}
	}
}
