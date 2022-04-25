package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/logrusorgru/aurora/v3"
	"github.com/x6r/yumeko/internal/common"
)

type Config struct {
	Instance  string `json:"instance,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`
}

var (
	PathYumeko = filepath.Join(xdg.ConfigHome, ".yumeko")
	PathConfig = filepath.Join(PathYumeko, "config.json")
)

func Init() Config {
	var cfg Config

	if err := os.MkdirAll(PathYumeko, 0755); err != nil {
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

		configJson, err := json.MarshalIndent(cfg, "", "	")
		if err != nil {
			common.Fatal(err)
		}
		if err := os.WriteFile(PathConfig, configJson, 0644); err != nil {
			common.Fatal(err)
		}
	}

	return cfg
}
