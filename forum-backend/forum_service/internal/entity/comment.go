package entity

import "time"

type Comment struct {
	ID        int       `json:"id" db:"id" exmaple:"1"`
	AuthorId  int       `json:"author_id" db:"author_id" exmaple:"1"`
	PostId    int       `json:"post_id" db:"post_id" exmaple:"1"`
	Content   string    `json:"content" db:"content" exmaple:"текст комментария"`
	CreatedAt time.Time `json:"created_at" exmaple:"22:00"`
}

func (c *Comment) Validate() error {
	if c.AuthorId <= 0 {
		return ErrInvalidAuthorID
	}
	if c.PostId <= 0 {
		return ErrInvalidPostID
	}
	if c.Content == "" {
		return ErrEmptyContent
	}
	return nil
}
