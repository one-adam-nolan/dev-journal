package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configKeyDirectory = "directory"

// ViperProvider loads and persists journal configuration from ~/.djconfig.
type ViperProvider struct {
	loaded    bool
	configDir string
}

// NewViperProvider returns a config provider backed by Viper.
func NewViperProvider() *ViperProvider {
	return &ViperProvider{}
}

// NewViperProviderWithConfigDir returns a provider that reads and writes
// .djconfig in dir instead of the user home directory. Intended for tests.
func NewViperProviderWithConfigDir(dir string) *ViperProvider {
	return &ViperProvider{configDir: dir}
}

// Load reads or creates the config file at ~/.djconfig.
func (p *ViperProvider) Load() error {
	if p.loaded {
		return nil
	}

	configDir, err := p.configDirPath()
	if err != nil {
		return err
	}

	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigName(".djconfig")
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set(configKeyDirectory, filepath.Join(configDir, "Documents", "Dev-Journal"))

			if err := viper.SafeWriteConfigAs(filepath.Join(configDir, ".djconfig")); err != nil {
				return err
			}

			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	p.loaded = true
	return nil
}

func (p *ViperProvider) configDirPath() (string, error) {
	if p.configDir != "" {
		return p.configDir, nil
	}
	return os.UserHomeDir()
}

func (p *ViperProvider) JournalDirectory() string {
	return viper.GetString(configKeyDirectory)
}

func (p *ViperProvider) SetJournalDirectory(path string) error {
	viper.Set(configKeyDirectory, path)

	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return fmt.Errorf("no config file found")
	}

	return viper.WriteConfigAs(configFile)
}

func (p *ViperProvider) ConfigFilePath() string {
	return viper.ConfigFileUsed()
}

func (p *ViperProvider) RawConfig() ([]byte, error) {
	configFile := p.ConfigFilePath()
	if configFile == "" {
		return nil, fmt.Errorf("no config file found")
	}

	return os.ReadFile(configFile)
}
