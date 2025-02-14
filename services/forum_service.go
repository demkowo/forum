package service

import (
	model "github.com/demkowo/forum/models"
	"github.com/demkowo/forum/repositories/postgres"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Forum interface {
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

type forum struct {
	repo postgres.ForumRepo
}

func NewForum(repository postgres.ForumRepo) Forum {
	return &forum{
		repo: repository,
	}
}

func (s *forum) CreateTableComments() string {
	log.Trace()

	return s.repo.CreateTableComments()
}

func (s *forum) CreateTableLikes() string {
	log.Trace()

	return s.repo.CreateTableLikes()
}

func (s *forum) CreateTableDislikes() string {
	log.Trace()

	return s.repo.CreateTableDislikes()
}

func (s *forum) CreateTableComplaints() string {
	log.Trace()

	return s.repo.CreateTableComplaints()
}

func (s *forum) AddComment(comment model.Comment) error {
	log.Trace()

	return s.repo.AddComment(comment)
}

func (s *forum) DeleteComment(commentId uuid.UUID) error {
	log.Trace()
	return s.repo.DeleteComment(commentId)
}

func (s *forum) GetComment(commentId uuid.UUID) (*model.Comment, error) {
	log.Trace()
	return s.repo.GetComment(commentId)
}

func (s *forum) FindComments() ([]model.Comment, error) {
	log.Trace()
	return s.repo.FindComments()
}

func (s *forum) FindCommentsByArticle(articleId uuid.UUID) ([]model.Comment, error) {
	log.Trace()
	return s.repo.FindCommentsByArticle(articleId)
}

func (s *forum) CountCommentsByArticle(articleId uuid.UUID) (int, error) {
	log.Trace()
	return s.repo.CountCommentsByArticle(articleId)
}

func (s *forum) AddLike(like model.Like) error {
	log.Trace()

	dislike := model.Dislike{
		CommentId: like.CommentId,
		UserId:    like.UserId,
	}

	if err := s.repo.DeleteDislike(dislike); err != nil {
		if err.Error() != "dislike not found" {
			return err
		}
	}

	return s.repo.AddLike(like)
}

func (s *forum) DeleteLike(like model.Like) error {
	log.Trace()
	return s.repo.DeleteLike(like)
}

func (s *forum) FindLikesByComment(commentId uuid.UUID) ([]model.Like, error) {
	log.Trace()
	return s.repo.FindLikesByComment(commentId)
}

func (f *forum) CountLikes(commentId uuid.UUID) (int, error) {
	log.Trace()
	return f.repo.CountLikes(commentId)
}

func (s *forum) AddDislike(dislike model.Dislike) error {
	log.Trace()

	like := model.Like{
		CommentId: dislike.CommentId,
		UserId:    dislike.UserId,
	}

	if err := s.repo.DeleteLike(like); err != nil {
		if err.Error() != "like not found" {
			return err
		}
	}

	return s.repo.AddDislike(dislike)
}

func (s *forum) DeleteDislike(dislike model.Dislike) error {
	log.Trace()
	return s.repo.DeleteDislike(dislike)
}

func (s *forum) FindDislikesByComment(commentId uuid.UUID) ([]model.Dislike, error) {
	log.Trace()
	return s.repo.FindDislikesByComment(commentId)
}

func (s *forum) CountDislikes(commentId uuid.UUID) (int, error) {
	log.Trace()
	return s.repo.CountDislikes(commentId)
}

func (s *forum) AddComplaint(complaint model.Complaint) error {
	log.Trace()
	return s.repo.AddComplaint(complaint)
}

func (s *forum) DeleteComplaint(id uuid.UUID) error {
	log.Trace()
	return s.repo.DeleteComplaint(id)
}

func (s *forum) FindComplaintsByComment(commentId uuid.UUID) ([]model.Complaint, error) {
	log.Trace()
	return s.repo.FindComplaintsByComment(commentId)
}

func (s *forum) CountComplaints(commentId uuid.UUID) (int, error) {
	log.Trace()
	return s.repo.CountComplaints(commentId)
}
