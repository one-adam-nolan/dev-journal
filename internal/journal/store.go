package journal

import (
	"context"
	"io/fs"
	"time"

	djv1 "dj/gen/go/dj/v1"
)

// JournalStore abstracts journal file storage operations.
type JournalStore interface {
	TodaysPath(ctx context.Context) (string, error)
	EnsureMonthFolder(ctx context.Context) error
	ReadDay(ctx context.Context, date time.Time) ([]byte, string, error)
	CreateToday(ctx context.Context) (path string, created bool, err error)
	AppendEntry(ctx context.Context, path, text string) error
	AppendBullet(ctx context.Context, path, text string) error
	ListMonthFolders(ctx context.Context) ([]string, error)
	ListFiles(ctx context.Context, monthFolder string) ([]string, error)
	ReadFile(ctx context.Context, monthFolder, fileName string) ([]byte, error)
	ListFolderEntries(ctx context.Context, folder string) ([]fs.DirEntry, error)
	PlanFolderMigration(ctx context.Context, root string) ([]*djv1.FolderRename, error)
	MigrateFolders(ctx context.Context, root string, dryRun bool) ([]*djv1.FolderRename, error)
}
