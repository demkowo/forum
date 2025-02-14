package postgres

import (
	"database/sql"
	"errors"

	model "github.com/demkowo/forum/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	CHECK_IF_EXIST_COMMENTS   = "SELECT to_regclass('public.comments')"
	CHECK_IF_EXIST_LIKES      = "SELECT to_regclass('public.likes')"
	CHECK_IF_EXIST_DISLIKES   = "SELECT to_regclass('public.dislikes')"
	CHECK_IF_EXIST_COMPLAINTS = "SELECT to_regclass('public.complaints')"
	CREATE_TABLE_COMMENTS     = `CREATE TABLE comments (
    id UUID PRIMARY KEY,
    article_id UUID NOT NULL,
	thread_id UUID NOT NULL,
    parent_id UUID,
    author varchar(255) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (parent_id) REFERENCES comments(id),
    FOREIGN KEY (article_id) REFERENCES articles(id),
    FOREIGN KEY (author) REFERENCES users(nickname)
	);`
	CREATE_TABLE_LIKES = `CREATE TABLE likes (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
	);`
	CREATE_TABLE_DISLIKES = `CREATE TABLE dislikes (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
	);`
	CREATE_TABLE_COMPLAINTS = `CREATE TABLE complaints (
    id UUID PRIMARY KEY,
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    message TEXT,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (comment_id, user_id)
	);`
)

type ForumRepo interface {
	CreateTableComments() string
	CreateTableLikes() string
	CreateTableDislikes() string
	CreateTableComplaints() string

	AddComment(comment model.Comment) error
	DeleteComment(commentId uuid.UUID) error
	GetComment(commentId uuid.UUID) (*model.Comment, error)
	FindComments() ([]model.Comment, error)
	FindCommentsByArticle(articleId uuid.UUID) ([]model.Comment, error)
	CountCommentsByArticle(articleId uuid.UUID) (int, error)

	AddLike(like model.Like) error
	DeleteLike(like model.Like) error
	FindLikesByComment(commentId uuid.UUID) ([]model.Like, error)
	CountLikes(commentId uuid.UUID) (int, error)

	AddDislike(dislike model.Dislike) error
	DeleteDislike(model.Dislike) error
	FindDislikesByComment(commentId uuid.UUID) ([]model.Dislike, error)
	CountDislikes(commentId uuid.UUID) (int, error)

	AddComplaint(model.Complaint) error
	DeleteComplaint(uuid.UUID) error
	FindComplaintsByComment(commentId uuid.UUID) ([]model.Complaint, error)
	CountComplaints(commentId uuid.UUID) (int, error)
}

type forumRepo struct {
	db *sql.DB
}

func NewForum(db *sql.DB) ForumRepo {
	log.Trace()

	return &forumRepo{
		db: db,
	}
}

func (r *forumRepo) CreateTableComments() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_EXIST_COMMENTS)
	if err != nil {
		log.Panicf("CHECK_IF_EXIST_COMMENTS failed: %v", err)
	}
	defer rows.Close()

	var tableName sql.NullString
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Panicf("CHECK_IF_EXIST_COMMENTS rows scan failed: %v", err)
		}
	}

	if tableName.Valid {
		return "DB comments ready to go"
	}

	_, err = r.db.Exec(CREATE_TABLE_COMMENTS)
	if err != nil {
		log.Panicf("CREATE_TABLE_COMMENTS failed: %v", err)
	}

	return "Table comments created, DB ready to go"

}

func (r *forumRepo) CreateTableLikes() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_EXIST_LIKES)
	if err != nil {
		log.Panicf("CHECK_IF_EXIST_LIKES failed: %v", err)
	}
	defer rows.Close()

	var tableName sql.NullString
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Panicf("CHECK_IF_EXIST_LIKES rows scan failed: %v", err)
		}
	}

	if tableName.Valid {
		return "DB likes ready to go"
	}

	_, err = r.db.Exec(CREATE_TABLE_LIKES)
	if err != nil {
		log.Panicf("CREATE_TABLE_LIKES failed: %v", err)
	}

	return "Table likes created, DB ready to go"

}

