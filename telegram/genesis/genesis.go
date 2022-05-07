package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
)

type Genesis interface {
	UpdateConfig(config map[string]interface{}) error
	UpdateAccounts(accounts []*Account, amount string) error
	UpdateGas(gasLimit uint64)
	UpdateExtraData(signers []common.Address, client string) error
}

func LoadGenesis(path, cluster string) (Genesis, error) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var gen Genesis

	var ethGen EthGenesis
	if err = json.Unmarshal(res, &ethGen); err != nil {
		return nil, err
	}

	gen = ethGen

	return gen, nil
}

func StoreGenesis(filename string, gen Genesis) error {
	bz, err := json.Marshal(gen)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bz, 0644)
}

func GetPeriod(client string, config map[string]interface{}) (int, error) {
	consensusName, _ := ConfigClient(client)

	consensus, _ := config[consensusName].(map[string]interface{})

	periodField, _ := BlockTimeFieldName(client)

	period, _ := consensus[periodField].(float64)

	return int(period), nil
}
