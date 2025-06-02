package domainservice

import (
	"context"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/dal/dao"
	"github/lhh-gh/go-mall/logic/do"
)

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

// GetDemos 配置GORM时的演示方法
func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.DemoOrder, 0, len(demos))
	// 后面会介绍工具, Model到Domain Object 可以一键转换
	for _, demo := range demos {
		demoOrders = append(demoOrders, &do.DemoOrder{
			Id:           demo.Id,
			UserId:       demo.UserId,
			BillMoney:    demo.BillMoney,
			OrderNo:      demo.OrderNo,
			OrderGoodsId: demo.OrderGoodsId,
			State:        demo.State,
			PaidAt:       demo.PaidAt,
			CreatedAt:    demo.CreatedAt,
			UpdatedAt:    demo.UpdatedAt,
		})
	}

	return demoOrders, nil
}
