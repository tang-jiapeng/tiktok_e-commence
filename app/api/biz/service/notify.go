package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	// payment "tiktok_e-commerce/api/hertz_gen/api/payment"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/go-pay/gopay/alipay"
)

type NotifyService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewNotifyService(Context context.Context, RequestContext *app.RequestContext) *NotifyService {
	return &NotifyService{RequestContext: RequestContext, Context: Context}
}

func (h *NotifyService) Run(RequestContext *app.RequestContext) (err error) {
	AlipayPublicContentPath, _ := os.Open(os.Getenv("ALIPAY_PUBLIC_CONTENT_PATH"))
	AlipayPublicContent, _ := io.ReadAll(AlipayPublicContentPath)

	// 解析异步通知的参数
	request, err := convertToHTTPRequest(&RequestContext.Request)
	if err != nil {
		hlog.Error("支付宝支付通知解析失败, 错误信息: %s", err)
		return
	}
	notifyReq, err := alipay.ParseNotifyToBodyMap(request)
	if err != nil {
		hlog.Error("支付宝支付通知解析失败, 错误信息: %s", err)
		return
	}
	hlog.Info("支付宝支付通知解析成功, 通知参数: %s", notifyReq)
	tradeStatus := notifyReq.GetString("trade_status")
	if tradeStatus != "TRADE_SUCCESS" {
		hlog.CtxErrorf(h.Context, "trade_status is not TRADE_SUCCESS")
	} else {
		//todo 支付成功逻辑

	}
	//// 支付宝异步通知验签（公钥模式）
	//ok, err = alipay.VerifySign(aliPayPublicKey, notifyReq)

	// 支付宝异步通知验签（公钥证书模式）
	_, err = alipay.VerifySignWithCert(AlipayPublicContent, notifyReq)
	if err != nil {
		hlog.CtxErrorf(h.Context, err.Error())
	}

	// 如果需要，可将 BodyMap 内数据，Unmarshal 到指定结构体指针 ptr
	//err = notifyReq.Unmarshal(ptr)

	// ====异步通知，返回支付宝平台的信息====
	// 文档：https://opendocs.alipay.com/open/203/105286
	// 程序执行完后必须打印输出“success”（不包含引号）。如果商户反馈给支付宝的字符不是success这7个字符，支付宝服务器会不断重发通知，直到超过24小时22分钟。一般情况下，25小时以内完成8次通知（通知的间隔频率一般是：4m,10m,10m,1h,2h,6h,15h）

	// 此写法是 gin 框架返回支付宝的写法
	//c.String(http.StatusOK, "%s", "success")

	// 此写法是 echo 框架返回支付宝的写法
	RequestContext.String(http.StatusOK, "success")
	return
}

// convertToHTTPRequest 将 Hertz 的 *protocol.Request 转换为标准库的 *http.Request
func convertToHTTPRequest(req *protocol.Request) (*http.Request, error) {
	// 读取请求体
	body := io.NopCloser(bytes.NewReader(req.Body()))

	// 创建 *http.Request
	httpReq, err := http.NewRequest(
		string(req.Method()),
		req.URI().String(),
		body,
	)
	if err != nil {
		return nil, err
	}

	// 复制 Header
	req.Header.VisitAll(func(key, value []byte) {
		httpReq.Header.Set(string(key), string(value))
	})

	return httpReq, nil
}
