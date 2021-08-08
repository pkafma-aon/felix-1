package cmd

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/spf13/cobra"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var isShuffle bool

var musicCmd = &cobra.Command{
	Use:   "music",
	Short: "scan music files",
	Long: `felix scan-music <PathMusicDir> <PathOutPutJson>  eg felix scan-music /data/music /data/music/index.json
编译环境
https://github.com/hajimehoshi/oto

`,
	Example: `felix scan-music <PathMusicDir> <PathOutPutJson>`,

	Run: func(cmd *cobra.Command, args []string) {
		list, err := scanMusicDir("/data/mojoMusic")
		if err != nil {
			log.Fatal(err)
		}
		if isShuffle {
			shuffle(list)
			for i, _ := range list {
				pay(list, i)
			}
			return
		}

		data := [][]string{}
		for i, m := range list {
			data = append(data, []string{fmt.Sprintf("%d", i), m.Name, m.Artist, m.Path})
		}

		showTable(data, []string{"ID", "歌曲", "歌手", "PATH"})

		if len(args) > 0 && args[0] != "" {

			parseInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			for i := int(parseInt); i < len(list); i++ {
				pay(list, i)
			}
		}

	},
}

func pay(list []music, idx int) {
	if idx > len(list) || idx < 0 {
		log.Fatal("无效歌曲ID")
	}
	one := list[idx]
	fmt.Println("播放:  ", idx, "    ", one.Name, "     \t", one.Artist)
	if err := run(one.Path); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(musicCmd)
	musicCmd.Flags().BoolVarP(&isShuffle, "shuffle", "s", false, "随机播放全部歌曲")
}

func run(musicFilePath string) error {
	f, err := os.Open(musicFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	//fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func shuffle(a []music) {
	rand.Seed(time.Now().UnixNano())
	for i := len(a) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
