package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) // no error
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

// 1.create a account
// 2.update the account
// 3.check error and empty
// 4.check if the update params equal to account's updated value
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	// all segment except Balance and CreateAt should be equal
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	// skip the first 5 records , and return the next 5
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
