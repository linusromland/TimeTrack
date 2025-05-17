package dtos

import "time"

type CreateTimeEntryInput struct {
	ProjectID string `json:"project_id" binding:"required,uuid"`
	Period    struct {
		Start time.Time `json:"start" binding:"required"`
		End   time.Time `json:"end" binding:"required"`
	} `json:"period" binding:"required"`
	Note string `json:"note" binding:"omitempty,max=1024"`
}

type UpdateTimeEntryInput struct {
	ProjectID *string `json:"project_id" binding:"omitempty,uuid"`
	Period    *struct {
		Start time.Time `json:"start" binding:"required"`
		End   time.Time `json:"end" binding:"required"`
	} `json:"period" binding:"omitempty"`
	Note *string `json:"note" binding:"omitempty,max=1024"`
}
