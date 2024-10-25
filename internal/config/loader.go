package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var reservedConfigFileNames = []string{
	"config.yaml",
	"config.local.yaml",
}

func Load(dir string) (*Config, error) {
	k := koanf.New(".")

	for i, f := range getConfigFilePaths(dir) {
		if i > 0 && !isFileExist(f) { // file is optional except first one
			continue
		}

		if err := loadFromFile(k, f); err != nil {
			return nil, fmt.Errorf("loading from file %s: %w", f, err)
		}
	}

	if err := loadFromEnv(k, "APP_"); err != nil {
		return nil, fmt.Errorf("loading from env: %w", err)
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, fmt.Errorf("unmarshaling into config struct: %w", err)
	}

	return &cfg, nil
}

func loadFromFile(k *koanf.Koanf, f string) error {
	if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
		return err
	}
	return nil
}

func loadFromEnv(k *koanf.Koanf, prefix string) error {
	// {prefix}FOO_BAR=baz -> foo.bar=baz
	if err := k.Load(env.Provider(prefix, "_", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil); err != nil {
		return err
	}
	return nil
}

func getConfigFilePaths(dir string) []string {
	files := make([]string, 0, len(reservedConfigFileNames))
	for _, name := range reservedConfigFileNames {
		files = append(files, filepath.Join(dir, name))
	}
	return files
}

func isFileExist(file string) bool {
	_, err := os.Stat(file)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false
	case err != nil:
		panic(fmt.Errorf("checking file existence: %w", err))
	}
	return true
}
