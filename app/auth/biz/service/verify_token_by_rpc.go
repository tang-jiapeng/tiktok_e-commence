package service

import (
	"context"
	"github.com/casbin/casbin/v2"
	casbinUtil "tiktok_e-commerce/auth/casbin"
	"tiktok_e-commerce/auth/conf"
	"tiktok_e-commerce/auth/utils/jwt"
	"tiktok_e-commerce/auth/utils/redis"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {

	// 检查accessToken是否为空
	if req.AccessToken == "" {
		return &auth.VerifyResp{
			StatusCode: 1004,
			StatusMsg:  constant.GetMsg(1004),
		}, nil
	}

	// 校验access token
	claims, status := jwt.ParseJWT(req.AccessToken)
	if status != jwt.TokenValid {
		return &auth.VerifyResp{
			StatusCode: 1004,
			StatusMsg:  constant.GetMsg(1004),
		}, nil
	}
	userId := int32(claims["userId"].(float64))

	if _, exist := conf.BannedUserList[userId]; exist {
		// 黑名单用户
		return &auth.VerifyResp{
			StatusCode: 1006,
			StatusMsg:  constant.GetMsg(1006),
		}, nil
	}

	role := claims["role"].(string)

	// 校验 redis 中的access token
	savedAccessToken, err := redis.GetVal(s.ctx, redis.GetAccessTokenKey(userId))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if savedAccessToken != req.AccessToken {
		return &auth.VerifyResp{
			StatusCode: 1005,
			StatusMsg:  constant.GetMsg(1005),
		}, nil
	}

	// casbin 鉴权
	hasPermisson := checkPermission(casbinUtil.Enforcer, role, req.Path, req.Method)
	if !hasPermisson {
		return &auth.VerifyResp{
			StatusCode: 1009,
			StatusMsg:  constant.GetMsg(1009),
		}, nil
	}

	return &auth.VerifyResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		UserId:     userId,
	}, nil
}

func checkPermission(enforcer *casbin.Enforcer, role, path, method string) bool {
	ok, err := enforcer.Enforce(role, path, method)
	if err != nil {
		klog.Errorf("权限校验失败: %v", err)
		return false
	}
	return ok
}
