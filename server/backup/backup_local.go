package backup

import (
	"errors"
	"io"
	"os"

	"github.com/mholt/archiver/v3"
	"github.com/pterodactyl/wings/system"
)

type LocalBackup struct {
	Backup
}

var _ BackupInterface = (*LocalBackup)(nil)

func NewLocal(uuid string, ignore string) *LocalBackup {
	return &LocalBackup{
		Backup{
			Uuid:    uuid,
			Ignore:  ignore,
			adapter: LocalBackupAdapter,
		},
	}
}

// LocateLocal finds the backup for a server and returns the local path. This
// will obviously only work if the backup was created as a local backup.
func LocateLocal(uuid string) (*LocalBackup, os.FileInfo, error) {
	b := &LocalBackup{
		Backup{
			Uuid:   uuid,
			Ignore: "",
		},
	}

	st, err := os.Stat(b.Path())
	if err != nil {
		return nil, nil, err
	}

	if st.IsDir() {
		return nil, nil, errors.New("invalid archive, is directory")
	}

	return b, st, nil
}

// Remove removes a backup from the system.
func (b *LocalBackup) Remove() error {
	return os.Remove(b.Path())
}

// WithLogContext attaches additional context to the log output for this backup.
func (b *LocalBackup) WithLogContext(c map[string]interface{}) {
	b.logContext = c
}

// Generate generates a backup of the selected files and pushes it to the
// defined location for this instance.
func (b *LocalBackup) Generate(basePath, ignore string) (*ArchiveDetails, error) {
	a := &Archive{
		BasePath: basePath,
		Ignore:   ignore,
	}

	b.log().Info("creating backup for server...")
	if err := a.Create(b.Path()); err != nil {
		return nil, err
	}
	b.log().Info("created backup successfully")

	return b.Details(), nil
}

// Restore will walk over the archive and call the callback function for each
// file encountered.
func (b *LocalBackup) Restore(_ io.Reader, callback RestoreCallback) error {
	return archiver.Walk(b.Path(), func(f archiver.File) error {
		if f.IsDir() {
			return nil
		}
		name, err := system.ExtractArchiveSourceName(f, "/")
		if err != nil {
			return err
		}
		return callback(name, f)
	})
}
