package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
	"gorm.io/gorm"
)

func (sv SvrRpc) AuthPw(ctx context.Context, in *pb.ReqAuthPassword) (*pb.User, error) {
	user, err := sv._authUserPassword(in.User.Account, in.Password)
	if err != nil {
		return nil, err
	}
	u := pb.User{
		Id:      fmt.Sprintf("%d", user.Id),
		Name:    user.Name,
		Account: user.Account,
		Email:   user.Email,
		Phone:   user.Phone,
		Role:    pb.UserRole_Reporter,
	}
	if r, ok := pb.UserRole_value[user.Role]; ok {
		u.Role = pb.UserRole(r)
	}

	return &u, nil
}

func (sv *SvrRpc) _authUserPassword(account string, password []byte) (user *model.User, err error) {
	user = new(model.User)
	err = sv.db.Where("account = ?", account).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user is not found")
	}
	if err != nil {
		return nil, err
	}

	if len(user.Password) == 0 {
		return nil, errors.New("user's password is not set yet")
	}
	// just string compare
	// no bcrypt because I just want my frontend show the user's raw password
	if user.Password != string(password) {
		return nil, errors.New("password is not correct")
	}

	return user, err
}
