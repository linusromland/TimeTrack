package dtos

import "time"

type TimePeriod struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
}

type CreateTimeEntryInput struct {
	ProjectID string     `json:"project_id" binding:"required,uuid"`
	Period    TimePeriod `json:"period" binding:"required"`
	Note      string     `json:"note" binding:"omitempty,max=1024"`
}

type UpdateTimeEntryInput struct {
	ProjectID *string     `json:"project_id" binding:"omitempty,uuid"`
	Period    *TimePeriod `json:"period" binding:"omitempty"`
	Note      *string     `json:"note" binding:"omitempty,max=1024"`
}
