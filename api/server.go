package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/silaselisha/bankapi/token"
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

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	envs, err := utils.Load("..")
	if err != nil {
		log.Panic(err)
		return nil
	}

	maker, err := token.NewJwtMaker(envs.JwtSecreteKey)
	if err != nil {
		log.Panic(err)
		return nil
	}

	authRouter := router.Group("/").Use(AuthorizationMiddleware(maker))

	authRouter.GET("/accounts", server.getAllAccounts)
	authRouter.POST("/accounts", server.createAccounts)
	authRouter.POST("/transfers", server.createTransfer)
	authRouter.GET("/accounts/:id", server.getAccountById)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
