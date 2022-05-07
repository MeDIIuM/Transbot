package command

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type connect struct {
	*ethclient.Client
}

func Transfer(userMap map[string]*ecdsa.PrivateKey, update tgbotapi.Update, ethConnect *connect) string {
	pk := userMap[update.Message.From.UserName]

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) != 2 {
		return "Недостаточно аргументов: логин_адресата количество"
	}

	pkTo := userMap[args[0]]
	addressTo := crypto.PubkeyToAddress(*pkTo.Public().(*ecdsa.PublicKey))
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Sprintf("amount incorrect %v", err.Error())
	}

	err = ethConnect.MoneyTransfer(pk, addressTo, big.NewInt(int64(amount)))
	if err != nil {
		return fmt.Sprintf("money transfer %v", err.Error())
	}

	return "Money transfer success"
}

func (c *connect) MoneyTransfer(pk *ecdsa.PrivateKey, addressTo common.Address, amount *big.Int) error {
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("nonce %w", err)
	}

	gasLimit := uint64(21000)
	gasPrice, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("gas price %w", err)
	}

	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &addressTo,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signTx, err := types.SignTx(tx, types.HomesteadSigner{}, pk)
	if err != nil {
		return fmt.Errorf("sign tx %w", err)
	}

	err = c.SendTransaction(context.Background(), signTx)
	if err != nil {
		return fmt.Errorf("send transaction %w", err)
	}

	return nil
}
