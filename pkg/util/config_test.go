package util

import (
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	initialConfig  Config
	path           string
	expectedError  string
	expectedConfig Config
}

func TestLoadConfig(t *testing.T) {
	okHome := Config{
		GRPCAddress: "localhost:1212",
		Insecure:    true,
		Token:       "UUID",
		TeamSide:    lugo.Team_HOME,
		Number:      4,
	}
	okAway := okHome
	okAway.TeamSide = lugo.Team_AWAY
	caseList := map[string]testCase{
		"ok": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
				TeamSide:    lugo.Team_AWAY,
			},
			path:           "testdata/config_test_ok.json",
			expectedConfig: okHome,
		},
		"ok_away": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:           "testdata/config_test_ok_away.json",
			expectedConfig: okAway,
		},
		"ok_team_cap": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:           "testdata/config_test_ok_team_capitals.json",
			expectedConfig: okHome,
		},
		"team undefined": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:          "testdata/config_test_invalid_home_undefined.json",
			expectedError: "invalid team option",
		},
		"number 0": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:          "testdata/config_test_invalid_number_0.json",
			expectedError: "invalid player number",
		},
		"number 12": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:          "testdata/config_test_invalid_number_12.json",
			expectedError: "invalid player number",
		},
		"file not found": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:          "testdata/no-file.json",
			expectedError: "no such file or director",
		},
		"invalid json": {
			initialConfig: Config{
				GRPCAddress: "localhost:1212",
				Number:      2,
			},
			path:          "testdata/config_test_invalid_json.json",
			expectedError: "error parsing the config",
		},
	}

	for caseName, tCase := range caseList {
		err := LoadConfig(tCase.path, &tCase.initialConfig)
		if err == nil {
			assert.Equal(t, tCase.expectedConfig, tCase.initialConfig, caseName)
		} else {
			assert.Contains(t, err.Error(), tCase.expectedError, caseName)
		}
	}
}
