package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.CurrencyValidator)
	}

	router.GET("/accounts", server.getAllAccounts)
	router.POST("/accounts", server.createAccounts)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)
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
