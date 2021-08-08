package rpc

import (
	"context"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) FetchAsset(ctx context.Context, arg *pb.ReqAssetsQuery) (*pb.AssetList, error) {
	var machines []model.Machine

	tx := sv.db.Model(new(model.Machine))
	if arg.Query != "" {
		tx.Where("`host` LIKE ? OR `name` LIKE ?", "%"+arg.Query+"%", "%"+arg.Query+"%")
	}
	err := tx.Find(&machines).Error
	if err != nil {
		return nil, err
	}
	pbAssets := make([]*pb.Asset, len(machines))
	for i, machine := range machines {
		one := pb.Asset{
			Id:       machine.ID(),
			Hostname: machine.Host,
			Alias:    machine.Name,
			Remark:   machine.Remark,
			ShhAddr:  machine.SshAddr(),
			Ip:       machine.Ip,
		}
		pbAssets[i] = &one
	}

	res := pb.AssetList{
		List: pbAssets,
	}
	return &res, nil

}
