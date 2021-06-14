package httpapi

import (
	"github.com/kelseyhightower/envconfig"

	"go.uber.org/zap"
)

type config struct {
	host    string
	port    string
	secret  []byte
	envName envName
}

type envName string

func (e envName) IsProd() bool {
	return e == "prod"
}

func newConfig() *config {
	envConfig := struct {
		Host    string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
		Port    string `envconfig:"HTTP_PORT" default:"8080"`
		Secret  string `envconfig:"AUTH_SECRET" default:"s3cr3t"`
		EnvName string `envconfig:"ENV_NAME" default:"local"`
	}{}

	_ = envconfig.Process("", &envConfig)

	return &config{
		host:    envConfig.Host,
		port:    envConfig.Port,
		secret:  []byte(envConfig.Secret),
		envName: envName(envConfig.EnvName),
	}
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	loggerConfig.DisableCaller = true
	loggerConfig.EncoderConfig.StacktraceKey = ""
	logger, _ := loggerConfig.Build()

	return logger.Named("api")
}
