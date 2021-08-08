package rpc

import (
	"context"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
	"strconv"
)

func (sv SvrRpc) Guacamole(ctx context.Context, arg *pb.ReqToken) (*pb.ResGuacamole, error) {
	machine, err := model.MachineFrom(sv.db, arg.Token, []byte(sv.cfg.Secret))
	if err != nil {
		return nil, err
	}

	return &pb.ResGuacamole{
		Protocol: machine.Protocol,
		Host:     machine.Host,
		Port:     strconv.Itoa(machine.Port),
		User:     machine.User,
		Password: machine.Password,
		Width:    1800,
		Height:   900,
	}, nil

}
