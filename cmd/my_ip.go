package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
	"net/http"
)

// jsonCmd represents the json command
var myIpCmd = &cobra.Command{
	Use:     "myip",
	Short:   "获取我的IP信息,内网外网信息",
	Long:    `获取本机的本地IP,外网出口IP`,
	Example: "felix myip",
	Run: func(cmd *cobra.Command, args []string) {
		var data [][]string
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, address := range addrs {
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				data = append(data, []string{ipNet.Network(), ipNet.IP.String()})
			}
		}
		data = append(data, []string{"ExternalIP", getMyExternalIp()})
		showTable(data, []string{"MASK", "IP"})
	},
}

func init() {
	rootCmd.AddCommand(myIpCmd)

}

type ipInfo struct {
	Query string
}

func getMyExternalIp() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip ipInfo
	err = json.Unmarshal(body, &ip)
	if err != nil {
		return err.Error()
	}
	return ip.Query
}
