package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id"`
	ArticleId uuid.UUID `json:"article_id"`
	ThreadId  uuid.UUID `json:"thread_id"`
	ParentId  uuid.UUID `json:"parent_id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Created   time.Time `json:"created"`
	Deleted   bool      `json:"deleted"`
}

type Like struct {
	Id        uuid.UUID `json:"id"`
	CommentId uuid.UUID `json:"comment_id"`
	UserId    uuid.UUID `json:"user_id"`
}

type Dislike struct {
	Id        uuid.UUID `json:"id"`
	CommentId uuid.UUID `json:"comment_id"`
	UserId    uuid.UUID `json:"user_id"`
}

type Complaint struct {
	Id        uuid.UUID `json:"id"`
	CommentId uuid.UUID `json:"comment_id"`
	UserId    uuid.UUID `json:"user_id"`
	Message   string    `json:"message"`
}
