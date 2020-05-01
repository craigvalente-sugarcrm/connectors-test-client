package common

// Action - represent a single calendar action (CREATE|UPDATE|DELETE)
type Action struct {
	Account  string        `json:"account,omitempty"`
	Calendar string        `json:"calendar,omitempty"`
	items    *[]Projection `json:"projection,omitempty"`
}
