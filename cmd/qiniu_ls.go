package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// jsonCmd represents the json command
var qiniuListCmd = &cobra.Command{
	Use:   "qnl",
	Short: "显示全部七牛文件里表",
	Long:  `felix qnl 列出自己私有bucket中的文件,当作文件中转站`,

	Run: func(cmd *cobra.Command, args []string) {
		mac := qbox.NewMac(cfg.QiniuAk, cfg.QiniuSk)
		storageCfg := storage.Config{
			// 是否使用https域名进行资源管理
			UseHTTPS: false,
		}
		// 指定空间所在的区域，如果不指定将自动探测
		// 如果没有特殊需求，默认不需要指定
		//storageCfg.Zone=&storage.ZoneHuabei

		bucketManager := storage.NewBucketManager(mac, &storageCfg)
		//https://developer.qiniu.com/kodo/1238/go#rs-list
		marker := ""
		data := [][]string{}
		for {
			entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(cfg.QiniuBucket, "", "", marker, 1000)
			if err != nil {
				fmt.Println("list error,", err)
				break
			}
			//print entries
			for _, ret := range entries {
				t := time.Unix(ret.PutTime/10000000, 0).Format("06-01-02T15:04:05")
				url := fmt.Sprintf(`%s/%s`, cfg.QiniuBucketEndPoint, ret.Key)
				size := fmt.Sprintf("%dKB", ret.Fsize>>10)
				data = append(data, []string{ret.Key, url, size, t})

			}
			if hasNext {
				marker = nextMarker
			} else {
				//list end
				break
			}
		}
		showTable(data, []string{"KEY", "URL", "Size", "Time"})

	},
}

func init() {
	rootCmd.AddCommand(qiniuListCmd)

}

func showTable(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
