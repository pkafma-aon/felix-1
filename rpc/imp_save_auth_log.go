package rpc

import (
	"context"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) SaveLogAuth(ctx context.Context, log *pb.ReqAuthLog) (*pb.ResStatus, error) {
	res := new(pb.ResStatus)
	return res, nil
}
