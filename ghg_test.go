package ghg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtract(t *testing.T) {
	const content = "executable contents"
	for _, ext := range []string{".zip", ".tar.gz", ".tgz"} {
		t.Run(ext, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "release"+ext)
			f, err := os.Create(src)
			if err != nil {
				t.Fatal(err)
			}
			if ext == ".zip" {
				w := zip.NewWriter(f)
				header := &zip.FileHeader{Name: "release/ghg"}
				header.SetMode(0755)
				entry, err := w.CreateHeader(header)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := entry.Write([]byte(content)); err != nil {
					t.Fatal(err)
				}
				if err := w.Close(); err != nil {
					t.Fatal(err)
				}
			} else {
				gz := gzip.NewWriter(f)
				w := tar.NewWriter(gz)
				if err := w.WriteHeader(&tar.Header{
					Name: "release/ghg",
					Mode: 0755,
					Size: int64(len(content)),
				}); err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write([]byte(content)); err != nil {
					t.Fatal(err)
				}
				if err := w.Close(); err != nil {
					t.Fatal(err)
				}
				if err := gz.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}

			dest := filepath.Join(dir, "extracted")
			if err := extract(src, dest); err != nil {
				t.Fatal(err)
			}
			bin := filepath.Join(dest, "release", "ghg")
			got, err := os.ReadFile(bin)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != content {
				t.Errorf("extracted content = %q, want %q", got, content)
			}
			info, err := os.Stat(bin)
			if err != nil {
				t.Fatal(err)
			}
			if !isWindows && info.Mode().Perm() != 0755 {
				t.Errorf("extracted mode = %o, want 755", info.Mode().Perm())
			}
		})
	}
}
