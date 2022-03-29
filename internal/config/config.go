package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/logrusorgru/aurora/v3"
	"github.com/x6r/sip/internal/common"
)

type Config struct {
	Instance  string `json:"instance,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`
}

var (
	PathSip    = filepath.Join(xdg.ConfigHome, "sip")
	PathConfig = filepath.Join(PathSip, "config.json")
)

func Init() Config {
	var cfg Config

	if err := os.MkdirAll(PathSip, 0755); err != nil {
		common.Fatal(err)
	}

	if common.FileExists(PathConfig) {
		f, err := os.ReadFile(PathConfig)
		if err != nil {
			common.Fatal(err)
		}
		if err := json.Unmarshal(f, &cfg); err != nil {
			common.Fatal(err)
		}
	} else {
		fmt.Printf("%s Enter api instance › ", aurora.Green("?"))
		fmt.Scanln(&cfg.Instance)
		fmt.Printf("%s Enter auth token › ", aurora.Green("?"))
		fmt.Scanln(&cfg.AuthToken)

		configJson, err := json.Marshal(cfg)
		if err != nil {
			common.Fatal(err)
		}
		if err := os.WriteFile(PathConfig, configJson, 0644); err != nil {
			common.Fatal(err)
		}
	}

	return cfg
}
