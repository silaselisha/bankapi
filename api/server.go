package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/silaselisha/bank-api/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}


func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.getAccounts)

	server.router = router
	return server
}

func errorResponse(err error) *gin.H {
	return &gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}