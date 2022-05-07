package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type Clique struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}
type CliqueConfig struct {
	Clique Clique `json:"clique"`
}

type IBFT2 struct {
	Period  uint64 `json:"blockperiodseconds"`
	Epoch   uint64 `json:"epochlength"`
	Timeout uint64 `json:"requesttimeoutseconds"`
}

type IBFT2Config struct {
	IBFT2 IBFT2 `json:"ibft2"`
}

type EthGenesis map[string]interface{}

func (g EthGenesis) UpdateConfig(newConfig map[string]interface{}) error {
	for name, val := range newConfig {
		switch config := g["config"].(type) {
		case map[string]interface{}:
			config[name] = val
		default:
		}
	}

	return nil
}

func (g EthGenesis) UpdateGas(newGasLimit uint64) {
	gasLimit, _ := math.HexOrDecimal64(newGasLimit).MarshalText()
	g["gasLimit"] = string(gasLimit)
}

func (g EthGenesis) UpdateAccounts(accounts []*Account, amount string) error {
	alloc := make(map[string]interface{})

	for _, account := range accounts {
		alloc[account.Address.Hex()[2:]] = map[string]string{
			"balance": amount,
		}
	}

	switch accounts := g["alloc"].(type) {
	case map[string]interface{}:
		for address, value := range alloc {
			accounts[address] = value
		}
	default:
	}

	return nil
}

func (g EthGenesis) UpdateExtraData(signers []common.Address, cluster string) error {
	consensus, err := NewEncoder(cluster, signers)
	if err != nil {
		return err
	}

	extraData, err := consensus.Encode()
	if err != nil {
		return err
	}

	g["extraData"] = extraData

	return nil
}
