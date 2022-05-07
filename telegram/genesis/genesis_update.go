package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"vanilla-achivs/api/cluster"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	givenClient := flag.String("client", "besu-ibft", "Name client")
	genesisFile := flag.String("src", "geth_genesis.json", "Path to original genesis")
	genesisDestFile := flag.String("dest", "ibft2Genesis_original.json", "Path to dest genesis")
	validators := flag.String("validators", "", "list of validators")
	gasLimit := flag.Uint64("gasLimit", 0, "Gas limit block")
	accounts := flag.String("accounts", "", "Alloc accounts")

	flag.Parse()

	gen, err := LoadGenesis(*genesisFile, *givenClient)
	if err != nil {
		return
	}

	if *gasLimit != 0 {
		gen.UpdateGas(*gasLimit)
	}

	//if len(*accounts) != 0 {
	allocs := strings.Split(*accounts, ",")
	accounts2 := make([]*Account, len(allocs))

	addr := cluster.Addresses{
		Addresses: make([]string, len(allocs)),
	}

	if len(*validators) != 0 {
		sealers := strings.Split(*validators, ",")
		signers := make([]common.Address, len(sealers))

		addr.Addresses[0] = sealers[0]

		for i, sealer := range sealers {
			signers[i] = common.HexToAddress(sealer)

			accounts2 = append(accounts2, &Account{
				Address: common.HexToAddress(sealer),
			})
		}

		if err := gen.UpdateExtraData(signers, *givenClient); err != nil {
			return
		}
	}

	for i, alloc := range allocs {
		accounts2[i] = &Account{
			Address: common.HexToAddress(alloc),
		}
	}

	if err := gen.UpdateAccounts(accounts2, "0x200000000000000000000000000000"); err != nil {

		return
	}
	//}

	if err := StoreGenesis(*genesisDestFile, gen); err != nil {

		return
	}

	addrJSON, err := json.Marshal(addr)
	if err != nil {
		log.Fatal(err)
	}

	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, addrJSON, "", "\t")
	if err != nil {
		log.Fatal(err)

	}

	err = ioutil.WriteFile("../addr.json", prettyJSON.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
