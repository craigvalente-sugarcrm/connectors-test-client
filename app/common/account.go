package common

// Account - user account used for test
type Account struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Calendar string `json:"calendar,omitempty"`
}
