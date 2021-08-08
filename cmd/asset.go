package cmd

import (
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
	"github.com/bytegang/sshd/sshterm"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "通过命令行来管理登录SSH",
	Long: `eg: felix asset  表格展示资产
eg: felix asset 2  命令行直接SSH登录ID=2机器
`,
	Run: func(c *cobra.Command, args []string) {
		//1. display machine
		webLis, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			log.Fatal(err)
		}
		displayMachineList(query, webLis)

		//命令行模式
		if len(args) > 0 {
			machineId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				log.Fatal("ssh ID must be a int:", err)
			}
			one := new(model.Machine)
			err = db.First(one, machineId).Error
			if err != nil {
				log.Fatal("ssh ID must be a int:", err)
			}
			arg := pb.SshConn{
				Addr:               fmt.Sprintf("[%s]:%d", one.Host, one.Port),
				User:               one.User,
				Password:           one.Password,
				PrivateKey:         one.PrivateKey,
				PrivateKeyPassword: one.PrivateKeyPassword,
			}
			err = sshterm.StartTerm(&arg, nil, time.Second*10)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
var query string

func init() {
	rootCmd.AddCommand(assetCmd)
	assetCmd.Flags().StringVarP(&query, "query", "q", "", "机器模糊搜索关键字")
}

func displayMachineList(q string, webSocket net.Listener) {
	list, err := model.MachineList(db, q)
	if err != nil {
		log.Fatal(err)
	}
	var data [][]string
	for _, row := range list {
		data = append(data, []string{fmt.Sprintf("%d", row.Id), row.Name, row.Protocol, row.Host, fmt.Sprintf("%d", row.Port), row.User})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"ID", "名称", "协议", "HOST", "端口", "用户"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
