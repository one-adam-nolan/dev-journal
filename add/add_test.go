package add

import (
	"testing"

	"dj/config"

	"github.com/spf13/cobra"
)

func Test_addbullet(t *testing.T) {
	rootCmd := &cobra.Command{}
	config.InitConfig(rootCmd)
}
