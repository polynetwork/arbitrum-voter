package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Config struct {
	sync.Once
	PolyConfig       PolyConfig
	ArbConfig        ArbConfig
	ForceConfig      ForceConfig
	BoltDbPath       string
	WhitelistMethods []string
	whitelistMethods map[string]bool
}

func (c *Config) IsWhitelistMethod(method string) bool {
	c.Do(func() {
		c.whitelistMethods = map[string]bool{}
		for _, m := range c.WhitelistMethods {
			c.whitelistMethods[m] = true
		}
	})

	return c.whitelistMethods[method]
}

type PolyConfig struct {
	RestURL    string
	WalletFile string
}

type ArbConfig struct {
	SideChainId         uint64
	ECCMContractAddress string
	RestURL             []string
}

type ForceConfig struct {
	ArbHeight uint64
}

func LoadConfig(confFile string) (config *Config, err error) {
	jsonBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return
	}

	config = &Config{}
	err = json.Unmarshal(jsonBytes, config)
	return
}
