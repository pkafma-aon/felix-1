package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"time"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "修改google-dns",
	Long:  `linux 机器修改dns 使用公司自带的翻墙`,
	Run: func(c *cobra.Command, args []string) {
		body := `
nameserver 10.10.18.10
nameserver 10.10.18.11
`
		if len(args) >= 1 && (args[0] == "g") {
			log.Println("google")
			body = `
nameserver 8.8.8.8
nameserver 8.8.4.4
`
		} else if len(args) >= 1 && (args[0] == "1") {
			log.Println("cloudflare")
			body = `
nameserver 1.1.1.1
nameserver 1.0.0.1
`
		} else {
			log.Println("使用公司DNS")
		}
		fd, err := os.Create("/etc/resolv.conf")
		if err != nil {
			log.Println(err)
			return
		}
		_, err = fd.Write([]byte(body))
		if err != nil {
			log.Println(err)
			return
		}
		fd.Close()
		log.Println("修改DNS结束")
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}

func lookup() {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 2,
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	ip, _ := r.LookupHost(context.Background(), "www.google.com")

	print(ip[0])
}
