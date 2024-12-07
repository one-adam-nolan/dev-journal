package add

import (
	"testing"

	"github.com/one-adam-nolan/dev-journal/cmds/config"
	"github.com/spf13/cobra"
)

func Test_addbullet(t *testing.T) {
	rootCmd := &cobra.Command{}
	config.InitConfig(rootCmd)

}
