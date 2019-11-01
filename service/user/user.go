package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"rpcx-demo/service/user/model"
)

var AuthMap = make(map[string]string, 10)

type UserService struct {
}

// New 新建一个 服务对象.
func New() *UserService {
	return &UserService{}
}

func (u *UserService) Login(ctx context.Context, req model.AuthRequest, res *model.AuthResponse) error {
	token, ok := IsLogin(req)
	if !ok {
		// res
		*res = model.AuthResponse{Token: token}
		return nil
	} else {
		token, err := checkPwd(req)
		if err != nil {
			return err
		}
		save_into_map(req, token)
		*res = model.AuthResponse{Token: token}
		return nil
	}

}

func (u *UserService) Say(ctx context.Context, req model.SayRequest, res *model.SayResponse) error {
	word := string(req)
	fmt.Println(word)
	*res = model.SayResponse(word)
	return nil

}

// ========= internal func ================
func checkPwd(request model.AuthRequest) (string, error) {
	if request.UserName == "user" {
		return uuid.GenerateUUID()
	}
	return "", errors.New("password is incorrect")
}

func IsLogin(request model.AuthRequest) (string, bool) {
	token, ok := AuthMap[request.UserName]
	if ok {
		return token, true
	}
	return "", false
}

func save_into_map(request model.AuthRequest, token string) {
	AuthMap[request.UserName] = token
}
