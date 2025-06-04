package appservice

import (
	"context"
	"github/lhh-gh/go-mall/api/reply"
	"github/lhh-gh/go-mall/api/request"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/util"
	"github/lhh-gh/go-mall/dal/cache"
	"github/lhh-gh/go-mall/logic/do"
	"github/lhh-gh/go-mall/logic/domainservice"
)

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}
func (das *DemoAppSvc) GetDemoIdentities() ([]int64, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]int64, 0, len(demos))

	for _, demo := range demos {
		identities = append(identities, demo.Id)
	}
	return identities, nil
}
func (das *DemoAppSvc) CreateDemoOrder(orderRequest *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	err := util.CopyProperties(demoOrderDo, orderRequest)
	if err != nil {
		errcode.Wrap("请求转换成demoOrderDo失败", err)
		return nil, err
	}
	demoOrderDo, err = das.demoDomainSvc.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知
	// TODO2 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知

	// 设置缓存和读取, 测试项目中缓存的使用, 没有其他任何意义
	cache.SetDemoOrder(das.ctx, demoOrderDo)
	cacheData, _ := cache.GetDemoOrder(das.ctx, demoOrderDo.OrderNo)
	logger.New(das.ctx).Info("redis data", "data", cacheData)
	replyDemoOrder := new(reply.DemoOrder)
	err = util.CopyProperties(replyDemoOrder, demoOrderDo)
	if err != nil {
		errcode.Wrap("demoOrderDo转换成replyDemoOrder失败", err)
		return nil, err
	}

	return replyDemoOrder, err
}
