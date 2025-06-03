package alipay

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

var (
	appid      string
	notifyUrl  string
	privateKey string

	// 订单标题
	subject = "tiktok_e-commerce-order"
	//产品码 沙箱环境仅支持value = FAST_INSTANT_TRADE_PAY
	productCode = "FAST_INSTANT_TRADE_PAY"

	// 支付宝公钥证书
	AlipayPublicContent []byte

	// 支付宝根证书
	AlipayRootContent []byte

	// 应用公钥证书
	AppPublicContent []byte

	Client *alipay.Client
)

func Init() {
	AlipayPublicContentPath, _ := os.Open(os.Getenv("ALIPAY_PUBLIC_CONTENT_PATH"))
	AlipayRootContentPath, _ := os.Open(os.Getenv("ALIPAY_ROOT_CONTENT_PATH"))
	AppPublicContentPath, _ := os.Open(os.Getenv("APP_PUBLIC_CONTENT_PATH"))
	defer AlipayPublicContentPath.Close()
	defer AlipayRootContentPath.Close()
	defer AppPublicContentPath.Close()

	appid = os.Getenv("APPID")
	notifyUrl = os.Getenv("NOTIFY_URL")
	privateKey = os.Getenv("PRIVATE_KEY")

	AlipayPublicContent, _ = io.ReadAll(AlipayPublicContentPath)
	AlipayRootContent, _ = io.ReadAll(AlipayRootContentPath)
	AppPublicContent, _ = io.ReadAll(AppPublicContentPath)

	// 初始化支付宝客户端
	// appid：应用ID
	// privateKey：应用私钥，支持PKCS1和PKCS8
	// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
	var err error
	Client, err = alipay.NewClient(appid, privateKey, false)
	if err != nil {
		klog.CtxErrorf(context.Background(), "初始化支付宝客户端失败, 错误信息: %s", err)
		return
	}

	// 打开Debug开关，输出日志，默认关闭
	Client.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	Client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).  // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2). // 设置签名类型，不设置默认 RSA2
							SetNotifyUrl(notifyUrl)   // 设置异步通知URL
	//SetReturnUrl("https://www.fmm.ink"). // 设置返回URL
	// 设置biz_content加密KEY，设置此参数默认开启加密（目前不可用，设置后会报错）
	//client.SetAESKey("1234567890123456")

	// 自动同步验签（只支持证书模式）
	// 传入 alipayPublicCert.crt 内容
	klog.CtxInfof(context.Background(), "支付宝自动同步验签,%+v", AlipayPublicContent)
	Client.AutoVerifySign(AlipayPublicContent)

	// 公钥证书模式，需要传入证书，以下两种方式二选一
	// 证书路径
	//err = client.SetCertSnByPath()
	// 证书内容
	err = Client.SetCertSnByContent(AppPublicContent, AlipayRootContent, AlipayPublicContent)

	if err != nil {
		klog.CtxErrorf(context.Background(), "设置支付宝证书失败, 错误信息: %s", err)
	}
}

func Pay(ctx context.Context, orderId int64, totalAmount float32) (result string, err error) {
	// 构建支付请求参数
	bodyMap := make(gopay.BodyMap)
	bodyMap.Set("out_trade_no", orderId)
	bodyMap.Set("total_amount", totalAmount)
	bodyMap.Set("subject", subject)
	bodyMap.Set("product_code", productCode)
	paymentUrl, err := Client.TradePagePay(ctx, bodyMap)
	if err != nil {
		klog.CtxErrorf(ctx, "支付宝支付失败, 错误信息: %s", err)
		return
	}
	// 跳转到支付页面
	return paymentUrl, nil
}

// 支付宝支付通知回调

func PayNotify(ctx *context.Context , notifyBody string , c *app.RequestContext) {
	// 解析异步通知的参数
	request , err := convertToHTTPRequest(&c.Request) 
	if err != nil {
		klog.Error("支付宝支付通知解析失败, 错误信息: %s", err)
		return
	}
	notifyReq , err := alipay.ParseNotifyToBodyMap(request)
	if err != nil {
		klog.Error("支付宝支付通知解析失败, 错误信息: %s", err)
		return
	}
	//或
	// value：url.Values
	//
	//notifyReq, err = alipay.ParseNotifyByURLValues()
	//if err != nil {
	//	klog.Error("支付宝支付通知解析失败, 错误信息: %s", err)
	//	return
	//}

	//// 支付宝异步通知验签（公钥模式）
	//ok, err = alipay.VerifySign(aliPayPublicKey, notifyReq)

	// 支付宝异步通知验签（公钥证书模式）
	_, err = alipay.VerifySignWithCert(AlipayPublicContent, notifyReq)
	if err != nil {
		klog.Error("支付宝支付通知验签失败, 错误信息: %s", err)
		return
	}

	// 如果需要，可将 BodyMap 内数据，Unmarshal 到指定结构体指针 ptr
	//err = notifyReq.Unmarshal(ptr)

	// ====异步通知，返回支付宝平台的信息====
	// 文档：https://opendocs.alipay.com/open/203/105286
	// 程序执行完后必须打印输出“success”（不包含引号）。如果商户反馈给支付宝的字符不是success这7个字符，支付宝服务器会不断重发通知，直到超过24小时22分钟。一般情况下，25小时以内完成8次通知（通知的间隔频率一般是：4m,10m,10m,1h,2h,6h,15h）

	// 此写法是 gin 框架返回支付宝的写法
	//c.String(http.StatusOK, "%s", "success")

	// 此写法是 echo 框架返回支付宝的写法
	c.String(http.StatusOK, "success")
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