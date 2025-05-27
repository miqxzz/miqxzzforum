package entity

type Post struct {
	ID       int    `json:"id" db:"id" example:"1" `
	AuthorId int    `json:"author_id" db:"author_id" example:"1" `
	Title    string `json:"title" db:"title" example:"Заголовк"`
	Content  string `json:"content" db:"content" example:"Текст"`
}
