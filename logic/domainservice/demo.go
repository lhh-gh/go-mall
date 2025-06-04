package domainservice

import (
	"context"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/util"
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
//
//	func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
//		demos, err := dds.DemoDao.GetAllDemos()
//		if err != nil {
//			err = errcode.Wrap("query entity error", err)
//			return nil, err
//		}
//
//		demoOrders := make([]*do.DemoOrder, 0, len(demos))
//		// 通过反射来完成数据复制
//		//
//		//copier 是 Go 语言的一个实现不同结构体之间成员复制的库。
//		//
//		//它可以实现不同结构体的对象到对象的复制、对象到 slice 的复制、slice 到 slice 的复制，支持同名方法和成员、同名成员和方法的复制，还支持 map 到 map 间的复制。后面会介绍工具, Model到Domain Object 可以一键转换
//		for _, demo := range demos {
//			demoOrders = append(demoOrders, &do.DemoOrder{
//				Id:           demo.Id,
//				UserId:       demo.UserId,
//				BillMoney:    demo.BillMoney,
//				OrderNo:      demo.OrderNo,
//				OrderGoodsId: demo.OrderGoodsId,
//				State:        demo.State,
//				PaidAt:       demo.PaidAt,
//				CreatedAt:    demo.CreatedAt,
//				UpdatedAt:    demo.UpdatedAt,
//			})
//		}
//
//		return demoOrders, nil
//	}
//
// GetDemos 配置GORM时的演示方法
func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.DemoOrder, 0, len(demos))
	for _, demo := range demos {
		demoOrder := new(do.DemoOrder)
		util.CopyProperties(demoOrder, demo)
		demoOrders = append(demoOrders, demoOrder)
	}

	return demoOrders, nil
}

func (dds *DemoDomainSvc) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
	// 生成订单号  先随便写个
	demoOrder.OrderNo = "20240627596615375920904456"

	demoOrderModel, err := dds.DemoDao.CreateDemoOrder(demoOrder)

	if err != nil {
		err = errcode.Wrap("创建DemoOrder失败", err)
		return nil, err
	}

	// TODO1: 写订单快照
	// 这里一般要在事务里写订单商品快照表, 这个等后面做需求时再演示
	err = util.CopyProperties(demoOrder, demoOrderModel)
	// 返回领域对象
	return demoOrder, err
}

//func (dds *DemoDomainSvc) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
//
//}
