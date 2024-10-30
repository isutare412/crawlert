package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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
crawls:
  - name: 성남시 판교수영장
    enabled: true
    interval: 10s
    target:
      http:
        method: POST
        url: https://foo.com/api/reserve
        header:
          content-type: application/json
        body: |-
          {"message":"hello, world!"}
    query:
      check: |-
        [ .items[] | select(.name == "tester") ] | length > 0
      variables:
        TESTER_INFOS: |-
          [ .items[] | select(.name == "tester") | {"address": .address} ]
  - name: 화담숲 모노레일
    enabled: false
    interval: 5s
    target:
      http:
        method: POST
        url: https://bar.com/api/reserve
        header:
          content-type: application/json
        body: |-
          {"message":"bye, world!"}
    query:
      check: |-
        [ .items[] | select(.name == "tester2") ] | length > 0
      variables:
        TESTER_INFOS: |-
          [ .items[] | select(.name == "tester2") | {"address": .address} ]
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
				Crawls: []CrawlConfig{
					{
						Name:     "성남시 판교수영장",
						Enabled:  true,
						Interval: 10 * time.Second,
						Target: CrawlTargetConfig{HTTP: CrawlHTTPTargetConfig{
							Method: "POST",
							URL:    "https://foo.com/api/reserve",
							Header: map[string]string{"content-type": "application/json"},
							Body:   `{"message":"hello, world!"}`,
						}},
						Query: CrawlQueryConfig{
							Check: `[ .items[] | select(.name == "tester") ] | length > 0`,
							Variables: map[string]string{
								"TESTER_INFOS": `[ .items[] | select(.name == "tester") | {"address": .address} ]`,
							},
						},
					},
					{
						Name:     "화담숲 모노레일",
						Enabled:  false,
						Interval: 5 * time.Second,
						Target: CrawlTargetConfig{HTTP: CrawlHTTPTargetConfig{
							Method: "POST",
							URL:    "https://bar.com/api/reserve",
							Header: map[string]string{"content-type": "application/json"},
							Body:   `{"message":"bye, world!"}`,
						}},
						Query: CrawlQueryConfig{
							Check: `[ .items[] | select(.name == "tester2") ] | length > 0`,
							Variables: map[string]string{
								"TESTER_INFOS": `[ .items[] | select(.name == "tester2") | {"address": .address} ]`,
							},
						},
					},
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
