package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnv(t *testing.T) {
	t.Setenv("FOO", "bar")
	if env("FOO", "def") != "bar" {
		t.Error("env should return set value")
	}
	if env("MISSING_XYZ", "def") != "def" {
		t.Error("env should return default when unset")
	}
	t.Setenv("EMPTY", "")
	if env("EMPTY", "def") != "def" {
		t.Error("env should treat empty as unset → default")
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("N", "42")
	if envInt("N", 7) != 42 {
		t.Error("envInt should parse int")
	}
	t.Setenv("BAD", "notanint")
	if envInt("BAD", 7) != 7 {
		t.Error("envInt should fall back to default on parse error")
	}
	if envInt("MISSING_N", 7) != 7 {
		t.Error("envInt should default when unset")
	}
}

func TestLoadDotenv(t *testing.T) {
	body := "# comment\nA=1\nB = \"two\"\n\nEXISTING=fromfile\n"
	p := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("EXISTING", "preset")
	loadDotenv(p)

	if os.Getenv("A") != "1" {
		t.Errorf("A = %q, want 1", os.Getenv("A"))
	}
	if os.Getenv("B") != "two" {
		t.Errorf("B = %q, want two (quotes + spaces trimmed)", os.Getenv("B"))
	}
	if os.Getenv("EXISTING") != "preset" {
		t.Error("loadDotenv must NOT override an already-set var")
	}
}

func TestLoadDotenv_MissingFileNoError(t *testing.T) {
	loadDotenv("no-such-file.env") // must not panic
}
