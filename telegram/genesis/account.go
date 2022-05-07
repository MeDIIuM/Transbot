package main

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Auth    *bind.TransactOpts
	Key     *ecdsa.PrivateKey
	Address common.Address
}

func GenerateAccount(gasLimit, gasPrice int) (*Account, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	acc := &Account{
		Auth:    bind.NewKeyedTransactor(key),
		Key:     key,
		Address: crypto.PubkeyToAddress(*key.Public().(*ecdsa.PublicKey)),
	}

	acc.Auth.GasLimit = uint64(gasLimit)
	acc.Auth.GasPrice = big.NewInt(int64(gasPrice))
	acc.Auth.Nonce = big.NewInt(0)

	return acc, nil
}
