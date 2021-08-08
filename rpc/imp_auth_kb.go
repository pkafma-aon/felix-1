package rpc

import (
	"context"
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
)

func (sv SvrRpc) AuthKb(ctx context.Context, arg *pb.ReqSshUser) (*pb.UserKb, error) {
	row := new(model.User)
	err := sv.db.First(row, "`account` = ?", arg.Account).Error
	if err != nil {
		return nil, err
	}
	u := pb.User{
		Id:      fmt.Sprintf("%d", row.Id),
		Account: row.Account,
		Name:    row.Name,
		Email:   row.Email,
		Phone:   row.Phone,
		Role:    pb.UserRole_Admin,
	}

	if r, ok := pb.UserRole_value[row.Role]; ok {
		u.Role = pb.UserRole(r)
	}
	//这里可以定义成发送手机验证码
	res := pb.UserKb{
		User:        row.Name,
		Instruction: "",
		Questions:   []string{"1+2=?\r\n", "你的邮箱是多少?\r\n"},
		Answers:     []string{"3", row.Email},
		Echos:       []bool{true, true},
		ResUser:     &u,
	}

	return &res, nil

}
