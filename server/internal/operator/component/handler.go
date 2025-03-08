package component

import (
	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/pkg/errors"
)

func (o *Operator) Create() (err error) {
	if err = o.checkComponent(); err != nil {
		err = errors.Wrapf(err, "检查组件参数失败")
		return
	}

	_, err = o.ComponentStore.Create(o.Component)
	if err != nil {
		err = errors.Wrapf(err, "新增组件 %s 失败", o.Component.Name)
		return
	}

	return
}

func (o *Operator) Get(id uint64) (component *models.Component, err error) {
	component, err = o.ComponentStore.Get(id)
	if err != nil {
		err = errors.Wrapf(err, "获取组件 %d 失败", id)
		return
	}

	return
}

func (o *Operator) List() (total int64, components []*models.Component, err error) {
	// 传递分页参数到存储层
	total, components, err = o.ComponentStore.List(o.QueryParma)
	if err != nil {
		err = errors.Wrapf(err, "获取组件列表失败")
		return
	}

	return
}

func (o *Operator) Delete() (err error) {
	if o.Component == nil {
		err = errors.New("组件不能为空")
		return
	}

	err = o.ComponentStore.Delete(o.Component)
	if err != nil {
		err = errors.Wrapf(err, "删除组件 %d 失败", o.Component.Id)
		return
	}

	return
}
