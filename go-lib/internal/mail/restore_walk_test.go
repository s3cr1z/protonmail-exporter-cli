// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
//
// Proton Export Tool is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Export Tool is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProtonMail/export-tool/internal/utils"
	"github.com/ProtonMail/go-proton-api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestWalkBackupDir_FindsEMLMessages(t *testing.T) {
	backupDir := t.TempDir()

	// Create a message with EML file
	msgID := "msg1"
	emlPath := filepath.Join(backupDir, msgID+emlExtension)
	metadataPath := filepath.Join(backupDir, msgID+jsonMetadataExtension)

	// Create EML file
	require.NoError(t, os.WriteFile(emlPath, []byte("EML content"), 0o600))

	// Create metadata file
	metadata := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{
			ID:   msgID,
			Time: 123456,
		},
		WriterType: MessageWriterTypeDecryptedAndBuilt,
	}
	metadataBytes, err := metadata.toBytes()
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(metadataPath, metadataBytes, 0o600))

	// Create restore task
	restoreTask := &RestoreTask{
		ctx:       context.Background(),
		backupDir: backupDir,
		log:       logrus.WithField("test", "test"),
	}

	// Walk and collect paths
	var foundPaths []string
	err = restoreTask.walkBackupDir(func(emlPath string) {
		foundPaths = append(foundPaths, emlPath)
	})
	require.NoError(t, err)

	// Should find the message
	require.Len(t, foundPaths, 1)
	require.Equal(t, emlPath, foundPaths[0])
}

func TestWalkBackupDir_FindsFailedToAssembleMessages(t *testing.T) {
	backupDir := t.TempDir()

	// Create a message that failed to assemble (no EML, just directory with body/attachments)
	msgID := "msg2"
	msgDir := filepath.Join(backupDir, msgID)
	metadataPath := filepath.Join(backupDir, msgID+jsonMetadataExtension)

	// Create message directory with body file
	require.NoError(t, os.MkdirAll(msgDir, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(msgDir, "body.txt"), []byte("body content"), 0o600))

	// Create metadata file
	metadata := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{
			ID:   msgID,
			Time: 123456,
		},
		WriterType: MessageWriterTypeFailedToAssemble,
	}
	metadataBytes, err := metadata.toBytes()
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(metadataPath, metadataBytes, 0o600))

	// Create restore task
	restoreTask := &RestoreTask{
		ctx:       context.Background(),
		backupDir: backupDir,
		log:       logrus.WithField("test", "test"),
	}

	// Walk and collect paths
	var foundPaths []string
	err = restoreTask.walkBackupDir(func(emlPath string) {
		foundPaths = append(foundPaths, emlPath)
	})
	require.NoError(t, err)

	// Should find the message (even without EML file)
	require.Len(t, foundPaths, 1)
	// The path will be the synthetic EML path
	expectedPath := filepath.Join(backupDir, msgID+emlExtension)
	require.Equal(t, expectedPath, foundPaths[0])
}

func TestWalkBackupDir_FindsBothMessageTypes(t *testing.T) {
	backupDir := t.TempDir()

	// Create message 1 with EML
	msg1ID := "msg1"
	msg1EMLPath := filepath.Join(backupDir, msg1ID+emlExtension)
	msg1MetadataPath := filepath.Join(backupDir, msg1ID+jsonMetadataExtension)
	require.NoError(t, os.WriteFile(msg1EMLPath, []byte("EML content"), 0o600))
	metadata1 := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{ID: msg1ID, Time: 100000},
		WriterType:      MessageWriterTypeDecryptedAndBuilt,
	}
	metadataBytes1, _ := metadata1.toBytes()
	require.NoError(t, os.WriteFile(msg1MetadataPath, metadataBytes1, 0o600))

	// Create message 2 without EML (failed to assemble)
	msg2ID := "msg2"
	msg2Dir := filepath.Join(backupDir, msg2ID)
	msg2MetadataPath := filepath.Join(backupDir, msg2ID+jsonMetadataExtension)
	require.NoError(t, os.MkdirAll(msg2Dir, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(msg2Dir, "body.pgp"), []byte("encrypted body"), 0o600))
	metadata2 := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{ID: msg2ID, Time: 200000},
		WriterType:      MessageWriterTypeNoAddrKey,
	}
	metadataBytes2, _ := metadata2.toBytes()
	require.NoError(t, os.WriteFile(msg2MetadataPath, metadataBytes2, 0o600))

	// Create restore task
	restoreTask := &RestoreTask{
		ctx:       context.Background(),
		backupDir: backupDir,
		log:       logrus.WithField("test", "test"),
	}

	// Walk and collect paths
	var foundPaths []string
	err := restoreTask.walkBackupDir(func(emlPath string) {
		foundPaths = append(foundPaths, emlPath)
	})
	require.NoError(t, err)

	// Should find both messages
	require.Len(t, foundPaths, 2)
	// Paths could be in any order, so check they're both present
	require.Contains(t, foundPaths, msg1EMLPath)
	require.Contains(t, foundPaths, filepath.Join(backupDir, msg2ID+emlExtension))
}

