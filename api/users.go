package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/silaselisha/bankapi/token"
)

type createUserRequestParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userResponse struct {
	Username  string    `db:"username"`
	Fullname  string    `db:"fullname"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type createUserResponseParams struct {
	Status int          `json:"status"`
	Token  string       `json:"token"`
	User   userResponse `json:"user"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Username: req.Username,
		Fullname: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	envs, err := utils.Load("../..")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	maker, err := token.NewJwtMaker(envs.JwtSecreteKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := maker.CreateToken(user.Username, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponseParams{
		Status: http.StatusCreated,
		Token:  token,
		User: userResponse{
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
	ctx.JSON(http.StatusCreated, response)
}

type loginRequestParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponseParams struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	envs, err := utils.Load("../..")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	maker, err := token.NewJwtMaker(envs.JwtSecreteKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := maker.CreateToken(user.Username, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginResponseParams{
		Status: http.StatusOK,
		Token:  token,
	}
	ctx.JSON(http.StatusOK, response)
}
