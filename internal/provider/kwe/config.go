package kwe

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	apiURL string
}

func newConfig() (*config, error) {
	envConfig := struct {
		URL string `envconfig:"KWE_URL" default:"http://0.0.0.0:5000/yake/"`
	}{}

	err := envconfig.Process("", &envConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env config - %s", err)
	}

	cfg := &config{apiURL: envConfig.URL}

	return cfg, nil
}
