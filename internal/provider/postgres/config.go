package postgres

import (
	"crypto/aes"
	"fmt"
	"net/url"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type config struct {
	secret        string
	connectionURL string
	logVerbose    bool
}

func newConfig() (*config, error) {
	envConfig := struct {
		URL        string `envconfig:"DB_URL" required:"true"`
		Secret     string `envconfig:"DB_SECRET" required:"true"`
		LogVerbose bool   `envconfig:"DB_LOG_VERBOSE" default:"false"`
	}{}

	err := envconfig.Process("", &envConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env config - %s", err)
	}

	// add sslmode=disable
	dbURL, err := url.Parse(envConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse DB_URL env - %s", err)
	}

	if dbURLQueryValues := dbURL.Query(); dbURLQueryValues.Get("sslmode") == "" {
		dbURLQueryValues.Set("sslmode", "disable")
		dbURL.RawQuery = dbURLQueryValues.Encode()
		envConfig.URL = dbURL.String()
	}

	_, err = aes.NewCipher([]byte(envConfig.Secret))
	if err != nil {
		return nil, err
	}

	cfg := &config{
		secret:        envConfig.Secret,
		connectionURL: envConfig.URL,
		logVerbose:    envConfig.LogVerbose,
	}

	return cfg, nil
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	loggerConfig.DisableCaller = true
	loggerConfig.EncoderConfig.StacktraceKey = ""
	logger, _ := loggerConfig.Build()

	return logger.Named("db")
}
