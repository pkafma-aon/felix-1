package cmd

import (
	"errors"
	"fmt"
	"github.com/bytegang/felix/util"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
	"log"
)

// qiniuFetchCmd represents the json command
var qiniuFetchCmd = &cobra.Command{
	Use:   "qnf",
	Short: "七牛远程下载",
	Long: `felix qnu <远程资源URL>
让七牛云下载好URL的文件,让后在从七牛云下载这个文件,
主要解决网络问题
`,
	Example: "felix qnf -e 5  https://download.jetbrains.com/go/goland-2021.2.1.tar.gz  远程下载文件 五天之后自动删除",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("第一个参数书必须是文件路径或则文件夹路径 eg: felix qnu <PATH_OF_FILE_OR_DIRECTORY>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		remoteResourceURL := args[0]

		mac := qbox.NewMac(cfg.QiniuAk, cfg.QiniuSk)
		storageCfg := storage.Config{}
		// 空间对应的机房
		storageCfg.Zone = &storage.ZoneHuanan
		// 是否使用https域名
		storageCfg.UseHTTPS = false
		// 上传是否使用CDN上传加速
		storageCfg.UseCdnDomains = false

		bmgr := storage.NewBucketManager(mac, &storageCfg)

		// 构建表单上传的对象
		key := util.RandString(8)
		ret, err := bmgr.Fetch(remoteResourceURL, cfg.QiniuBucket, key)
		if err != nil {
			log.Fatalln(err)
		}
		err = bmgr.DeleteAfterDays(cfg.QiniuBucket, ret.Key, deleteAfterDays)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("云端下载地址: wget '%s/%s'   Hash:%s\r\n", cfg.QiniuBucketEndPoint, ret.Key, ret.Hash)

	},
}

var deleteAfterDays = 1

func init() {
	rootCmd.AddCommand(qiniuFetchCmd)
	qiniuFetchCmd.Flags().IntVarP(&deleteAfterDays, "expire", "e", 1, "远程下载过期时间,之后自动被删除")

}
