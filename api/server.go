package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/silaselisha/bankapi/database/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	
	router.GET("/accounts", server.getAllAccounts)
	router.POST("/accounts", server.createAccounts)
	router.GET("/accounts/:id", server.getAccountById)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
