package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var nesDir string
var ymlFile string

var nesScanCmd = &cobra.Command{
	Use:   "nes",
	Short: "扫描.nes文件生成静态文件和json index",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		yamlFileW, err := os.Create("/code/tech.mojotv.cn/_data/game.yml")
		if err != nil {
			log.Fatal(err)
		}
		nesDir = "/code/tech.mojotv.cn/nes"
		filepath.Walk(nesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".nes") {
				thisDir := filepath.Dir(path)
				fileName := strings.ToLower(info.Name())
				fileName = strings.Replace(fileName, "final_fight", "快打旋风 ", -1)
				fileName = strings.Replace(fileName, "contra", "魂斗罗", -1)
				fileName = strings.Replace(fileName, "donkey_kong", "金刚", -1)
				fileName = strings.Replace(fileName, "double_dragon", "双截龙", -1)
				fileName = strings.Replace(fileName, "super_mario_bros", "超级玛丽兄弟", -1)
				fileName = strings.Replace(fileName, "dragon_ball", "龙珠", -1)
				fileName = strings.Replace(fileName, "final_fantasy", "最终幻想", -1)
				fileName = strings.Replace(fileName, "transformers", "变形金刚", -1)
				fileName = strings.Replace(fileName, "dragon_warrior", "龙战士", -1)
				fileName = strings.Replace(fileName, "arkanoid", "快打砖块", -1)
				fileName = strings.Replace(fileName, "indiana_jones", "印第安纳琼斯", -1)
				fileName = strings.Replace(fileName, "balloon_fight", "气球大战", -1)
				fileName = strings.Replace(fileName, "dale_rescue_rangers", "松鼠大作战", -1)
				fileName = strings.Replace(fileName, "battle_city", "坦克大战", -1)
				fileName = strings.Replace(fileName, "bomberman", "炸弹人", -1)
				fileName = strings.Replace(fileName, "commando", "古巴兄弟", -1)
				fileName = strings.Replace(fileName, "doraemon", "哆啦A梦", -1)
				fileName = strings.Replace(fileName, "excitebike", "越野摩托车", -1)
				fileName = strings.Replace(fileName, "hudsons_adventure_island", "冒险岛", -1)
				fileName = strings.Replace(fileName, "mario_bros", "马里奥兄弟", -1)
				fileName = strings.Replace(fileName, "thunderbirds", "雷鸟战机", -1)
				fileName = strings.Replace(fileName, "!", "", -1)
				fileName = strings.Replace(fileName, "'", "", -1)
				fileName = strings.Replace(fileName, " ", "_", -1)
				fileName = strings.Replace(fileName, "-", "_", -1)
				fileName = strings.Replace(fileName, "___", "_", -1)
				fileName = strings.Replace(fileName, "__", "_", -1)
				fileName = strings.Replace(fileName, "_.", ".", -1)

				//t := strings.TrimRight(fileName,".nes")
				//zh, err := googleTranslateEn2ZH(t)
				//if err != nil {
				//	log.Println(err)
				//	fileName = t+".nes"
				//}else {
				//	fileName = zh+".nes"
				//
				//}
				//time.Sleep(time.Millisecond*100)
				//
				//if matched, _ := regexp.MatchString(`fc\d{4}`,path);matched{
				//	os.Remove(path)
				//	return nil
				//}

				nPath := filepath.Join(thisDir, fileName)
				err = os.Rename(path, nPath)
				if err != nil {
					log.Println("重命名失败", path, nPath, err)
				}
				fn := fmt.Sprintf(`- "%s"`, strings.ReplaceAll(nPath, nesDir+"/", ""))
				fn = strings.ReplaceAll(fn, `\`, "/") + "\n"

				fmt.Fprintf(yamlFileW, fn)
			}
			return nil
		})
		yamlFileW.Close()

	},
}

func init() {
	rootCmd.AddCommand(nesScanCmd)
	nesScanCmd.Flags().StringVarP(&nesDir, "src", "s", "/code/tech.mojotv.cn", "nes文件夹")
	nesScanCmd.Flags().StringVarP(&ymlFile, "file", "f", "", "yaml索引文件名")
}
