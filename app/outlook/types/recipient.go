package types

// Recipient represents the Microsoft Graph API's recipient resource type.
type Recipient struct {
	EmailAddress *EmailAddress `json:"emailAddress,omitempty"`
}
