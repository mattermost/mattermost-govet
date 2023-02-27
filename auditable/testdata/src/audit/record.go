// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package audit

// Record provides a consistent set of fields used for all audit logging.
type Record struct {
	EventName string                 `json:"event_name"`
	Status    string                 `json:"status"`
	EventData EventData              `json:"event"`
	Meta      map[string]interface{} `json:"meta"`
}

// EventData contains all event specific data about the modified entity
type EventData struct {
	Parameters  map[string]interface{} `json:"parameters"`      // Payload and parameters being processed as part of the request
	PriorState  map[string]interface{} `json:"prior_state"`     // Prior state of the object being modified, nil if no prior state
	ResultState map[string]interface{} `json:"resulting_state"` // Resulting object after creating or modifying it
	ObjectType  string                 `json:"object_type"`     // String representation of the object type. eg. "post"
}

// Auditable for sensitive object classes, consider implementing Auditable and include whatever the
// AuditableObject returns. For example: it's likely OK to write a user object to the
// audit logs, but not the user password in cleartext or hashed form
type Auditable interface {
	Auditable() map[string]interface{}
}

// Success marks the audit record status as successful.
func (rec *Record) Success() {
	rec.Status = Success
}

// Fail marks the audit record status as failed.
func (rec *Record) Fail() {
	rec.Status = Fail
}

// AddEventParameter adds a parameter, e.g. query or post body, to the event
func (rec *Record) AddEventParameter(key string, val interface{}) {
	if rec.EventData.Parameters == nil {
		rec.EventData.Parameters = make(map[string]interface{})
	}

	if auditableVal, ok := val.(Auditable); ok {
		rec.EventData.Parameters[key] = auditableVal.Auditable()
	} else {
		rec.EventData.Parameters[key] = val
	}
}

// AddEventResultState adds the result state of the modified object to the audit record
func (rec *Record) AddEventResultState(object Auditable) {
	rec.EventData.ResultState = object.Auditable()
}

// AddEventObjectType adds the object type of the modified object to the audit record
func (rec *Record) AddEventObjectType(objectType string) {
	rec.EventData.ObjectType = objectType
}
