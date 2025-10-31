// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
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
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"context"

	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/go-proton-api"
	"github.com/bradenaw/juniper/xslices"
	"github.com/sirupsen/logrus"
)

type MetadataFileChecker interface {
	HasMessage(msgID string) (bool, error)
}

type MetadataStage struct {
	client    apiclient.Client
	log       *logrus.Entry
	outputCh  chan []proton.MessageMetadata
	pageSize  int
	splitSize int
	filter    *Filter // Filter for messages (nil = no filtering)
}

func NewMetadataStage(
	client apiclient.Client,
	entry *logrus.Entry,
	pageSize int,
	splitSize int,
	filter *Filter,
) *MetadataStage {
	return &MetadataStage{
		client:    client,
		log:       entry.WithField("stage", "metadata"),
		outputCh:  make(chan []proton.MessageMetadata),
		pageSize:  pageSize,
		splitSize: splitSize,
		filter:    filter,
	}
}

func (m *MetadataStage) Run(
	ctx context.Context,
	errReporter StageErrorReporter,
	mfc MetadataFileChecker,
	reporter Reporter,
) {
	m.log.Debug("Starting")
	defer m.log.Debug("Exiting")
	defer close(m.outputCh)

	client := m.client

	// Determine filter strategy
	var serverFilter *proton.MessageFilter
	needsClientFiltering := false

	if m.filter != nil && !m.filter.IsEmpty() {
		serverFilter = m.filter.ToServerFilter()
		needsClientFiltering = m.filter.NeedsClientFiltering()

		if serverFilter != nil {
			m.log.Info("Using server-side filtering")
		}
		if needsClientFiltering {
			m.log.Info("Using client-side filtering")
		}
	}

	var lastMessageID string

	for {
		if ctx.Err() != nil {
			return
		}

		var metadata []proton.MessageMetadata

		// Build the message filter for this page
		var pageFilter proton.MessageFilter
		if serverFilter != nil {
			pageFilter = *serverFilter
		} else {
			pageFilter = proton.MessageFilter{
				Desc: true,
			}
		}

		if lastMessageID != "" {
			pageFilter.EndID = lastMessageID
		}

		meta, err := client.GetMessageMetadataPage(ctx, 0, m.pageSize, pageFilter)
		if err != nil {
			errReporter.ReportStageError(err)
			return
		}

		// If there's only one message and it matches EndID, skip it (pagination overlap)
		if lastMessageID != "" && len(meta) != 0 && meta[0].ID == lastMessageID {
			meta = meta[1:]
		}

		metadata = meta

		// Nothing left to do
		if len(metadata) == 0 {
			return
		}

		lastMessageID = metadata[len(metadata)-1].ID

		initialLen := len(metadata)
		metadata = xslices.Filter(metadata, func(t proton.MessageMetadata) bool {
			isPresent, err := mfc.HasMessage(t.ID)
			if err != nil {
				errReporter.ReportStageError(err)
				return false
			}

			// Skip if already present
			if isPresent {
				return false
			}

			// Apply client-side filtering if needed
			if needsClientFiltering {
				return m.filter.MatchesMetadata(t)
			}

			return true
		})

		if len(metadata) != initialLen {
			reporter.OnProgress(initialLen - len(metadata))
		}

		if len(metadata) == 0 {
			continue
		}

		for _, chunk := range xslices.Chunk(metadata, m.splitSize) {
			select {
			case <-ctx.Done():
				return
			case m.outputCh <- chunk:
			}
		}
	}
}

type alwaysMissingMetadataFileChecker struct{}

func (a alwaysMissingMetadataFileChecker) HasMessage(string) (bool, error) {
	return false, nil
}
