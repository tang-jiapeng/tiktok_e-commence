package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type LockProductQuantityService struct {
	ctx context.Context
} // NewLockProductQuantityService new LockProductQuantityService
func NewLockProductQuantityService(ctx context.Context) *LockProductQuantityService {
	return &LockProductQuantityService{ctx: ctx}
}

// Run create note info
func (s *LockProductQuantityService) Run(req *product.ProductLockQuantityRequest) (resp *product.ProductLockQuantityResponse, err error) {
	originProducts := req.Products
	var ids []int64 = make([]int64, 0)
	var productQuantityMap map[int64]int64 = make(map[int64]int64)
	for _, pro := range originProducts {
		ids = append(ids, pro.Id)
		productQuantityMap[pro.Id] = pro.Quantity
	}
	productList, err := model.SelectProductList(mysql.DB, context.Background(), ids)
	//确定当前库存是否足够
	var canLock bool = true
	for _, pro := range productList {
		//如果真实库存小于下单的数量，则库存锁定失败
		if pro.RealStock < productQuantityMap[pro.ProductId] {
			canLock = false
			break
		}
	}
	//如果库存锁定失败，则返回失败信息
	if !canLock {
		resp = &product.ProductLockQuantityResponse{
			StatusCode: 2022,
			StatusMsg:  constant.GetMsg(2022),
		}
		return
	}
	//如果库存锁定成功，则更新库存信息
	err = model.UpdateLockStock(mysql.DB, context.Background(), productQuantityMap)
	resp = &product.ProductLockQuantityResponse{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
