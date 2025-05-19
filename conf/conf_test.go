package conf

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnvVarName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		envVarNameDelim string
		delim           string
		layers          map[string]any
		name            string
		expected        string
	}{
		"nil": {
			envVarNameDelim: "_",
			delim:           ".",
			layers:          map[string]any{},
			name:            "",
			expected:        "",
		},

		"level 0": {
			envVarNameDelim: "_",
			delim:           ".",
			layers:          map[string]any{},
			name:            "a_b_c",
			expected:        "a_b_c",
		},

		"level 1": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{},
			},
			name:     "a_b_c",
			expected: "a.b_c",
		},

		"level 2": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": map[string]any{},
				},
			},
			name:     "a_b_c_d",
			expected: "a.b.c_d",
		},
		"level 2 with struct": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": struct{}{},
				},
			},
			name:     "a_b_c_d",
			expected: "a.b.c_d",
		},
		"level 2 with map": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": map[string]any{},
				},
			},
			name:     "a_b_c_d",
			expected: "a.b.c_d",
		},

		"level 3 edge": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{},
					},
				},
			},
			name:     "a_b_c",
			expected: "a.b.c",
		},

		"level 3 not found": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{},
					},
				},
			},
			name:     "a_b_d",
			expected: "a.b.d",
		},

		"level 3 with nil": {
			envVarNameDelim: "_",
			delim:           ".",
			layers: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": nil,
					},
				},
			},
			name:     "a_b_c_d",
			expected: "a.b.c_d",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := parseEnvVarName("_", ".", tc.layers, tc.name)
			if got != tc.expected {
				t.Errorf("expected %s got %s", tc.expected, got)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	ctx := context.Background()
	lvl := slog.Level(100)
	if testing.Verbose() {
		lvl = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	dr := t.TempDir()

	yamlFile := filepath.Join(dr, "config.yaml")
	jsonFile := filepath.Join(dr, "config.json")
	envFile := filepath.Join(dr, ".env")

	if err := os.WriteFile(yamlFile, []byte(yamlConfig), 0o600); err != nil {
		t.Fatalf("error writing yaml file: %s", err.Error())
	}
	if err := os.WriteFile(jsonFile, []byte(jsonConfig), 0o600); err != nil {
		t.Fatalf("error writing json file: %s", err.Error())
	}
	if err := os.WriteFile(envFile, []byte(envConfig), 0o600); err != nil {
		t.Fatalf("error writing env file: %s", err.Error())
	}

	t.Run("json file", func(t *testing.T) {
		k, err := Load(ctx, log, LoadArguments{Delimiter: ".", FileJSONPath: jsonFile}, map[string]any{})
		require.NoError(t, err)
		require.NotNil(t, k)

		assert.Equal(t, "dev", k.String("env_type"))
		assert.Equal(t, "info", k.String("log.level"))
		assert.Equal(t, "text", k.String("log.format"))
	})

	t.Run("yaml file", func(t *testing.T) {
		k, err := Load(ctx, log, LoadArguments{Delimiter: ".", FileYAMLPath: yamlFile}, map[string]any{})
		require.NoError(t, err)
		require.NotNil(t, k)

		assert.Equal(t, "production", k.String("env_type"))
		assert.Equal(t, "debug", k.String("log.level"))
		assert.Equal(t, "json", k.String("log.format"))
	})

	t.Run("env file", func(t *testing.T) {
		k, err := Load(
			ctx,
			log,
			LoadArguments{
				Delimiter:     ".",
				FileEnvPath:   envFile,
				EnvVarsPrefix: "",
				EnvVarsDelim:  "_",
				EnvVarsLayers: map[string]any{"log": struct{}{}},
			},
			map[string]any{},
		)
		require.NoError(t, err)
		require.NotNil(t, k)

		assert.Equal(t, "qa", k.String("env_type"))
		assert.Equal(t, "warn", k.String("log.level"))
		assert.Equal(t, "yaml", k.String("log.format"))
	})

	t.Run("mix", func(t *testing.T) {
		t.Setenv("LOG_FORMAT", "abc")

		k, err := Load(
			ctx,
			log,
			LoadArguments{
				Delimiter:     ".",
				FileJSONPath:  jsonFile,
				FileYAMLPath:  yamlFile,
				EnvVarsPrefix: "",
				EnvVarsDelim:  "_",
				EnvVarsLayers: map[string]any{"log": struct{}{}},
			},
			map[string]any{
				"log.opt": "ef",
			},
		)
		require.NoError(t, err)
		require.NotNil(t, k)

		assert.Equal(t, "production", k.String("env_type"))
		assert.Equal(t, "debug", k.String("log.level"))
		assert.Equal(t, "abc", k.String("log.format"))
		assert.Equal(t, "ef", k.String("log.opt"))
	})
}

const yamlConfig = `env_type: production
log:
  level: debug
  format: json
`

const jsonConfig = `{
  "env_type": "dev",
  "log": {
    "level": "info",
    "format": "text"
  }
}`

const envConfig = `ENV_TYPE=qa
LOG_LEVEL=warn
LOG_FORMAT=yaml
`
