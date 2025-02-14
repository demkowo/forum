package handler

import (
	"fmt"
	"net/http"
	"time"

	model "github.com/demkowo/forum/models"
	service "github.com/demkowo/forum/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Forum interface {
	CreateTableComments()
	CreateTableLikes()
	CreateTableDislikes()
	CreateTableComplaints()

	AddComment(c *gin.Context)
	DeleteComment(c *gin.Context)
	GetComment(c *gin.Context)
	FindComments(c *gin.Context)
	FindCommentsByArticle(c *gin.Context)
	CountComments(c *gin.Context)

	AddLike(c *gin.Context)
	DeleteLike(c *gin.Context)
	FindLikesByComment(c *gin.Context)
	CountLikes(c *gin.Context)

	AddDislike(c *gin.Context)
	DeleteDislike(c *gin.Context)
	FindDislikesByComment(c *gin.Context)
	CountDislikes(c *gin.Context)

	AddComplaint(c *gin.Context)
	DeleteComplaint(c *gin.Context)
	FindComplaintsByComment(c *gin.Context)
	CountComplaints(c *gin.Context)
}

type CommentNode struct {
	ID      uuid.UUID      `json:"id"`
	Comment *model.Comment `json:"comment"`
	Childs  []CommentNode  `json:"childs"`
}
type forum struct {
	service service.Forum
}

func NewForum(service service.Forum) Forum {
	log.Trace()

	return &forum{
		service: service,
	}
}

func (h *forum) CreateTableComments() {
	log.Trace()

	log.Info(h.service.CreateTableComments())
}

func (h *forum) CreateTableLikes() {
	log.Trace()

	log.Info(h.service.CreateTableLikes())
}

func (h *forum) CreateTableDislikes() {
	log.Trace()

	log.Info(h.service.CreateTableDislikes())
}

func (h *forum) CreateTableComplaints() {
	log.Trace()

	log.Info(h.service.CreateTableComplaints())
}

