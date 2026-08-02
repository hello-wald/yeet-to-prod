package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, name, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadConfig_OK(t *testing.T) {
	p := writeTemp(t, "config.json", `{
		"ID": {"timezone":"Asia/Jakarta","weekend_days":[6,0],"work_hours":{"start":9,"end":17},"holidays":["2026-08-17"]}
	}`)
	c, err := loadConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c["ID"]; !ok || c["ID"].WorkHours.End != 17 {
		t.Errorf("parsed config wrong: %+v", c)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	if _, err := loadConfig("does-not-exist.json"); err == nil {
		t.Error("want error for missing file, got nil")
	}
}

func TestLoadConfig_BadJSON(t *testing.T) {
	p := writeTemp(t, "config.json", `{not json`)
	if _, err := loadConfig(p); err == nil {
		t.Error("want error for bad JSON, got nil")
	}
}

func TestLoadConfig_Empty(t *testing.T) {
	p := writeTemp(t, "config.json", `{}`)
	if _, err := loadConfig(p); err == nil {
		t.Error("want error for empty config (no countries), got nil")
	}
}
