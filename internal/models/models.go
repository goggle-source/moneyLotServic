package models

import "time"

type HealthCheakDatabase struct {
	ConnDB     bool
	LeadTimeDB time.Duration
}
