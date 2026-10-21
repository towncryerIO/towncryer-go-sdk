package models

// ScheduleInfo is the response from Messages.SendBulkMessages.
//
// The wire field is "jobId", not "id" - the committed OpenAPI spec disagrees with the
// backend's actual struct tag, so this is modeled against the real backend source.
type ScheduleInfo struct {
	JobId  string `json:"jobId"`
	Status string `json:"status"`
}
