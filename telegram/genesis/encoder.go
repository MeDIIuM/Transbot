package main

import (
	"bytes"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

type Encoder interface {
	Encode() (string, error)
}

type CliqueExtraData struct {
	validators []common.Address
}

func NewCliqueExtraData(validators []common.Address) CliqueExtraData {
	return CliqueExtraData{
		validators: validators,
	}
}

func NewEncoder(cluster string, signers []common.Address) (Encoder, error) {
	validators := make([]common.Address, len(signers))
	copy(validators, signers)

	return NewCliqueExtraData(validators), nil
}

func (clique CliqueExtraData) Encode() (string, error) {
	const (
		prefix = 32 // vanity
		suffix = 65
	)

	sort.SliceStable(clique.validators, func(i, j int) bool {
		return bytes.Compare(clique.validators[i][:], clique.validators[j][:]) > 0
	})

	bz := make([]byte, prefix+len(clique.validators)*common.AddressLength+suffix)

	for i, validator := range clique.validators {
		copy(bz[prefix+i*common.AddressLength:], validator[:])
	}

	return hexutil.Encode(bz), nil
}

type IBFT2ExtraData struct {
	validators []common.Address
}

func NewIBFT2ExtraData(validators []common.Address) IBFT2ExtraData {
	return IBFT2ExtraData{
		validators: validators,
	}
}

func (ibft2 IBFT2ExtraData) Encode() (string, error) {
	const (
		vanity = 32
		vote   = 0
		round  = 4
	)
	// last one is []seals
	payload, err := rlp.EncodeToBytes([]interface{}{[vanity]byte{}, &ibft2.validators, [vote]byte{}, [round]byte{}, []string{}})
	if err != nil {
		return "", err
	}

	return hexutil.Encode(payload), nil
}
