package config

// Provider exposes journal configuration.
type Provider interface {
	Load() error
	JournalDirectory() string
	SetJournalDirectory(path string) error
	ConfigFilePath() string
	RawConfig() ([]byte, error)
}
