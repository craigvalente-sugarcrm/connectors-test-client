package types

// Calendar represents the Microsoft Graph API's Calendar resource.
type Calendar struct {
	CanEdit             bool          `json:"canEdit,omitempty"`
	CanShare            bool          `json:"canShare,omitempty"`
	CanViewPrivateItems bool          `json:"canViewPrivateItems,omitempty"`
	ChangeKey           string        `json:"changeKey,omitempty"`
	Events              []*Event      `json:"events,omitempty"`
	ID                  string        `json:"id,omitempty"`
	Name                string        `json:"name,omitempty"`
	Owner               *EmailAddress `json:"owner,omitempty"`
}

// Calendars represents a list of calendars from the Microsoft Graph API's
// Calendar resource.
type Calendars struct {
	Calendars []*Calendar `json:"value,omitempty"`
}
