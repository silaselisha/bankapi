package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/silaselisha/bankapi/database/sqlc"
)

type transferRequestParams struct {
	ToAccountId   int64  `json:"to_account_id"`
	FromAccountId int64  `json:"from_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency" binding:"currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequestParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, isValid := s.validateAccount(ctx, req.FromAccountId, req.Currency) 
	if !isValid {
		return
	}

	_, isValid = s.validateAccount(ctx, req.ToAccountId, req.Currency) 
	if !isValid {
		return
	}

	args := db.TransferTxParams{
		Amount:        int32(req.Amount),
		ToAccountId:   req.ToAccountId,
		FromAccountId: req.FromAccountId,
	}

	transfer, err := s.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(s.validateAccount(ctx, req.FromAccountId, req.Currency))
	ctx.JSON(http.StatusCreated, transfer)
}

func (s *Server) validateAccount(ctx *gin.Context, id int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return db.Account{}, false
		}
	
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Account{}, false
	}

	if account.Currency != currency {
		fmt.Println("not equal")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Account{}, false
	}

	return account, true
}