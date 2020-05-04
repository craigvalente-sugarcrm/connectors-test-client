package common

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

// Action - represent a single calendar action (CREATE|UPDATE|DELETE)
type Action struct {
	Account Account       `json:"account,omitempty"`
	Items   []*Projection `json:"items,omitempty"`
}

// GenerateProjection - creates a new projection for test
func (a *Action) GenerateProjection() *Projection {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1440 * 90) // 90 days
	t := time.Now().UTC().Add(time.Minute * time.Duration(r))

	return &Projection{
		CaseID:   uuid.New().String(),
		Type:     "create",
		Owner:    a.Account.Email,
		Subject:  fmt.Sprintf("%s: %v", TestEventPrefix, t.Unix()),
		Start:    t,
		End:      t.Add(time.Minute * time.Duration(30)),
		Location: randomdata.City(),
	}
}

// UpdateProjection - updates an existing projection
func (a *Action) UpdateProjection(p *Projection) *Projection {
	return &Projection{
		CaseID:   p.CaseID,
		Type:     "update",
		Owner:    a.Account.Email,
		Subject:  p.Subject,
		Start:    p.Start,
		End:      p.End.Add(time.Minute * time.Duration(30)),
		Location: p.Location,
	}
}
