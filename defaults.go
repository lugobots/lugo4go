package lugo4go

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
)

// DefaultLogger creates a logger that is compatible with the lugo4go.rawBot expected logger.
// The bots are NOT obligated to use this logger though. You may implement your own logger.
func DefaultLogger(config Config) *zap.SugaredLogger {
	configZap := zap.NewDevelopmentConfig()
	configZap.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog, _ := configZap.Build()
	// no need to check the error since this configuration won cause any
	return zapLog.Sugar().Named(fmt.Sprintf("%s-%d", config.TeamSide, config.Number))
}

// DefaultInitBundle created a basic configuration that may be used by the client to connect to the server.
// It also creates a logger that is compatible with the lugo4go.rawBot.
func DefaultInitBundle() (Config, mapper.Mapper, *zap.SugaredLogger, error) {
	config := Config{}

	if err := config.loadConfig(os.Args[1:]); err != nil {
		return config, nil, nil, fmt.Errorf("did not parsed well the flags for config: %s", err)
	}

	// default initial position
	defaultMapper, _ := mapper.NewMapper(16, 8, config.TeamSide)
	defaultFieldMap := DefaultRoleMap["initial"]
	defaultInitialRegion, _ := defaultMapper.GetRegion(defaultFieldMap[int(config.Number)].Col, defaultFieldMap[int(config.Number)].Row)
	config.InitialPosition = defaultInitialRegion.Center()

	config.PlayerPositionFn = func(playerNumber int, inspector SnapshotInspector) *proto.Point {
		ballRegion, _ := defaultMapper.GetPointRegion(inspector.GetBall().GetPosition())
		regionCol := ballRegion.Col()

		teamState := "neutral"
		if regionCol > 7 {
			teamState = "regionCol"
		} else if regionCol > 4 {
			teamState = "neutral"
		}
		pos := DefaultRoleMap[teamState][playerNumber]
		reg, _ := defaultMapper.GetRegion(pos.Col, pos.Row)
		return reg.Center()
	}

	logger := DefaultLogger(config)
	return config, defaultMapper, logger, nil
}

var DefaultRoleMap = map[string]map[int]struct {
	Col int
	Row int
}{

	"initial": {
		2:  {Col: 1, Row: 2},
		3:  {Col: 1, Row: 3},
		4:  {Col: 1, Row: 4},
		5:  {Col: 1, Row: 5},
		6:  {Col: 4, Row: 1},
		7:  {Col: 4, Row: 3},
		8:  {Col: 4, Row: 4},
		9:  {Col: 4, Row: 6},
		10: {Col: 6, Row: 3},
		11: {Col: 6, Row: 4},
	},
	"defensive": {
		2:  {Col: 1, Row: 2},
		3:  {Col: 1, Row: 3},
		4:  {Col: 1, Row: 4},
		5:  {Col: 1, Row: 5},
		6:  {Col: 4, Row: 1},
		7:  {Col: 4, Row: 3},
		8:  {Col: 4, Row: 4},
		9:  {Col: 4, Row: 6},
		10: {Col: 6, Row: 3},
		11: {Col: 6, Row: 4},
	},
	"neutral": {
		2:  {Col: 3, Row: 1},
		3:  {Col: 3, Row: 3},
		4:  {Col: 3, Row: 4},
		5:  {Col: 3, Row: 6},
		6:  {Col: 6, Row: 1},
		7:  {Col: 6, Row: 3},
		8:  {Col: 6, Row: 4},
		9:  {Col: 6, Row: 6},
		10: {Col: 10, Row: 2},
		11: {Col: 10, Row: 5},
	},
	"offensive": {
		2:  {Col: 5, Row: 1},
		3:  {Col: 4, Row: 3},
		4:  {Col: 4, Row: 4},
		5:  {Col: 5, Row: 6},
		6:  {Col: 9, Row: 2},
		7:  {Col: 11, Row: 1},
		8:  {Col: 9, Row: 5},
		9:  {Col: 11, Row: 6},
		10: {Col: 13, Row: 2},
		11: {Col: 13, Row: 5},
	},
}
