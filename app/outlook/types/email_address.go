package types

// EmailAddress represents the Microsoft Graph API's emailAddress resource type.
type EmailAddress struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}