func (h *forum) AddComment(c *gin.Context) {
	log.Trace()

	var input struct {
		ArticleID string `json:"article_id" binding:"required"`
		ThreadID  string `json:"thread_id"`
		ParentID  string `json:"parent_id"`
		Author    string `json:"author" binding:"required"`
		Content   string `json:"content" binding:"required"`
		ReplyTo   string `json:"reply_to"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	articleId, err := uuid.Parse(input.ArticleID)
	if err != nil {
		log.Errorf("Invalid article_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article_id format"})
		return
	}

	id := uuid.New()

	threadId := id
	if input.ThreadID != "" {
		threadId, err = uuid.Parse(input.ThreadID)
		if err != nil {
			log.Errorf("Invalid thread_id UUID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread_id format"})
			return
		}
	}

	var parentId uuid.UUID
	if input.ParentID != "" {
		parentId, err = uuid.Parse(input.ParentID)
		if err != nil {
			log.Errorf("Invalid parent_id UUID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent_id format"})
			return
		}
	}

	if input.ReplyTo != "" {
		input.Content = fmt.Sprintf("@%s: %s", input.ReplyTo, input.Content)
	}

	comment := &model.Comment{
		Id:        id,
		ArticleId: articleId,
		ThreadId:  threadId,
		ParentId:  parentId,
		Author:    input.Author,
		Content:   input.Content,
		Created:   time.Now(),
		Deleted:   false,
	}

	if err := h.service.AddComment(*comment); err != nil {
		log.Errorf("Failed to add comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	newNode := CommentNode{
		ID:      comment.Id,
		Comment: comment,
		Childs:  []CommentNode{},
	}

	c.JSON(http.StatusOK, gin.H{
		"comment_added": newNode,
	})
}

func (h *forum) DeleteComment(c *gin.Context) {
	log.Trace()

	idStr := c.Param("comment_id")
	commentId, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := h.service.DeleteComment(commentId); err != nil {
		log.Errorf("Failed to delete comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (h *forum) GetComment(c *gin.Context) {
	log.Trace()

	idStr := c.Param("comment_id")
	commentId, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.service.GetComment(commentId)
	if err != nil {
		log.Errorf("Failed to retrieve comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": comment})
}

func (h *forum) FindComments(c *gin.Context) {
	log.Trace()

	res, err := h.service.FindComments()
	if err != nil {
		log.Errorf("Failed to retrieve comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	commentMap := make(map[uuid.UUID]*CommentNode)

	for _, val := range res {
		node := &CommentNode{
			ID:      val.Id,
			Comment: &val,
			Childs:  []CommentNode{},
		}
		commentMap[val.Id] = node
	}

	var roots []CommentNode

	for _, val := range res {
		if val.ThreadId == val.Id {
			roots = append(roots, *commentMap[val.Id])
		} else {
			parentNode, found := commentMap[val.ThreadId]
			if found {
				parentNode.Childs = append(parentNode.Childs, *commentMap[val.Id])
			} else {
				log.Warnf("Parent comment with ID %s not found for child %s", val.ThreadId, val.Id)
				roots = append(roots, *commentMap[val.Id])
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": roots,
		"count":    len(res),
	})
}

func (h *forum) FindCommentsByArticle(c *gin.Context) {
	log.Trace()

	articleIdStr := c.Param("article_id")
	articleId, err := uuid.Parse(articleIdStr)
	if err != nil {
		log.Errorf("Invalid article ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	res, err := h.service.FindCommentsByArticle(articleId)
	if err != nil {
		log.Errorf("Failed to retrieve comments for article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	commentMap := make(map[uuid.UUID]*CommentNode)

	for _, val := range res {
		node := &CommentNode{
			ID:      val.Id,
			Comment: &val,
			Childs:  []CommentNode{},
		}
		commentMap[val.Id] = node
	}

	var roots []CommentNode

	for _, val := range res {
		if val.ThreadId == val.Id {
			roots = append(roots, *commentMap[val.Id])
		} else {
			parentNode, found := commentMap[val.ThreadId]
			if found {
				parentNode.Childs = append(parentNode.Childs, *commentMap[val.Id])
			} else {
				log.Warnf("Parent comment with ID %s not found for child %s", val.ThreadId, val.Id)
				roots = append(roots, *commentMap[val.Id])
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": roots,
		"count":    len(res),
	})
}

func (h *forum) CountComments(c *gin.Context) {
	log.Trace()

	idStr := c.Param("article_id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	nr, err := h.service.CountCommentsByArticle(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "counting number of articles failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments_amount": nr})
}

func (h *forum) AddLike(c *gin.Context) {
	log.Trace()

	var input struct {
		CommentID string `json:"comment_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	commentId, err := uuid.Parse(input.CommentID)
	if err != nil {
		log.Errorf("Invalid comment_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id format"})
		return
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Errorf("Invalid user_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	like := model.Like{
		Id:        uuid.New(),
		CommentId: commentId,
		UserId:    userId,
	}

	if err := h.service.AddLike(like); err != nil {
		log.Errorf("Failed to add like: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add like",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like added successfully"})
}

func (h *forum) DeleteLike(c *gin.Context) {
	log.Trace()

	var input struct {
		CommentID string `json:"comment_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	commentId, err := uuid.Parse(input.CommentID)
	if err != nil {
		log.Errorf("Invalid comment_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id format"})
		return
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Errorf("Invalid user_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	like := model.Like{
		CommentId: commentId,
		UserId:    userId,
	}

	if err := h.service.DeleteLike(like); err != nil {
		log.Errorf("Failed to remove like: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove like",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}

func (h *forum) FindLikesByComment(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	likes, err := h.service.FindLikesByComment(commentId)
	if err != nil {
		log.Errorf("Failed to retrieve likes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}

func (h *forum) CountLikes(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	count, err := h.service.CountLikes(commentId)
	if err != nil {
		log.Errorf("Failed to count likes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count likes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"number_of_likes": count})
}

func (h *forum) AddDislike(c *gin.Context) {
	log.Trace()

	var input struct {
		CommentID string `json:"comment_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	commentId, err := uuid.Parse(input.CommentID)
	if err != nil {
		log.Errorf("Invalid comment_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id format"})
		return
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Errorf("Invalid user_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	dislike := model.Dislike{
		Id:        uuid.New(),
		CommentId: commentId,
		UserId:    userId,
	}

	if err := h.service.AddDislike(dislike); err != nil {
		log.Errorf("Failed to add dislike: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add dislike",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dislike added successfully",
	})
}

func (h *forum) DeleteDislike(c *gin.Context) {
	log.Trace()

	var input struct {
		CommentID string `json:"comment_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	commentId, err := uuid.Parse(input.CommentID)
	if err != nil {
		log.Errorf("Invalid comment_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id format"})
		return
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Errorf("Invalid user_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	dislike := model.Dislike{
		CommentId: commentId,
		UserId:    userId,
	}

	if err := h.service.DeleteDislike(dislike); err != nil {
		log.Errorf("Failed to remove dislike: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove dislike",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dislike removed successfully"})
}

func (h *forum) FindDislikesByComment(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	dislikes, err := h.service.FindDislikesByComment(commentId)
	if err != nil {
		log.Errorf("Failed to retrieve dislikes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dislikes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dislikes": dislikes})
}

func (h *forum) CountDislikes(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	count, err := h.service.CountDislikes(commentId)
	if err != nil {
		log.Errorf("Failed to count dislikes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count dislikes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"number_of_dislikes": count})
}

func (h *forum) AddComplaint(c *gin.Context) {
	log.Trace()

	var input struct {
		CommentID string `json:"comment_id" binding:"required"`
		UserID    string `json:"user_id" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	commentId, err := uuid.Parse(input.CommentID)
	if err != nil {
		log.Errorf("Invalid comment_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment_id format"})
		return
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		log.Errorf("Invalid user_id UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	complaint := model.Complaint{
		Id:        uuid.New(),
		CommentId: commentId,
		UserId:    userId,
		Message:   input.Message,
	}

	if err := h.service.AddComplaint(complaint); err != nil {
		log.Errorf("Failed to add complaint: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add complaint",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Complaint added successfully"})
}

func (h *forum) DeleteComplaint(c *gin.Context) {
	log.Trace()

	idStr := c.Param("complaint_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("Invalid complaint ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}

	if err := h.service.DeleteComplaint(id); err != nil {
		log.Errorf("Failed to remove complaint: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove complaint",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Complaint removed successfully"})
}

func (h *forum) FindComplaintsByComment(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	complaints, err := h.service.FindComplaintsByComment(commentId)
	if err != nil {
		log.Errorf("Failed to retrieve complaints: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve complaints"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complaints": complaints})
}

func (h *forum) CountComplaints(c *gin.Context) {
	log.Trace()

	commentIdStr := c.Param("comment_id")

	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		log.Errorf("Invalid comment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	count, err := h.service.CountComplaints(commentId)
	if err != nil {
		log.Errorf("Failed to count complaints: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count complaints"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"number_of_complaints": count})
}
