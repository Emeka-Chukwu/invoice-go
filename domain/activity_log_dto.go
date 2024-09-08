package domain

import "time"

type Activity struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateActivityDTO struct {
	UserID     int    `json:"user_id"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
}

type IDParamPayload struct {
	ID int `uri:"id" binding:"required"`
}
