package rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bytegang/felix/model"
	"github.com/bytegang/pb"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func (sv SvrRpc) AuthPk(ctx context.Context, in *pb.ReqAuthPublicKey) (*pb.User, error) {
	user, err := sv._authUserPublicKey(in.User.Account, in.PublicKey)
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

func (sv *SvrRpc) _authUserPublicKey(account string, key []byte) (user *model.User, err error) {
	inAuthorizedKey, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return nil, err
	}
	//fetch user
	user = new(model.User)
	err = sv.db.Where("account = ?", account).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user is not found")
	}
	if err != nil {
		return nil, err
	}

	if len(user.PublicKey) == 0 {
		return nil, errors.New("user's github public key is not fetched yet")
	}
	//parse db public keys
	authorizedKeysBytes := []byte(user.PublicKey)
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing key %v", err)
		}
		authorizedKeysBytes = rest

		if bytes.Equal(inAuthorizedKey.Marshal(), pubKey.Marshal()) {
			return user, nil
		}
	}
	return nil, errors.New("no github public key matched")
}
