module github.com/bytegang/felix

go 1.16

require (
	github.com/PuerkitoBio/goquery v1.7.1
	github.com/abadojack/whatlanggo v1.0.1
	github.com/bytegang/pb v0.0.0-20210618014102-a6c30fc28e32
	github.com/bytegang/sshd v0.0.0-00010101000000-000000000000
	github.com/dhowden/tag v0.0.0-20201120070457-d52dcb253c63
	github.com/fatih/color v1.12.0 // indirect
	github.com/gin-gonic/gin v1.7.3
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/hajimehoshi/go-mp3 v0.3.2
	github.com/hajimehoshi/oto v0.6.1
	github.com/mattn/go-runewidth v0.0.13
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/qiniu/go-sdk/v7 v7.9.8
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/grpc v1.38.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)

replace github.com/bytegang/sshd => ../sshd
