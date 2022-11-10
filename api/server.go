package api

import (
	"shortURL/db/redis"
	db "shortURL/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Querier
	redis  *redis.RedisQueries
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Querier, redis *redis.RedisQueries) *Server {
	server := &Server{
		store: store,
		redis: redis,
	}

	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/short", server.createShortURL)  // 建立短連結
	router.GET("/:short_url", server.getRedirect) // 導向長連結

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
