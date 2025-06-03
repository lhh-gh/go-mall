package dao

import (
	"context"
	"github/lhh-gh/go-mall/comon/util"
	"github/lhh-gh/go-mall/dal/model"
	"github/lhh-gh/go-mall/logic/do"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.DemoOrder, err error) {

	err = DB().WithContext(demo.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}

	return demos, err
}

func (demo *DemoDao) CreateDemoOrder(demoOrder *do.DemoOrder) (*model.DemoOrder, error) {
	model := new(model.DemoOrder)
	err := util.CopyProperties(model, demoOrder)
	if err != nil {
		return nil, err
	}
	err = DB().WithContext(demo.ctx).
		Create(model).Error
	return model, err
}