func TestWalkBackupDir_SkipsSubdirectories(t *testing.T) {
	backupDir := t.TempDir()

	// Create a properly structured message in root
	msg1ID := "msg1"
	msg1MetadataPath := filepath.Join(backupDir, msg1ID+jsonMetadataExtension)
	metadata1 := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{ID: msg1ID, Time: 100000},
		WriterType:      MessageWriterTypeDecryptedAndBuilt,
	}
	metadataBytes1, _ := metadata1.toBytes()
	require.NoError(t, os.WriteFile(msg1MetadataPath, metadataBytes1, 0o600))

	// Create a timestamped subdirectory with a message (should be skipped)
	subDir := filepath.Join(backupDir, "mail_20241031_074412")
	require.NoError(t, os.MkdirAll(subDir, 0o700))
	msg2ID := "msg2"
	msg2MetadataPath := filepath.Join(subDir, msg2ID+jsonMetadataExtension)
	metadata2 := MessageMetadata{
		MessageMetadata: proton.MessageMetadata{ID: msg2ID, Time: 200000},
		WriterType:      MessageWriterTypeDecryptedAndBuilt,
	}
	metadataBytes2, _ := metadata2.toBytes()
	require.NoError(t, os.WriteFile(msg2MetadataPath, metadataBytes2, 0o600))

	// Create restore task
	restoreTask := &RestoreTask{
		ctx:       context.Background(),
		backupDir: backupDir,
		log:       logrus.WithField("test", "test"),
	}

	// Walk and collect paths
	var foundPaths []string
	err := restoreTask.walkBackupDir(func(emlPath string) {
		foundPaths = append(foundPaths, emlPath)
	})
	require.NoError(t, err)

	// Should only find the message in root, not in subdirectory
	require.Len(t, foundPaths, 1)
	expectedPath := filepath.Join(backupDir, msg1ID+emlExtension)
	require.Equal(t, expectedPath, foundPaths[0])
}

func TestGetTimestampedBackupDirs(t *testing.T) {
	backupDir := t.TempDir()

	// Create valid timestamped directories
	dir1 := filepath.Join(backupDir, "mail_20241031_074412")
	dir2 := filepath.Join(backupDir, "mail_20231225_120000")
	require.NoError(t, os.MkdirAll(dir1, 0o700))
	require.NoError(t, os.MkdirAll(dir2, 0o700))

	// Create invalid directories (should be ignored)
	invalidDir1 := filepath.Join(backupDir, "invalid_dir")
	invalidDir2 := filepath.Join(backupDir, "mail_20241031") // wrong format
	invalidDir3 := filepath.Join(backupDir, "mail_notadate_074412")
	require.NoError(t, os.MkdirAll(invalidDir1, 0o700))
	require.NoError(t, os.MkdirAll(invalidDir2, 0o700))
	require.NoError(t, os.MkdirAll(invalidDir3, 0o700))

	// Create a file (should be ignored)
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, "somefile.txt"), []byte("content"), 0o600))

	// Create restore task
	restoreTask := &RestoreTask{
		ctx:       context.Background(),
		backupDir: backupDir,
		log:       logrus.WithField("test", "test"),
	}

	// Get timestamped directories
	dirs, err := restoreTask.getTimestampedBackupDirs()
	require.NoError(t, err)

	// Should find exactly 2 valid timestamped directories
	require.Len(t, dirs, 2)
	require.Contains(t, dirs, dir1)
	require.Contains(t, dirs, dir2)
}
