package lugo4go

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
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

	if err := config.loadConfig(os.Args[1:]); err != nil {
		return config, nil, fmt.Errorf("did not parsed well the flags for config: %s", err)
	}

	// default initial position
	tempMapper, _ := mapper.NewMapper(11, 7, config.TeamSide)
	defaultFieldMap := DefaultRoleMap["initial"]
	defaultInitialRegion, _ := tempMapper.GetRegion(defaultFieldMap[int(config.Number)].Col, defaultFieldMap[int(config.Number)].Row)
	config.InitialPosition = defaultInitialRegion.Center()

	config.PlayerPositionFn = func(playerNumber int, inspector SnapshotInspector) *proto.Point {
		ballRegion, _ := tempMapper.GetPointRegion(inspector.GetBall().GetPosition())
		regionCol := ballRegion.Col()

		teamState := "neutral"
		if regionCol > 7 {
			teamState = "regionCol"
		} else if regionCol > 4 {
			teamState = "neutral"
		}
		pos := DefaultRoleMap[teamState][playerNumber]
		reg, _ := tempMapper.GetRegion(pos.Col, pos.Row)
		return reg.Center()
	}

	logger, err := DefaultLogger(config)
	if err != nil {
		return config, nil, err
	}
	return config, logger, nil
}

var DefaultRoleMap = map[string]map[int]struct {
	Col int
	Row int
}{

	"initial": {
		2:  {Col: 1, Row: 1},
		3:  {Col: 1, Row: 3},
		4:  {Col: 1, Row: 4},
		5:  {Col: 1, Row: 6},
		6:  {Col: 2, Row: 2},
		7:  {Col: 2, Row: 3},
		8:  {Col: 2, Row: 4},
		9:  {Col: 2, Row: 5},
		10: {Col: 3, Row: 3},
		11: {Col: 3, Row: 4},
	},
	"defensive": {
		2:  {Col: 1, Row: 1},
		3:  {Col: 2, Row: 2},
		4:  {Col: 2, Row: 3},
		5:  {Col: 1, Row: 4},
		6:  {Col: 3, Row: 1},
		7:  {Col: 3, Row: 2},
		8:  {Col: 3, Row: 3},
		9:  {Col: 3, Row: 4},
		10: {Col: 4, Row: 3},
		11: {Col: 4, Row: 2},
	},
	"neutral": {
		2:  {Col: 2, Row: 1},
		3:  {Col: 4, Row: 2},
		4:  {Col: 4, Row: 3},
		5:  {Col: 2, Row: 4},
		6:  {Col: 6, Row: 1},
		7:  {Col: 8, Row: 2},
		8:  {Col: 8, Row: 3},
		9:  {Col: 6, Row: 4},
		10: {Col: 7, Row: 4},
		11: {Col: 7, Row: 1},
	},
	"offensive": {
		2:  {Col: 3, Row: 1},
		3:  {Col: 5, Row: 2},
		4:  {Col: 5, Row: 3},
		5:  {Col: 3, Row: 4},
		6:  {Col: 7, Row: 1},
		7:  {Col: 8, Row: 2},
		8:  {Col: 8, Row: 3},
		9:  {Col: 7, Row: 4},
		10: {Col: 9, Row: 4},
		11: {Col: 9, Row: 1},
	},
}
