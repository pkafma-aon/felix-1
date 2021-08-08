package rpc

import (
	"context"
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/felix/util"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) WebXtermSsh(c context.Context, arg *pb.ReqToken) (*pb.ResSshConnCfg, error) {
	machine, err := model.MachineFrom(sv.db, arg.Token, []byte(sv.cfg.Secret))
	if err != nil {
		return nil, err
	}

	return &pb.ResSshConnCfg{
		Uuid: util.RandomString(12),
		AssetConn: &pb.SshConn{
			Addr:               fmt.Sprintf("%s:%d", machine.Host, machine.Port),
			User:               machine.User,
			Password:           machine.Password,
			PrivateKey:         machine.PrivateKey,
			PrivateKeyPassword: machine.PrivateKeyPassword,
		},
		ProxyConn: nil,
	}, nil

}
