package cmd

import (
	"fmt"
	"github.com/bytegang/felix/rpc"
	"github.com/bytegang/felix/util"
	"github.com/bytegang/felix/web"
	"github.com/bytegang/sshd/sshweb"
	"github.com/spf13/cobra"
	"log"
	"net"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "打开网页编辑SQLite3数据,和其他的一些操作,windows操作性不兼容代替功能",
	Long: `

`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.AddrWebUi = addr
		//1.start rpc
		lisRpc, err := net.Listen("tcp", cfg.AddrRpc)
		if err != nil {
			log.Fatal(err)
		}
		cfg.AddrRpc = lisRpc.Addr().String()
		go rpc.NewSvrRpc(db, cfg).Run(lisRpc)

		//2. start web sshd
		lisWebSshd, err := net.Listen("tcp", cfg.AddrWebSshd)
		if err != nil {
			log.Fatal(err)
		}
		cfg.AddrWebSshd = lisWebSshd.Addr().String()
		go runWebSshd(lisWebSshd, lisRpc, cfg.AddrGuacamole)
		//3. start felix web admin
		lisUI, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		cfg.AddrWebUi = lisUI.Addr().String()

		//4. open browser
		tab := "/"
		if len(args) > 0 {
			tab += args[0]
		}
		url := "http://" + cfg.AddrWebUi + tab
		util.BrowserOpen(url)
		fmt.Println("open browser: ", url)
		if err := web.RunWebUI(db, cfg, lisUI); err != nil {
			log.Fatal(err)
		}
	},
}

var addr string

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:0", "制定网页服务启动的地址(127.0.0.1:0代表打开网页随机空闲端口,之可以被本机访问,当然你也可以使用 :8080对外网可见)")

}

// 启动webTerminal网页服务
func runWebSshd(lisWeb, lisRpc net.Listener, addrGuacad string) {
	rpc := lisRpc.Addr().String()
	app := sshweb.NewApp(rpc, addrGuacad)
	err := app.Run(lisWeb)
	if err != nil {
		log.Fatal(err)
	}
}
