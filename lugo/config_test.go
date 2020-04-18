package lugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	path           string
	expectedError  string
	expectedConfig Config
}

func TestLoadConfig(t *testing.T) {
	okHome := Config{
		GRPCAddress: "localhost:1212",
		Insecure:    true,
		Token:       "UUID",
		TeamSide:    Team_HOME,
		Number:      4,
	}
	okAway := okHome
	okAway.TeamSide = Team_AWAY
	caseList := map[string]testCase{
		"ok": {
			path:           "testdata/config_test_ok.json",
			expectedConfig: okHome,
		},
		"ok_away": {
			path:           "testdata/config_test_ok_away.json",
			expectedConfig: okAway,
		},
		"ok_team_cap": {
			path:           "testdata/config_test_ok_team_capitals.json",
			expectedConfig: okHome,
		},
		"team undefined": {
			path:          "testdata/config_test_invalid_home_undefined.json",
			expectedError: "invalid team option",
		},
		"number 0": {
			path:          "testdata/config_test_invalid_number_0.json",
			expectedError: "invalid player number",
		},
		"number 12": {
			path:          "testdata/config_test_invalid_number_12.json",
			expectedError: "invalid player number",
		},
		"file not found": {
			path:          "testdata/no-file.json",
			expectedError: "no such file or director",
		},
		"invalid json": {
			path:          "testdata/config_test_invalid_json.json",
			expectedError: "error parsing the config",
		},
	}

	for caseName, tCase := range caseList {
		config, err := LoadConfig(tCase.path)
		if err == nil {
			assert.Equal(t, tCase.expectedConfig, config, caseName)
		} else {
			assert.Contains(t, err.Error(), tCase.expectedError, caseName)
		}
	}
}
