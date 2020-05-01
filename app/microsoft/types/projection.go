package types

import "time"

// Projection - holds calendar event data
type Projection struct {
	Deleted     bool      `json:"deleted,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	Description string    `json:"description,omitempty"`
	Owner       string    `json:"owner,omitempty"`
	Start       time.Time `json:"start,omitempty"`
	End         time.Time `json:"end,omitempty"`
	Location    string    `json:"location,omitempty"`
}
