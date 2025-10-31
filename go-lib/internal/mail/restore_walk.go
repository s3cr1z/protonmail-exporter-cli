package mail

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func (r *RestoreTask) walkBackupDir(fn func(emlPath string)) error {
	return filepath.Walk(r.backupDir, func(path string, info fs.FileInfo, err error) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		if err != nil {
			logrus.WithError(err).WithField("path", path).Warn("Cannot inspect path. Skipping.")
			return nil
		}

		if info.IsDir() {
			if path != r.backupDir {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process metadata files
		if !strings.HasSuffix(path, jsonMetadataExtension) {
			return nil
		}

		// Convert metadata path to what the EML path would be for compatibility
		// with existing callers that expect EML paths
		metadataPath := path
		emlPath := strings.TrimSuffix(metadataPath, jsonMetadataExtension) + emlExtension

		fn(emlPath)

		return nil
	})
}

func (r *RestoreTask) getTimestampedBackupDirs() ([]string, error) {
	var result []string
	err := filepath.Walk(r.backupDir, func(path string, info fs.FileInfo, err error) error {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
		}

		if err != nil {
			return nil //nolint:nilerr // we proceed in case of errors
		}

		name := info.Name()
		if (err != nil) || !info.IsDir() || (path == r.backupDir) {
			return nil //nolint:nilerr // ignore errors, files, and the walk's root folder
		}

		if mailFolderRegExp.MatchString(name) {
			result = append(result, filepath.Join(r.backupDir, name))
		}

		return fs.SkipDir // we do not recurse into dirs
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
