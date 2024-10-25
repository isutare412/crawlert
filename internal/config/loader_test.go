package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/crawlert/internal/log"
)

var (
	testConfig = `
log:
  format: text # text / json
  level: debug # debug / info / warn / error
  caller: true
alerts:
  telegram:
    bot-token: test-bot-token
    chat-ids:
      - test-chat-id
`

	testLocalConfig = `
log:
  level: error
`
)

func TestLoad(t *testing.T) {
	type args struct {
		cfg      string
		localCfg string
		envs     map[string]string
	}
	tests := []struct {
		name         string
		args         args
		want         *Config
		wantLogLevel log.Level
	}{
		{
			name: "load_from_file",
			args: args{
				cfg: testConfig,
			},
			want: &Config{
				Log: LogConfig{
					Format: log.FormatText,
					Level:  log.LevelDebug,
					Caller: true,
				},
				Alerts: AlertsConfig{
					Telegram: TelegramConfig{
						BotToken: "test-bot-token",
						ChatIDs:  []string{"test-chat-id"},
					},
				},
			},
		},
		{
			name: "override_by_local_config",
			args: args{
				cfg:      testConfig,
				localCfg: testLocalConfig,
			},
			wantLogLevel: log.LevelError,
		},
		{
			name: "override_by_envs",
			args: args{
				cfg: testConfig,
				envs: map[string]string{
					"APP_LOG_LEVEL": "error",
				},
			},
			wantLogLevel: log.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := prepareTestEnvironment(t, tt.args.cfg, tt.args.localCfg, tt.args.envs)

			got, err := Load(testDir)
			require.NoError(t, err)

			switch {
			case tt.wantLogLevel != "":
				assert.Equal(t, tt.wantLogLevel, got.Log.Level)
			case tt.want != nil:
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func prepareTestEnvironment(t *testing.T, cfg, localCfg string, envs map[string]string) (testDir string) {
	tempDir := t.TempDir()

	writeToFile := func(fileName, cfgBody string) {
		file, err := os.Create(filepath.Join(tempDir, fileName))
		require.NoError(t, err)
		defer file.Close()

		_, err = file.WriteString(cfgBody)
		require.NoError(t, err)
	}

	writeToFile("config.yaml", cfg)
	if localCfg != "" {
		writeToFile("config.local.yaml", localCfg)
	}

	for k, v := range envs {
		t.Setenv(k, v)
	}

	return tempDir
}
