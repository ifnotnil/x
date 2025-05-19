package conf

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type LoadArguments struct {
	Delimiter     string
	EnvVarsPrefix string
	EnvVarsDelim  string
	EnvVarsLayers map[string]any
	FileYAMLPath  string
	FileJSONPath  string
	FileEnvPath   string
}

func Load(ctx context.Context, logger *slog.Logger, la LoadArguments, defaults map[string]any) (*koanf.Koanf, error) {
	k := koanf.New(la.Delimiter)

	errs := make([]error, 0, 4)

	// Load default values.
	if len(defaults) > 0 {
		if err := k.Load(confmap.Provider(defaults, la.Delimiter), nil); err != nil {
			logger.DebugContext(ctx, "error during config loading from defaults", slog.String("error", err.Error()))
			errs = append(errs, err)
		}
	}

	// Load json config.
	if len(la.FileJSONPath) > 0 {
		if err := k.Load(file.Provider(la.FileJSONPath), json.Parser()); err != nil {
			logger.WarnContext(ctx, "error during config loading from json file", slog.String("error", err.Error()))
			errs = append(errs, err)
		}
	}

	// Load YAML config.
	if len(la.FileYAMLPath) > 0 {
		if err := k.Load(file.Provider(la.FileYAMLPath), yaml.Parser()); err != nil {
			logger.WarnContext(ctx, "error during config loading from yaml file", slog.String("error", err.Error()))
			errs = append(errs, err)
		}
	}

	mp := envVarNamesMapper(la)

	if err := k.Load(env.Provider(la.EnvVarsPrefix, la.Delimiter, mp), nil); err != nil {
		logger.WarnContext(ctx, "error during config loading from env vars", slog.String("error", err.Error()))
		errs = append(errs, err)
	}

	// Load .env file
	if len(la.FileEnvPath) > 0 {
		if err := k.Load(file.Provider(la.FileEnvPath), dotenv.ParserEnv(la.EnvVarsPrefix, la.Delimiter, mp)); err != nil {
			logger.WarnContext(ctx, "error during config loading from dot env file", slog.String("error", err.Error()))
			errs = append(errs, err)
		}
	}

	return k, errors.Join(errs...)
}

func envVarNamesMapper(la LoadArguments) func(string) string {
	return func(s string) string {
		s = strings.ToLower(s)
		s = strings.TrimPrefix(s, la.EnvVarsPrefix)
		return parseEnvVarName(la.EnvVarsDelim, la.Delimiter, la.EnvVarsLayers, s)
	}
}

func parseEnvVarName(envVarNameDelim string, delim string, layers map[string]any, name string) string {
	parts := strings.Split(name, envVarNameDelim)

	if len(parts) == 0 {
		return name
	}

	s := strings.Builder{}

	m := layers
	for i, p := range parts {
		if i > 0 {
			_, _ = s.WriteString(delim)
		}

		if len(m) == 0 {
			_, _ = s.WriteString(strings.Join(parts[i:], envVarNameDelim))
			break
		}

		inner, exists := m[p]
		if !exists || inner == nil {
			_, _ = s.WriteString(strings.Join(parts[i:], envVarNameDelim))
			break
		}

		asMap, is := inner.(map[string]any)
		if is {
			m = asMap
		}

		_, _ = s.WriteString(p)
	}

	return s.String()
}
