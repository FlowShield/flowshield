package service

import (
	"github.com/flowshield/flowshield/provider/internal/dao/provider/dao"
	"github.com/flowshield/flowshield/provider/internal/dao/provider/model"
)

func ListProvider(param *model.Provider) (model.Providers, error) {
	list, err := dao.NewProvider().ListProvider(param)
	if err != nil {
		return nil, err
	}
	return list, err
}

func AddProvider(item *model.Provider) error {
	data, err := dao.NewProvider().GetProviderByUuid(item.Uuid)
	if err != nil {
		return err
	}
	if data.ID > 0 {
		item.ID = data.ID
		item.CreatedAt = data.CreatedAt
		err = dao.NewProvider().EditProvider(item)
		if err != nil {
			return err
		}
	} else {
		err = dao.NewProvider().AddProvider(item)
		if err != nil {
			return err
		}
	}
	return nil
}
