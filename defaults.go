package client

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func DefaultLogger(config Config) (lugo.Logger, error) {
	configZap := zap.NewDevelopmentConfig()
	configZap.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, err := configZap.Build()
	if err != nil {
		return nil, fmt.Errorf("could not initiliase looger: %s", err)
	}
	return zapLog.Sugar().Named(fmt.Sprintf("%s-%d", config.TeamSide, config.Number)), nil
}

func DefaultConfigurator() (Config, error) {
	config, err := LoadConfig("./config.json")
	if err != nil {
		return Config{}, fmt.Errorf("did not load the config: %s", err)
	}
	if err := config.ParseConfigFlags(); err != nil {
		return Config{}, fmt.Errorf("did not parsed well the flags for config: %s", err)
	}

	return config, nil
}

func DefaultBundle() (Config, lugo.Logger, error) {
	config, err := DefaultConfigurator()
	if err != nil {
		return config, nil, err
	}
	logger, err := DefaultLogger(config)
	if err != nil {
		return config, nil, err
	}
	return config, logger, nil
}
