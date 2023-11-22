package token

import (
	"fmt"
	"testing"
	"time"

	db "github.com/silaselisha/bankapi/db/sqlc"
	"github.com/silaselisha/bankapi/db/utils"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	firstName, err := utils.RandomString(6)
	require.NoError(t, err)
	lastName, err := utils.RandomString(6)
	require.NoError(t, err)

	pass, err := utils.RandomString(8)
	require.NoError(t, err)
	hashedPassword, err := utils.GenerateHashedPassword(pass)
	require.NoError(t, err)

	user := db.User{
		Username:  firstName,
		Fullname:  fmt.Sprintf("%s %s", firstName, lastName),
		Email:     fmt.Sprintf("%s%d@outlook.com", lastName, utils.RandomAmount(1, 100)),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC),
	}

	key, err := utils.RandomString(32)
	require.NoError(t, err)
	maker, err := NewJwtMaker(key)
	require.NoError(t, err)

	token, err := maker.CreateToken(user.Username, 15 * time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println(token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
}