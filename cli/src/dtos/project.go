package dtos

type IntegrationInfo struct {
	Type       string `bson:"type" json:"type" binding:"required,oneof=jira"`
	Key        string `bson:"key" json:"key"`
	ExternalID string `bson:"external_id" json:"external_id"`
}

type CreateProjectInput struct {
	Name        string          `json:"name" binding:"required,min=1"`
	Integration IntegrationInfo `bson:"integration" json:"integration" binding:"omitempty"`
}

type UpdateProjectInput struct {
	Name        *string          `json:"name" binding:"omitempty,min=1"`
	Integration *IntegrationInfo `bson:"integration" json:"integration" binding:"omitempty"`
}
