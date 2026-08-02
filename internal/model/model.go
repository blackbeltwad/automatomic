// Set up structs that will be used throughout

package model

import (
	"time"
)

type User struct {

	ID string `json:"id"` 
	GitHubID int64 `json:"github_id"` 
	Email string `json:"email"`
	Username string `json:"username"`
	Role string `json:"role"`
	AvatarURL string `json:"avatar_url"` // profile picture url
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`


}

type Claims struct {
	UserID string   `json:"user_id"`
	Role   string   `json:"role"`
	Scopes []string `json:"scopes"`
}

const (
	ScopePipelineRead  = "pipeline:read"
	ScopePipelineWrite = "pipeline:write"
	ScopeAdmin         = "admin"
)