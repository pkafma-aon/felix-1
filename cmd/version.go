package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// jsonCmd represents the json command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("编译时间:", buildAt)
		log.Println("Git Info:", gitHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
