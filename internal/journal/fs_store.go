package journal

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	djv1 "dj/gen/go/dj/v1"

	"dj/internal/config"
)

var (
	newMonthFolderPattern    = regexp.MustCompile(`^\d{4}-[A-Z][a-z]{2}$`)
	legacyMonthFolderPattern = regexp.MustCompile(`^([A-Za-z]+)-(\d{4})$`)
)

var legacyMonthNames = map[string]int{
	"jan": 1, "january": 1,
	"feb": 2, "february": 2,
	"mar": 3, "march": 3,
	"apr": 4, "april": 4,
	"may": 5,
	"jun": 6, "june": 6,
	"jul": 7, "july": 7,
	"aug": 8, "august": 8,
	"sep": 9, "september": 9,
	"oct": 10, "october": 10,
	"nov": 11, "november": 11,
	"dec": 12, "december": 12,
}

// FSJournalStore implements JournalStore using the local filesystem.
type FSJournalStore struct {
	config config.Provider
	now    func() time.Time
}

// NewFSStore returns a filesystem-backed journal store.
func NewFSStore(cfg config.Provider) *FSJournalStore {
	return &FSJournalStore{
		config: cfg,
		now:    time.Now,
	}
}

// MonthFolderName returns the journal month folder name for t (YYYY-MMM).
func MonthFolderName(t time.Time) string {
	return t.Format("2006-Jan")
}

// IsNewMonthFolderFormat reports whether name matches YYYY-MMM.
func IsNewMonthFolderFormat(name string) bool {
	return newMonthFolderPattern.MatchString(name)
}

// LegacyMonthFolderTarget converts a legacy Month-Year folder name to YYYY-MMM.
func LegacyMonthFolderTarget(name string) (string, bool) {
	matches := legacyMonthFolderPattern.FindStringSubmatch(name)
	if matches == nil {
		return "", false
	}

	month, ok := legacyMonthNames[strings.ToLower(matches[1])]
	if !ok {
		return "", false
	}

	year, err := time.Parse("2006", matches[2])
	if err != nil {
		return "", false
	}

	t := time.Date(year.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return MonthFolderName(t), true
}

func (s *FSJournalStore) baseDir() string {
	return s.config.JournalDirectory()
}

func (s *FSJournalStore) TodaysPath(_ context.Context) (string, error) {
	return s.dayPath(s.now()), nil
}

func (s *FSJournalStore) EnsureMonthFolder(_ context.Context) error {
	monthFolder := filepath.Join(s.baseDir(), MonthFolderName(s.now()))
	return os.MkdirAll(monthFolder, os.ModePerm)
}

func (s *FSJournalStore) ReadDay(_ context.Context, date time.Time) ([]byte, string, error) {
	path := s.dayPath(date)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, path, err
	}
	return content, path, nil
}

func (s *FSJournalStore) CreateToday(_ context.Context) (string, bool, error) {
	if err := s.EnsureMonthFolder(context.Background()); err != nil {
		return "", false, err
	}

	path := s.dayPath(s.now())
	if _, err := os.Stat(path); err == nil {
		return path, false, ErrDayAlreadyExists
	} else if !os.IsNotExist(err) {
		return "", false, err
	}

	file, err := os.Create(path)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	if _, err := file.WriteString(formatDayHeader(s.now())); err != nil {
		return "", false, err
	}

	return path, true, nil
}

func (s *FSJournalStore) AppendEntry(_ context.Context, path, text string) error {
	return s.appendToFile(path, formatEntry(text, s.now()))
}

func (s *FSJournalStore) AppendBullet(_ context.Context, path, text string) error {
	return s.appendToFile(path, formatBullet(text, s.now()))
}

func (s *FSJournalStore) ListMonthFolders(_ context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.baseDir())
	if err != nil {
		return nil, err
	}

	sorted := SortFoldersForHistory(entries)
	names := make([]string, 0, len(sorted))
	for _, entry := range sorted {
		names = append(names, entry.Name())
	}
	return names, nil
}

