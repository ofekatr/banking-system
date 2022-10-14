package db

import (
	"context"
	"testing"

	"github.com/ofekatr/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, fromAccount)

	toAccount, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, toAccount)

	n := 5
	amount := util.RandomMoney()

	errs := make(chan error)
	results := make(chan *TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		require.Equal(t, fromAccount.ID, result.FromAccount.ID)
		dbFromAccount, err := store.GetAccount(context.Background(), result.FromAccount.ID)
		require.NoError(t, err)
		require.Equal(t, dbFromAccount.ID, result.FromAccount.ID)

		require.Equal(t, toAccount.ID, result.ToAccount.ID)
		dbToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)
		require.Equal(t, dbToAccount.ID, result.ToAccount.ID)

		require.NotZero(t, result.Transfer.ID)
		dbTransfer, err := store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.Equal(t, dbTransfer.Amount, result.Transfer.Amount)
		require.Equal(t, result.Transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, result.Transfer.ToAccountID, toAccount.ID)

		require.NotZero(t, result.FromEntry.ID)
		dbFromEntry, err := store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		require.Equal(t, dbFromEntry.Amount, -amount)
		require.Equal(t, dbFromEntry.AccountID, fromAccount.ID)

		require.NotZero(t, result.ToEntry.ID)
		dbToEntry, err := store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		require.Equal(t, dbToEntry.Amount, amount)
		require.Equal(t, dbToEntry.AccountID, toAccount.ID)

		require.NotZero(t, result.Transfer.CreatedAt)
	}
}
