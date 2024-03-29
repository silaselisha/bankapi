package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/token"
)

type createAccountsParams struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccounts(ctx *gin.Context) {
	var req createAccountsParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.Owner != authPayload.Username {
		err := errors.New("unauthorized user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

type accountIdParams struct {
	Id int64 `uri:"id" binding:"required"`
}

func (s *Server) getAccountById(ctx *gin.Context) {
	var req accountIdParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsParams struct {
	Owner    string `form:"owner" binding:"required"`
	PageId   int64  `form:"page_id" binding:"required,min=1"`
	PageSize int64  `form:"page_size" binding:"required,min=3,max=10"`
}

func (s *Server) getAllAccounts(ctx *gin.Context) {
	var reqQuery listAccountsParams
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)
	if authPayload.Username != reqQuery.Owner {
		err := errors.New("unauthorized request")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  int32(reqQuery.PageSize),
		Offset: (int32(reqQuery.PageId) - 1) * int32(reqQuery.PageSize),
	}

	accounts, err := s.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
