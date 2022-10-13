package db

import (
	"context"
	"testing"

	"github.com/ofekatr/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
}

func TestGetEntry(t *testing.T) {
	entry, err := createRandomEntry()
	require.NoError(t, err)
	require.NotEmpty(t, *entry)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, entry.Amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	account, err := createRandomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, account)

	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}

		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func createRandomEntry() (*Entry, error) {
	account, err := createRandomAccount()
	if err != nil {
		return nil, err
	}

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
