package rpc

import (
	"context"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/felix/util"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) FetchAssetSshConfig(ctx context.Context, arg *pb.ReqAssetUser) (*pb.ResSshConnCfg, error) {
	row := new(model.Machine)
	err := sv.db.First(row, arg.AssetId).Error
	if err != nil {
		return nil, err
	}

	res := pb.ResSshConnCfg{
		Uuid: util.RandString(12),
		AssetConn: &pb.SshConn{
			Addr:               row.SshAddr(),
			User:               row.User,
			Password:           row.Password,
			PrivateKey:         row.PrivateKey,
			PrivateKeyPassword: row.PrivateKeyPassword,
		},
		ProxyConn: nil, // ssh 跳板服务器
	}

	return &res, nil

}
