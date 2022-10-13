package db

import (
	"context"
	"testing"

	"github.com/ofekatr/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	account2, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
}

func TestGetTransfer(t *testing.T) {
	transfer1, err := createRandomTransfer()
	require.NoError(t, err)
	require.NotEmpty(t, *transfer1)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
}

func TestListTransfers(t *testing.T) {
	account1, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	account2, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	for i := 0; i < 10; i++ {
		arg := CreateTransferParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        util.RandomMoney(),
		}

		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTransfersParams{
		Limit:         5,
		Offset:        5,
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func createRandomTransfer() (*Transfer, error) {
	account1, err := createRandomAccount()
	if err != nil {
		return nil, err
	}

	account2, err := createRandomAccount()
	if err != nil {
		return nil, err
	}

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}
