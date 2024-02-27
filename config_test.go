package lugo4go

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lugobots/lugo4go/v3/proto"
)

func TestLoadConfig(t *testing.T) {
	caseList := map[string]struct {
		args           []string
		env            map[string]string
		expectedError  string
		expectedConfig Config
	}{
		"default values": {
			expectedError: "",
			expectedConfig: Config{
				GRPCAddress: "localhost:5000",
				TeamSide:    proto.Team_HOME,
				Number:      0,
				Token:       "",
				Insecure:    false,
			},
		},
		"every thing changed from flags": {
			args: []string{
				"--grpc_address", "localhost:1212",
				"--team", "home",
				"--number", "7",
				"--token", "a-token",
				"--insecure", "no",
			},
			expectedError: "",
			expectedConfig: Config{
				GRPCAddress: "localhost:1212",
				TeamSide:    proto.Team_AWAY,
				Number:      7,
				Token:       "a-token",
				Insecure:    false,
			},
		},
		"every thing changed from environment variables": {
			env: map[string]string{
				EnvVarBotGrpcUrl:      "another-host:5000",
				EnvVarBotTeam:         "HOME",
				EnvVarBotNumber:       "5",
				EnvVarBotToken:        "another-token",
				EnvVarBotGrpcInsecure: "false",
			},
			expectedError: "",
			expectedConfig: Config{
				GRPCAddress: "another-host:5000",
				TeamSide:    proto.Team_HOME,
				Number:      4,
				Token:       "another-token",
				Insecure:    false,
			},
		},
		"env variables should override params": {
			args: []string{
				"--grpc_address", "localhost:1212",
				"--team", "away",
				"--number", "7",
				"--token", "a-token",
				"--insecure", "yes",
			},
			env: map[string]string{
				EnvVarBotGrpcUrl:      "another-host:5000",
				EnvVarBotTeam:         "HOME",
				EnvVarBotNumber:       "5",
				EnvVarBotToken:        "another-token",
				EnvVarBotGrpcInsecure: "false",
			},
			expectedError: "",
			expectedConfig: Config{
				GRPCAddress: "another-host:5000",
				TeamSide:    proto.Team_HOME,
				Number:      4,
				Token:       "another-token",
				Insecure:    false,
			},
		},
		"invalid team": {
			env: map[string]string{
				EnvVarBotTeam: "another",
			},
			expectedError:  "invalid team option",
			expectedConfig: Config{},
		},
		"invalid player number - greater than 11": {
			env: map[string]string{
				EnvVarBotNumber: "12",
			},
			expectedError:  "invalid player number",
			expectedConfig: Config{},
		},
		"invalid player number - less than 1": {
			env: map[string]string{
				EnvVarBotNumber: "0",
			},
			expectedError:  "invalid player number",
			expectedConfig: Config{},
		},
		"invalid insecure flag": {
			env: map[string]string{
				EnvVarBotGrpcInsecure: "maybe",
			},
			expectedError:  "invalid gRPC insecure flag",
			expectedConfig: Config{},
		},
	}

	for caseName, tc := range caseList {
		t.Run(caseName, func(t *testing.T) {
			defer func() {
				for varName := range tc.env {
					_ = os.Unsetenv(varName)
				}
			}()
			for varName, varValue := range tc.env {
				_ = os.Setenv(varName, varValue)
			}

			config := Config{}
			err := config.loadConfig(tc.args)
			if err == nil {
				assert.Equal(t, tc.expectedConfig.GRPCAddress, config.GRPCAddress, caseName)
			} else {
				assert.Contains(t, err.Error(), tc.expectedError, caseName)
			}
		})

	}
}