func (s *FSJournalStore) ListFiles(_ context.Context, monthFolder string) ([]string, error) {
	dirItems, err := os.ReadDir(filepath.Join(s.baseDir(), monthFolder))
	if err != nil {
		return nil, err
	}

	sorted := SortDescendingByLastModified(dirItems)
	names := make([]string, 0, len(sorted))
	for _, item := range sorted {
		if !item.IsDir() && !strings.HasSuffix(item.Name(), "swp") {
			names = append(names, item.Name())
		}
	}
	return names, nil
}

func (s *FSJournalStore) ReadFile(_ context.Context, monthFolder, fileName string) ([]byte, error) {
	path := filepath.Join(s.baseDir(), monthFolder, fileName)
	return os.ReadFile(path)
}

func (s *FSJournalStore) ListFolderEntries(_ context.Context, folder string) ([]fs.DirEntry, error) {
	return os.ReadDir(folder)
}

func (s *FSJournalStore) PlanFolderMigration(_ context.Context, root string) ([]*djv1.FolderRename, error) {
	var renames []*djv1.FolderRename

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}

		base := d.Name()
		if strings.HasPrefix(base, ".") {
			return fs.SkipDir
		}
		if IsNewMonthFolderFormat(base) {
			return nil
		}

		newName, ok := LegacyMonthFolderTarget(base)
		if !ok {
			return nil
		}

		renames = append(renames, &djv1.FolderRename{
			OldPath: path,
			NewPath: filepath.Join(filepath.Dir(path), newName),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(renames, func(i, j int) bool {
		return len(renames[i].OldPath) > len(renames[j].OldPath)
	})

	return renames, nil
}

func (s *FSJournalStore) MigrateFolders(ctx context.Context, root string, dryRun bool) ([]*djv1.FolderRename, error) {
	renames, err := s.PlanFolderMigration(ctx, root)
	if err != nil {
		return nil, err
	}

	for _, rename := range renames {
		if _, err := os.Stat(rename.NewPath); err == nil {
			return nil, fmt.Errorf("target already exists: %s", rename.NewPath)
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		if dryRun {
			continue
		}

		if err := os.Rename(rename.OldPath, rename.NewPath); err != nil {
			return nil, fmt.Errorf("rename %s -> %s: %w", rename.OldPath, rename.NewPath, err)
		}
	}

	return renames, nil
}

func (s *FSJournalStore) dayPath(date time.Time) string {
	monthPath := filepath.Join(s.baseDir(), MonthFolderName(date))
	fileName := fmt.Sprintf("%s.md", date.Format("02-Monday"))
	return filepath.Join(monthPath, fileName)
}

func (s *FSJournalStore) appendToFile(path, entry string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrJournalNotFound
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(entry); err != nil {
		return err
	}

	return nil
}

// SortDescendingByLastModified sorts directory entries by modification time descending.
func SortDescendingByLastModified(fileItems []fs.DirEntry) []fs.DirEntry {
	sort.SliceStable(fileItems, func(i, j int) bool {
		first, _ := fileItems[i].Info()
		second, _ := fileItems[j].Info()
		return first.ModTime().After(second.ModTime())
	})
	return fileItems
}

// SortFoldersForHistory sorts month folders descending by name, then other dirs by mtime.
func SortFoldersForHistory(fileItems []fs.DirEntry) []fs.DirEntry {
	monthFolders := make([]fs.DirEntry, 0)
	otherFolders := make([]fs.DirEntry, 0)

	for _, item := range fileItems {
		if !item.IsDir() || strings.HasPrefix(item.Name(), ".") {
			continue
		}
		if IsNewMonthFolderFormat(item.Name()) {
			monthFolders = append(monthFolders, item)
		} else {
			otherFolders = append(otherFolders, item)
		}
	}

	sort.SliceStable(monthFolders, func(i, j int) bool {
		return monthFolders[i].Name() > monthFolders[j].Name()
	})

	otherFolders = SortDescendingByLastModified(otherFolders)
	return append(monthFolders, otherFolders...)
}
