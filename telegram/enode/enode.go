package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	filenamePrivate     = "key"
	filenamePublic      = "key.pub"
	enodeFormat         = "Enode: enode://%s@"
	publicKeyFormat     = "\nPublic key: %s"
	publicAddressFormat = "\nAddress: %s"
)

func main() {
	importPath := flag.String("import", "", "the path to the file that contains the private key of the node")
	exportPath := flag.String("generate", "", "the path where you want to save the generated keys")

	flag.Parse()

	var (
		privateKey      *ecdsa.PrivateKey
		publicKeyString string
		err             error
		hexkey          []byte
	)

	switch true {
	case *importPath != "":
		hexkey, err = ioutil.ReadFile(*importPath)
		if err != nil {
			fmt.Println("can't read file with private key: ", err)

			return
		}

		privateKey, err = crypto.HexToECDSA(string(hexkey))
		if err != nil {
			fmt.Println("wrong private key: ", err)

			return
		}
	case *exportPath != "":
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			fmt.Println("can't generate key", err)

			return
		}

		defer func() {
			fmt.Printf(publicKeyFormat, publicKeyString)

			files := map[string]string{
				filenamePrivate: hexutil.Encode(crypto.FromECDSA(privateKey))[2:],
				filenamePublic:  publicKeyString,
			}

			for file, key := range files {
				err = ioutil.WriteFile(path.Join(*exportPath, file), []byte(key), 0755)
				if err != nil {
					fmt.Println("Can't write file: ", err)
				}
			}
		}()
	default:
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		fmt.Println("error casting public key to ECDSA")

		return
	}

	publicKeyString = hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))[4:]

	fmt.Printf(enodeFormat, publicKeyString)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf(publicAddressFormat, address)
}
