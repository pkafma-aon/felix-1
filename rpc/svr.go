package rpc

import (
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type SvrRpc struct {
	pb.UnimplementedByteGangsterServer
	db  *gorm.DB
	cfg *model.CfgSchema
}

func NewSvrRpc(db *gorm.DB, cfg *model.CfgSchema) *SvrRpc {
	return &SvrRpc{db: db, cfg: cfg}
}

func (sv *SvrRpc) Run(lis net.Listener) {
	log.Println("gRPC at: ", lis.Addr().String())
	gSvr := grpc.NewServer()
	pb.RegisterByteGangsterServer(gSvr, sv)
	if err := gSvr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
