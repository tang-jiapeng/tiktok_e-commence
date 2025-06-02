package alipay

import (
	"context"
	"io"
	"os"

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
	notifyUrl = "http://example.com/payment/charge"
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

//Todo 支付宝支付通知回调
