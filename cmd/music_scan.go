package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var scanMusicCmd = &cobra.Command{
	Use:     "scan-music",
	Short:   "scan music files 生产音乐index.json",
	Long:    `felix scan-music <PathMusicDir> <PathOutPutJson>  eg felix scan-music /data/music /data/music/index.json`,
	Example: `felix scan-music <PathMusicDir> <PathOutPutJson>`,

	Run: func(cmd *cobra.Command, args []string) {
		//https://github.com/dhowden/tag music tag
		//baseDir := "/code/tech.mojotv.cn"
		baseDir := args[0]
		list, err := scanMusicDir(baseDir)
		if err != nil {
			log.Fatalln(err)
		}

		bs, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
			return
		}
		jsonPath := args[1]
		err = os.WriteFile(jsonPath, bs, 0644)
		if err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(scanMusicCmd)
}

type music struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
	Cover  string `json:"cover"`
	Path   string `json:"-"`
}

func b64UriImage(tag tag.Metadata) string {
	return "/assets/image/logo00.png"
	pic := tag.Picture()
	if pic == nil {
		return ""
	}

	return fmt.Sprintf("data:%s;%s", pic.MIMEType, base64.StdEncoding.EncodeToString(pic.Data))
}

func scanMusicDir(baseDir string) ([]music, error) {
	list := []music{}
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		mFile, err := os.Open(path)
		if err != nil {
			return err
		}
		m, err := tag.ReadFrom(mFile)
		if err != nil {
			return err
		}

		one := music{
			Name:   m.Title(),
			Artist: m.Artist(),
			Url:    strings.ReplaceAll(path, baseDir, ""),
			Cover:  b64UriImage(m),
			Path:   path,
		}

		list = append(list, one)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
