package rpc

import (
	"context"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) SaveLogSshSession(ctx context.Context, arg *pb.ReqSshdData) (*pb.ResStatus, error) {

	res := new(pb.ResStatus)
	return res, nil
}
