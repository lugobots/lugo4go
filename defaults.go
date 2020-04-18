package lugo4go

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func DefaultLogger(config lugo.Config) (*zap.SugaredLogger, error) {
	configZap := zap.NewDevelopmentConfig()
	configZap.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, err := configZap.Build()
	if err != nil {
		return nil, fmt.Errorf("could not initiliase looger: %s", err)
	}
	return zapLog.Sugar().Named(fmt.Sprintf("%s-%d", config.TeamSide, config.Number)), nil
}

func DefaultConfigurator() (lugo.Config, error) {
	config, err := lugo.LoadConfig("./config.json")
	if err != nil {
		return lugo.Config{}, fmt.Errorf("did not load the config: %s", err)
	}
	if err := config.ParseConfigFlags(); err != nil {
		return lugo.Config{}, fmt.Errorf("did not parsed well the flags for config: %s", err)
	}

	return config, nil
}

func DefaultInitBundle() (lugo.Config, *zap.SugaredLogger, error) {
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
