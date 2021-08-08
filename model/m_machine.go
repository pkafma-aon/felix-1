package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytegang/felix/util"
	"gorm.io/gorm"
	"time"
)

type Machine struct {
	BaseModel
	Protocol           string `json:"protocol"` // ssh rdp vnc
	Name               string `json:"name"`
	Host               string `json:"host"`
	Ip                 string `json:"ip"`
	Port               int    `json:"port"`
	Remark             string `json:"remark"`
	User               string `json:"user"`
	Password           string `json:"password"`
	PrivateKey         string `json:"private_key"`
	PrivateKeyPassword string `json:"private_key_password"`
	WebSshURL          string `gorm:"-" json:"web_ssh_url"`
}

func (m Machine) ID() string {
	return fmt.Sprintf("%d", m.Id)
}

func (m Machine) SshAddr() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func MachineList(db *gorm.DB, q string) (list []Machine, err error) {
	tx := db.Model(new(Machine))
	if q != "" {
		tx.Where("name LIKE ? OR host LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	err = db.Find(&list).Error
	return
}

type machinePayLoad struct {
	Id        uint  `json:"i"`
	NotBefore int64 `json:"f"`
	NotAfter  int64 `json:"t"`
}

func (receiver machinePayLoad) Valid() error {
	from := time.Unix(receiver.NotBefore, 0)
	to := time.Unix(receiver.NotAfter, 0)
	if time.Now().After(to) || time.Now().Before(from) {
		return errors.New("token is expired")
	}
	return nil
}

func (m *Machine) GenerateToken(key []byte, expireIn time.Duration) (string, error) {
	now := time.Now()
	data, err := json.Marshal(machinePayLoad{
		Id:        m.Id,
		NotBefore: now.Unix(),
		NotAfter:  now.Add(expireIn).Unix(),
	})
	if err != nil {
		return "", err
	}
	return util.AesEncrypt(data, key)
}
func MachineFrom(db *gorm.DB, secret string, key []byte) (ins *Machine, err error) {
	bytes, err := util.AesDecrypt(secret, key)
	if err != nil {
		return nil, err
	}
	arg := new(machinePayLoad)
	err = json.Unmarshal(bytes, arg)
	if err != nil {
		return nil, err
	}
	err = arg.Valid()
	if err != nil {
		return nil, err
	}

	ins = new(Machine)
	err = db.First(ins, arg.Id).Error
	return ins, err
}