func (r *forumRepo) CreateTableDislikes() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_EXIST_DISLIKES)
	if err != nil {
		log.Panicf("CHECK_IF_EXIST_DISLIKES failed: %v", err)
	}
	defer rows.Close()

	var tableName sql.NullString
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Panicf("CHECK_IF_EXIST_DISLIKES rows scan failed: %v", err)
		}
	}

	if tableName.Valid {
		return "DB dislikes ready to go"
	}

	_, err = r.db.Exec(CREATE_TABLE_DISLIKES)
	if err != nil {
		log.Panicf("CREATE_TABLE_DISLIKES failed: %v", err)
	}

	return "Table dislikes created, DB ready to go"

}

func (r *forumRepo) CreateTableComplaints() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_EXIST_COMPLAINTS)
	if err != nil {
		log.Panicf("CHECK_IF_EXIST_COMPLAINTS failed: %v", err)
	}
	defer rows.Close()

	var tableName sql.NullString
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Panicf("CHECK_IF_EXIST_COMPLAINTS rows scan failed: %v", err)
		}
	}

	if tableName.Valid {
		return "DB complaints ready to go"
	}

	_, err = r.db.Exec(CREATE_TABLE_COMPLAINTS)
	if err != nil {
		log.Panicf("CREATE_TABLE_COMPLAINTS failed: %v", err)
	}

	return "Table complaints created, DB ready to go"

}

