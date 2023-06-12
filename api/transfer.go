package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/silaselisha/bank-api/db/sqlc"
)

type createTransferTxParams struct {
	FromAccountId int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountId int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferTxParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID: req.ToAccountId,
		Amount: req.Amount,
	}

	if !server.validate(ctx, args.FromAccountID, req.Currency) {
		err := fmt.Errorf("invalid currency %s", req.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validate(ctx, args.ToAccountID, req.Currency) {
		err := fmt.Errorf("invalid currency %s", req.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}

func (server *Server) validate(ctx *gin.Context, id int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		return false
	}
	return true
}