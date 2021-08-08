package cmd

import (
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
	"os"
	"runtime"
)

var (
	buildAt, gitHash       string
	verbose, isShowVersion bool
	db                     *gorm.DB
	cfg                    = new(model.CfgSchema)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "felix",
	Short: "命令来自我儿子的英文名字",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if isShowVersion {
			fmt.Println("Golang Env: ", runtime.Version(), runtime.GOOS, runtime.GOARCH)
			fmt.Println("UTC build time:", buildAt)
			fmt.Println("Build from Github repo version: https://github.com/bytegang/felix/commit/", gitHash)
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initFunc)
	rootCmd.Flags().BoolVarP(&isShowVersion, "version", "v", false, "show binary build information")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose")
}

func initFunc() {
	log.SetFlags(log.Lshortfile)
	db = model.LoadDatabase(verbose)
	model.CfgSyncDatabase(db, cfg)
}
