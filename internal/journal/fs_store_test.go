package journal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMonthFolderName(t *testing.T) {
	tests := []struct {
		date     time.Time
		expected string
	}{
		{time.Date(2024, 3, 27, 0, 0, 0, 0, time.UTC), "2024-Mar"},
		{time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC), "2026-Sep"},
		{time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC), "2022-Feb"},
	}

	for _, tt := range tests {
		if got := MonthFolderName(tt.date); got != tt.expected {
			t.Fatalf("MonthFolderName(%v) = %q, want %q", tt.date, got, tt.expected)
		}
	}
}

func TestIsNewMonthFolderFormat(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"2024-Mar", true},
		{"2026-Sep", true},
		{"March-2024", false},
		{"Jan-2023", false},
		{"Landmarks", false},
	}

	for _, tt := range tests {
		if got := IsNewMonthFolderFormat(tt.name); got != tt.expected {
			t.Fatalf("IsNewMonthFolderFormat(%q) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestLegacyMonthFolderTarget(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		ok       bool
	}{
		{"March-2024", "2024-Mar", true},
		{"Jan-2023", "2023-Jan", true},
		{"feb-2022", "2022-Feb", true},
		{"September-2026", "2026-Sep", true},
		{"Landmarks", "", false},
		{"2024-Mar", "", false},
	}

	for _, tt := range tests {
		got, ok := LegacyMonthFolderTarget(tt.name)
		if got != tt.expected || ok != tt.ok {
			t.Fatalf("LegacyMonthFolderTarget(%q) = (%q, %v), want (%q, %v)", tt.name, got, ok, tt.expected, tt.ok)
		}
	}
}

func TestPlanFolderMigration(t *testing.T) {
	root := t.TempDir()

	legacyDirs := []string{"March-2024", "Jan-2023", "feb-2022", "Landmarks"}
	for _, dir := range legacyDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	store := &FSJournalStore{}
	renames, err := store.PlanFolderMigration(t.Context(), root)
	if err != nil {
		t.Fatal(err)
	}

	if len(renames) != 3 {
		t.Fatalf("expected 3 renames, got %d", len(renames))
	}

	targets := map[string]bool{
		"2024-Mar": true,
		"2023-Jan": true,
		"2022-Feb": true,
	}

	for _, rename := range renames {
		base := filepath.Base(rename.NewPath)
		if !targets[base] {
			t.Fatalf("unexpected rename target %q", base)
		}
	}
}

func TestMigrateFoldersDryRun(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "March-2024"), 0o755); err != nil {
		t.Fatal(err)
	}

	store := &FSJournalStore{}
	renames, err := store.MigrateFolders(t.Context(), root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(renames) != 1 {
		t.Fatalf("expected 1 rename, got %d", len(renames))
	}

	if _, err := os.Stat(filepath.Join(root, "March-2024")); err != nil {
		t.Fatal("legacy folder should still exist after dry run")
	}
}

func TestMigrateFoldersApply(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "March-2024"), 0o755); err != nil {
		t.Fatal(err)
	}

	store := &FSJournalStore{}
	renames, err := store.MigrateFolders(t.Context(), root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(renames) != 1 {
		t.Fatalf("expected 1 rename, got %d", len(renames))
	}

	if _, err := os.Stat(filepath.Join(root, "2024-Mar")); err != nil {
		t.Fatal("expected migrated folder 2024-Mar")
	}
}