func (r *forumRepo) AddComment(comment model.Comment) error {
	log.Trace()

	COMMENTS_ADD := "INSERT INTO comments (id, article_id, thread_id, parent_id, author, content, created, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	var parent interface{}

	if comment.ParentId.String() == "00000000-0000-0000-0000-000000000000" {
		parent = nil
	} else {
		parent = comment.ParentId
	}

	_, err := r.db.Exec(COMMENTS_ADD, comment.Id, comment.ArticleId, comment.ThreadId, parent, comment.Author, comment.Content, comment.Created, comment.Deleted)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *forumRepo) DeleteComment(commentId uuid.UUID) error {
	log.Trace()

	query := `UPDATE comments SET deleted = TRUE WHERE id = $1`

	_, err := r.db.Exec(query, commentId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *forumRepo) GetComment(commentId uuid.UUID) (*model.Comment, error) {
	log.Trace()

	query := `
        SELECT id, article_id, thread_id, parent_id, author, content, created, deleted
        FROM comments
        WHERE id = $1
    `
	row := r.db.QueryRow(query, commentId)
	var comment model.Comment
	err := row.Scan(&comment.Id, &comment.ArticleId, &comment.ThreadId, &comment.ParentId, &comment.Author, &comment.Content, &comment.Created, &comment.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(err)
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}
	return &comment, nil
}

func (r *forumRepo) FindComments() ([]model.Comment, error) {
	log.Trace()

	query := `
        SELECT id, article_id, thread_id, parent_id, author, content, created, deleted
        FROM comments
        WHERE deleted = FALSE
		ORDER by created DESC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.ArticleId, &comment.ThreadId, &comment.ParentId, &comment.Author, &comment.Content, &comment.Created, &comment.Deleted)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *forumRepo) FindCommentsByArticle(articleId uuid.UUID) ([]model.Comment, error) {
	log.Trace()

	query := `
        SELECT id, article_id, thread_id, parent_id, author, content, created, deleted
        FROM comments
        WHERE article_id = $1 AND deleted = FALSE
		ORDER by created DESC
    `
	rows, err := r.db.Query(query, articleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.ArticleId, &comment.ThreadId, &comment.ParentId, &comment.Author, &comment.Content, &comment.Created, &comment.Deleted)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *forumRepo) CountCommentsByArticle(articleId uuid.UUID) (int, error) {
	log.Trace()

	query := `
        SELECT COUNT(*)
        FROM comments
        WHERE article_id = $1
    `
	var count int
	err := r.db.QueryRow(query, articleId).Scan(&count)
	if err != nil {
		log.Warn("db.QueryRow failed: ", err)
		return 0, nil
	}

	return count, nil

}

func (r *forumRepo) AddLike(like model.Like) error {
	log.Trace()

	query := `
        INSERT INTO likes (id, comment_id, user_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (comment_id, user_id) DO NOTHING
    `

	_, err := r.db.Exec(query, like.Id, like.CommentId, like.UserId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *forumRepo) DeleteLike(like model.Like) error {
	log.Trace()

	query := `
        DELETE FROM likes
        WHERE comment_id = $1 AND user_id = $2
    `
	result, err := r.db.Exec(query, like.CommentId, like.UserId)
	if err != nil {
		log.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *forumRepo) FindLikesByComment(commentId uuid.UUID) ([]model.Like, error) {
	log.Trace()

	query := `
        SELECT id, comment_id, user_id
        FROM likes
        WHERE comment_id = $1
    `
	rows, err := r.db.Query(query, commentId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var likes []model.Like
	for rows.Next() {
		var like model.Like
		err := rows.Scan(&like.Id, &like.CommentId, &like.UserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (r *forumRepo) CountLikes(commentId uuid.UUID) (int, error) {
	log.Trace()

	query := `
        SELECT COUNT(*)
        FROM likes
        WHERE comment_id = $1
    `
	var count int
	err := r.db.QueryRow(query, commentId).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, nil
	}

	return count, nil
}

func (r *forumRepo) AddDislike(dislike model.Dislike) error {
	log.Trace()

	query := `
        INSERT INTO dislikes (id, comment_id, user_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (comment_id, user_id) DO NOTHING
    `
	_, err := r.db.Exec(query, dislike.Id, dislike.CommentId, dislike.UserId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *forumRepo) DeleteDislike(dislike model.Dislike) error {
	log.Trace()

	query := `
        DELETE FROM dislikes
        WHERE comment_id = $1 AND user_id = $2
    `
	result, err := r.db.Exec(query, dislike.CommentId, dislike.UserId)
	if err != nil {
		log.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}

	if rowsAffected == 0 {
		log.Warn("dislike not found")
		return errors.New("dislike not found")
	}

	return nil
}

func (r *forumRepo) FindDislikesByComment(commentId uuid.UUID) ([]model.Dislike, error) {
	log.Trace()

	query := `
        SELECT id, comment_id, user_id
        FROM dislikes
        WHERE comment_id = $1
    `
	rows, err := r.db.Query(query, commentId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var dislikes []model.Dislike
	for rows.Next() {
		var dislike model.Dislike
		err := rows.Scan(&dislike.Id, &dislike.CommentId, &dislike.UserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		dislikes = append(dislikes, dislike)
	}

	return dislikes, nil
}

func (r *forumRepo) CountDislikes(commentId uuid.UUID) (int, error) {
	log.Trace()

	query := `
        SELECT COUNT(*)
        FROM dislikes
        WHERE comment_id = $1
    `
	var count int
	err := r.db.QueryRow(query, commentId).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return count, nil
}

func (r *forumRepo) AddComplaint(complaint model.Complaint) error {
	log.Trace()

	query := `
        INSERT INTO complaints (id, comment_id, user_id, message)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (comment_id, user_id)
		DO UPDATE SET message = complaints.message || E'\n' || EXCLUDED.message
    `
	_, err := r.db.Exec(query, complaint.Id, complaint.CommentId, complaint.UserId, complaint.Message)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *forumRepo) DeleteComplaint(id uuid.UUID) error {
	log.Trace()

	query := `DELETE FROM complaints WHERE id = $1
    `
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}

	if rowsAffected == 0 {
		log.Error("complaint not found")
		return errors.New("complaint not found")
	}

	return nil
}

func (r *forumRepo) FindComplaintsByComment(commentId uuid.UUID) ([]model.Complaint, error) {
	log.Trace()

	query := `
        SELECT id, comment_id, user_id, message
        FROM complaints
        WHERE comment_id = $1
    `
	rows, err := r.db.Query(query, commentId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var complaints []model.Complaint
	for rows.Next() {
		var complaint model.Complaint
		err := rows.Scan(&complaint.Id, &complaint.CommentId, &complaint.UserId, &complaint.Message)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

func (r *forumRepo) CountComplaints(commentId uuid.UUID) (int, error) {
	log.Trace()

	query := `
        SELECT COUNT(*)
        FROM complaints
        WHERE comment_id = $1
    `
	var count int
	err := r.db.QueryRow(query, commentId).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return count, nil
}
