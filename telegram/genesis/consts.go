package main

type GenesysPath struct {
	Original string
	Dest     string
}

func GenesysPaths(c string) (GenesysPath, bool) {
	var res *GenesysPath

	res = &GenesysPath{
		"clients/geth/docker/genesis_original.json",
		"clients/geth/docker/genesis.json",
	}

	if res == nil {
		return GenesysPath{}, false
	}

	return *res, true
}

func ConfigClient(c string) (string, bool) {
	const clique = "clique"

	return clique, true
}

func BlockTimeFieldName(c string) (string, bool) {
	const period = "period"

	return period, true
}
