package app

import (
	"database/sql"
	"os"

	handler "github.com/demkowo/forum/handlers"
	postgres "github.com/demkowo/forum/repositories/postgres"
	service "github.com/demkowo/forum/services"
	logger "github.com/demkowo/forum/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	portNumber = ":5000"
)

var (
	dbConnection = os.Getenv("DB_CONNECTION")
	router       = gin.Default()
)

func init() {
	logger.Start.BasicConfig()
}

func Start() {
	log.Trace()

	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	forumRepo := postgres.NewForum(db)
	forumService := service.NewForum(forumRepo)
	forumHandler := handler.NewForum(forumService)
	addForumRoutes(forumHandler)

	forumHandler.CreateTableComments()
	forumHandler.CreateTableComplaints()
	forumHandler.CreateTableLikes()
	forumHandler.CreateTableDislikes()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(portNumber)
}
