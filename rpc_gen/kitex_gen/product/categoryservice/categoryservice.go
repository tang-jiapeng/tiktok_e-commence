// Code generated by Kitex v0.9.1. DO NOT EDIT.

package categoryservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"SelectCategory": kitex.NewMethodInfo(
		selectCategoryHandler,
		newSelectCategoryArgs,
		newSelectCategoryResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	categoryServiceServiceInfo                = NewServiceInfo()
	categoryServiceServiceInfoForClient       = NewServiceInfoForClient()
	categoryServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return categoryServiceServiceInfo
}

// for client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return categoryServiceServiceInfoForStreamClient
}

// for stream client
func serviceInfoForClient() *kitex.ServiceInfo {
	return categoryServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "CategoryService"
	handlerType := (*product.CategoryService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "product",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.9.1",
		Extra:           extra,
	}
	return svcInfo
}

func selectCategoryHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.CategorySelectReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.CategoryService).SelectCategory(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *SelectCategoryArgs:
		success, err := handler.(product.CategoryService).SelectCategory(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*SelectCategoryResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newSelectCategoryArgs() interface{} {
	return &SelectCategoryArgs{}
}

func newSelectCategoryResult() interface{} {
	return &SelectCategoryResult{}
}

type SelectCategoryArgs struct {
	Req *product.CategorySelectReq
}

func (p *SelectCategoryArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.CategorySelectReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *SelectCategoryArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *SelectCategoryArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *SelectCategoryArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *SelectCategoryArgs) Unmarshal(in []byte) error {
	msg := new(product.CategorySelectReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var SelectCategoryArgs_Req_DEFAULT *product.CategorySelectReq

func (p *SelectCategoryArgs) GetReq() *product.CategorySelectReq {
	if !p.IsSetReq() {
		return SelectCategoryArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *SelectCategoryArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *SelectCategoryArgs) GetFirstArgument() interface{} {
	return p.Req
}

type SelectCategoryResult struct {
	Success *product.CategorySelectResp
}

var SelectCategoryResult_Success_DEFAULT *product.CategorySelectResp

func (p *SelectCategoryResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.CategorySelectResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *SelectCategoryResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *SelectCategoryResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *SelectCategoryResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *SelectCategoryResult) Unmarshal(in []byte) error {
	msg := new(product.CategorySelectResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *SelectCategoryResult) GetSuccess() *product.CategorySelectResp {
	if !p.IsSetSuccess() {
		return SelectCategoryResult_Success_DEFAULT
	}
	return p.Success
}

func (p *SelectCategoryResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.CategorySelectResp)
}

func (p *SelectCategoryResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *SelectCategoryResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) SelectCategory(ctx context.Context, Req *product.CategorySelectReq) (r *product.CategorySelectResp, err error) {
	var _args SelectCategoryArgs
	_args.Req = Req
	var _result SelectCategoryResult
	if err = p.c.Call(ctx, "SelectCategory", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
