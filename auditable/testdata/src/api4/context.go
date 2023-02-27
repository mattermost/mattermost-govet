// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import "audit"

type Context struct {
}

// LogAuditRec logs an audit record using default LevelAPI.
func (c *Context) LogAuditRec(rec *audit.Record) {
}

// MakeAuditRecord creates a audit record pre-populated with data from this context.
func (c *Context) MakeAuditRecord(event string, initialStatus string) *audit.Record {
	rec := &audit.Record{
		EventName: event,
		Status:    initialStatus,
	}

	return rec
}
