package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_creates_config(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".djconfig")

	cfg := NewViperProviderWithConfigDir(dir)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load failed: %s", err)
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load did not create the default config file")
	}

	expected := filepath.Join(dir, "Documents", "Dev-Journal")
	if cfg.JournalDirectory() != expected {
		t.Fatalf("JournalDirectory() = %q, want %q", cfg.JournalDirectory(), expected)
	}
}

func TestSetJournalDirectory_updates_config_file(t *testing.T) {
	dir := t.TempDir()

	cfg := NewViperProviderWithConfigDir(dir)
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load failed: %s", err)
	}

	directory := "/my/new/directory"
	if err := cfg.SetJournalDirectory(directory); err != nil {
		t.Fatalf("SetJournalDirectory failed: %s", err)
	}

	if cfg.JournalDirectory() != directory {
		t.Fatalf("JournalDirectory() = %q, want %q", cfg.JournalDirectory(), directory)
	}

	content, err := os.ReadFile(cfg.ConfigFilePath())
	if err != nil {
		t.Fatalf("ReadFile failed: %s", err)
	}

	if string(content) == "" {
		t.Fatal("expected config file to contain updated directory")
	}
}

func TestLoad_uses_real_home_by_default(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir failed: %s", err)
	}

	cfg := NewViperProvider()
	if cfg.configDir != "" {
		t.Fatal("expected empty configDir for default provider")
	}

	dir, err := cfg.configDirPath()
	if err != nil {
		t.Fatalf("configDirPath failed: %s", err)
	}

	if dir != home {
		t.Fatalf("configDirPath() = %q, want %q", dir, home)
	}
}
