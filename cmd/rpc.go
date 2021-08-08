package cmd

import (
	"github.com/bytegang/felix/rpc"
	"github.com/spf13/cobra"
	"log"
	"net"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "启动RPC服务",
	Long: `为bytegang/sshd 提供开发测试环境
`,
	Run: func(c *cobra.Command, args []string) {
		lisRpc, err := net.Listen("tcp", cfg.AddrRpc)
		if err != nil {
			log.Fatal(err)
		}
		cfg.AddrRpc = lisRpc.Addr().String()
		rpc.NewSvrRpc(db, cfg).Run(lisRpc)
	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)
}
