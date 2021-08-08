package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytegang/felix/util"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cobra"
	"path/filepath"
)

// qiniuUploadCmd represents the json command
var qiniuUploadCmd = &cobra.Command{
	Use:   "qnu",
	Short: "七牛上传",
	Long:  `felix qnu <文件/文件夹> 到七牛文件中转站, 如果是文件夹则被zip`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("第一个参数书必须是文件路径或则文件夹路径 eg: felix qnu <PATH_OF_FILE_OR_DIRECTORY>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sourceFileOrDir := args[0]

		putPolicy := storage.PutPolicy{
			Scope: cfg.QiniuBucket,
		}
		mac := qbox.NewMac(cfg.QiniuAk, cfg.QiniuSk)
		upToken := putPolicy.UploadToken(mac)

		storageCfg := storage.Config{}
		// 空间对应的机房
		storageCfg.Zone = &storage.ZoneHuanan
		// 是否使用https域名
		storageCfg.UseHTTPS = false
		// 上传是否使用CDN上传加速
		storageCfg.UseCdnDomains = false

		// 构建表单上传的对象
		formUploader := storage.NewFormUploader(&storageCfg)
		ret := storage.PutRet{}

		// 可选配置
		putExtra := storage.PutExtra{
			Params: map[string]string{
				"x:hostname": util.GetHostName(),
				"x:username": util.GetUserName(),
			},
		}
		fileKeyName := fmt.Sprintf("%s_%s", util.RandomString(8), filepath.Base(sourceFileOrDir))

		//upload a dir
		isDirectory, err := util.IsDirectory(sourceFileOrDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		if isDirectory {
			z := util.NewZippie(sourceFileOrDir)
			defer z.Close()
			zipFilePath, _, err := z.Zip()
			if err != nil {
				fmt.Println("创建压缩文件失败", err)
				return
			}
			sourceFileOrDir = zipFilePath
			fileKeyName += ".zip"
			putExtra.Params["x:zip"] = z.Base
		}

		err = formUploader.PutFile(context.Background(), &ret, upToken, fileKeyName, sourceFileOrDir, &putExtra)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("文件下载地址: wget '%s/%s'   Hash:%s\r\n", cfg.QiniuBucketEndPoint, ret.Key, ret.Hash)

	},
}

func init() {
	rootCmd.AddCommand(qiniuUploadCmd)

}
