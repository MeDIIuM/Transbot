package command

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func Ethbalance(userMap map[string]*ecdsa.PrivateKey, update tgbotapi.Update, ethConnect *connect) string {
	balance, err := ethConnect.Balance(crypto.PubkeyToAddress(*userMap[update.Message.From.UserName].Public().(*ecdsa.PublicKey)))
	if err != nil {
		return fmt.Sprintf("balance failed %v", err)
	}

	return balance.String() + " ETH"
}

func (c *connect) Balance(address common.Address) (*big.Int, error) {
	balance, err := c.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return &big.Int{}, err
	}

	return balance, err
}
