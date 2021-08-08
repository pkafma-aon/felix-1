package cmd

import (
	"github.com/bytegang/sshd/svr"
	"github.com/spf13/cobra"
)

var sshdCmd = &cobra.Command{
	Use:   "sshd",
	Short: "启动SSH-server",
	Long: `启动堡垒机
`,
	Run: func(c *cobra.Command, args []string) {

		svr.Run(&svr.DefaultCfg)

	},
}

func init() {
	rootCmd.AddCommand(sshdCmd)
}
