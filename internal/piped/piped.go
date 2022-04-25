package piped

import (
	"fmt"
	"io"
	"net/http"

	"github.com/x6r/yumeko/internal/config"
)

func GetInstanceFromConfig() (string, string) {
	cfg := config.Init()
	feed := fmt.Sprintf("%s/feed?authToken=%s", cfg.Instance, cfg.AuthToken)
	return cfg.Instance, feed
}

func Get(url string) (*http.Response, []byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res, nil, err
	}

	return res, body, nil
}
