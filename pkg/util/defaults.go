package util

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// DefaultLogger creates a logger that is compatible with the lugo4go.Handler expected logger.
// The bots are NOT obligated to use this logger though. You may implement your own logger.
func DefaultLogger(config Config) (*zap.SugaredLogger, error) {
	configZap := zap.NewDevelopmentConfig()
	configZap.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, err := configZap.Build()
	if err != nil {
		return nil, fmt.Errorf("could not initiliase looger: %s", err)
	}
	return zapLog.Sugar().Named(fmt.Sprintf("%s-%d", config.TeamSide, config.Number)), nil
}

// DefaultInitBundle created a basic configuration that may be used by the client to connect to the server.
// It also creates a logger that is compatible with the lugo4go.Handler.
func DefaultInitBundle() (Config, *zap.SugaredLogger, error) {
	config := Config{}

	if err := config.LoadConfig(os.Args[1:]); err != nil {
		return config, nil, fmt.Errorf("did not parsed well the flags for config: %s", err)
	}

	logger, err := DefaultLogger(config)
	if err != nil {
		return config, nil, err
	}
	return config, logger, nil
}
