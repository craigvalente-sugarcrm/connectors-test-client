package common

import "time"

// Projection represents the set of fields that are common between Sugar
// and a remote calendar provider.
type Projection struct {
	ID          bool      `json:"id,omitempty"`
	Type        string    `json:"type,omitempty"` // create|update|delete
	Subject     string    `json:"subject,omitempty"`
	Description string    `json:"description,omitempty"`
	Owner       string    `json:"owner,omitempty"`
	Start       time.Time `json:"start,omitempty"`
	End         time.Time `json:"end,omitempty"`
	Location    string    `json:"location,omitempty"`
	Attendees   *[]string `json:"attendees,omitempty"`
}
