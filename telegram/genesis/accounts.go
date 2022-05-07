package main

import (
	"context"
	"fmt"
)

func CreateAccounts(ctx context.Context, countAccount int, gasLimit int, gasPrice int) ([]*Account, error) {
	accounts := make([]*Account, countAccount)

	for i := 0; i < countAccount; i++ {
		account, err := GenerateAccount(gasLimit, gasPrice)
		if err != nil {
			return nil, fmt.Errorf("can't generate an account %w", err)
		}

		accounts[i] = account
	}

	return accounts, nil
}
