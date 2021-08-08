package cmd

import (
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/rpc"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
	"strings"
)

// qiniuRmCmd represents the json command
var qiniuRmCmd = &cobra.Command{
	Use:   "qnrm",
	Short: "七牛删除文件",
	Long:  `felix qnrm <KEY>... 需要删除的七牛远程文件, 需要配合 felix qnl 来查询 KEY`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("第一个参数书必须是KEY")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		keys := args
		mac := qbox.NewMac(cfg.QiniuAk, cfg.QiniuSk)
		storageCfg := storage.Config{
			// 是否使用https域名进行资源管理
			UseHTTPS: false,
		}
		// 指定空间所在的区域，如果不指定将自动探测
		// 如果没有特殊需求，默认不需要指定
		//storageCfg.Zone=&storage.ZoneHuabei

		bucketManager := storage.NewBucketManager(mac, &storageCfg)
		//每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送

		deleteOps := make([]string, 0, len(keys))
		for _, key := range keys {
			deleteOps = append(deleteOps, storage.URIDelete(cfg.QiniuBucket, key))
		}

		rets, err := bucketManager.Batch(deleteOps)
		if err != nil {
			// 遇到错误
			if _, ok := err.(*rpc.ErrorInfo); ok {
				for _, ret := range rets {
					// 200 为成功
					if ret.Code != 200 {
						fmt.Printf("%s\n", ret.Data.Error)
					} else {
						fmt.Printf("文件:%s删除失败", strings.Join(keys, ", "))
					}
				}
			} else {
				fmt.Printf("batch error, %s", err)
			}
		} else {
			// 完全成功
			for _, ret := range rets {
				// 200 为成功
				fmt.Printf("文件:%s删除成功", ret.Data.Hash)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(qiniuRmCmd)
}
