package app

import (
	handler "github.com/demkowo/forum/handlers"
	log "github.com/sirupsen/logrus"
)

func addForumRoutes(h handler.Forum) {
	log.Trace()

	public := router.Group("/api/v1/")
	auth := router.Group("/api/v1/")

	auth.POST("/comments/add", h.AddComment)
	auth.DELETE("/comments/delete/:comment_id", h.DeleteComment)
	auth.GET("/comments/get/:comment_id", h.GetComment)
	auth.GET("/comments/find", h.FindComments)
	public.GET("/comments/find/:article_id", h.FindCommentsByArticle)
	public.GET("/comments/count/:article_id", h.CountComments)

	auth.POST("/likes/add", h.AddLike)
	auth.DELETE("/likes/delete", h.DeleteLike)
	public.GET("/likes/count/:comment_id", h.CountLikes)
	public.GET("/likes/find/:comment_id", h.FindLikesByComment)

	auth.POST("/dislikes/add", h.AddDislike)
	auth.DELETE("dislikes/delete", h.DeleteDislike)
	public.GET("/dislikes/count/:comment_id", h.CountDislikes)
	public.GET("/dislikes/find/:comment_id", h.FindDislikesByComment)

	auth.POST("/complaints/add", h.AddComplaint)
	auth.DELETE("/complaints/delete/:complaint_id", h.DeleteComplaint)
	public.GET("/complaints/count/:comment_id", h.CountComplaints)
	public.GET("/complaints/find/:comment_id", h.FindComplaintsByComment)
}
