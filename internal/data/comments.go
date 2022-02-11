package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Comment struct {
	ID      int    `json:"id"`
	GroupID int    `json:"group_id"`
	Text    string `json:"text"`
}

type Atmosphere struct {
	Negative int `json:"negative,omitempty"`
	Positive int `json:"positive,omitempty"`
	Neutral int `json:"neutral,omitempty"`
}

type Comments struct {
	ID        int64         `json:"id,omitempty"`
	Count     int           `json:"count,omitempty"`
	UserID    int64         `json:"user_id,omitempty"`
	Items     []interface{}	`json:"items,omitempty"`
	ToxicityScale int32 	`json:"toxicity_scale,omitempty"`
	Atmosphere *Atmosphere 	`json:"atmosphere,omitempty"`
	CreatedAt time.Time     `json:"-"`
	Version   int32         `json:"version,omitempty"`
}

type CommentModel struct {
	DB *sql.DB
}

func (m CommentModel) Insert(comments *Comments) error {
	query := `INSERT INTO posts (count, items, user_id)
	VALUES ($1, $2, $3)
	RETURNING id,created_at,version
	`

	items, err := json.Marshal(comments.Items)

	if err != nil {
		return err
	}

	args := []interface{}{comments.Count, items, comments.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&comments.ID, &comments.CreatedAt, &comments.Version)
}
