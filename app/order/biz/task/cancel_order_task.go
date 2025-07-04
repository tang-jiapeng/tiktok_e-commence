package task

import (
	"context"
	"sync"
	"tiktok_e-commerce/order/biz/dal/mysql"
	"tiktok_e-commerce/order/biz/model"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xxl-job/xxl-job-executor-go"
)

var (
	maxConcurrency = 10
)

// CancelOrderTask 定时取消超时的订单，作为mq延时取消订单的兜底
func CancelOrderTask(ctx context.Context, param *xxl.RunReq) (msg string) {
	// 查询已超时的订单
	now := time.Now()
	orderIdList, err := model.GetOverdueOrder(ctx, mysql.DB, now.Add(-10*time.Minute))
	if err != nil {
		klog.Errorf("定时任务查询超时订单失败: %v", err)
		return "查询超时订单失败" + err.Error()
	}

	// 取消超时订单的支付
	concurrentCancelCharges(ctx, &sync.WaitGroup{}, make(chan struct{}, maxConcurrency), orderIdList)

	affectedRows, err := model.CancelOrderList(ctx, mysql.DB, orderIdList)
	if err != nil {
		klog.Errorf("定时任务取消超时订单失败: %v", err)
		return "取消超时订单失败"
	}
	klog.Infof("定时任务取消超时订单成功，耗时%.1f秒，本次执行的超时订单：%v，成功取消%d个订单", time.Since(now).Seconds(), orderIdList, affectedRows)

	return "success"
}

func concurrentCancelCharges(ctx context.Context, wg *sync.WaitGroup, guard chan struct{}, orderIdList []string) {
	for _, orderId := range orderIdList {
		wg.Add(1)
		guard <- struct{}{}

		go func(orderId string) {
			defer wg.Done()
			defer func() {
				<-guard
			}()
			if err := cancelCharge(ctx, orderId); err != nil {
				klog.Errorf("定时任务取消超时订单的支付失败，订单ID：%s, err： %v", orderId, err)
			}
		}(orderId)
	}
}

func cancelCharge(ctx context.Context, orderId string) error {
	// TODO 取消支付
	return nil
}
