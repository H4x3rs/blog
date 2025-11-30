package service

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/entity"
)

type sSettings struct{}

var Settings = sSettings{}

func (s *sSettings) Get(ctx context.Context, key string) (value string, err error) {
	var setting *entity.Settings
	err = dao.Settings.Ctx(ctx).Where(dao.Settings.Columns().Key, key).Scan(&setting)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", nil
	}
	return setting.Value, nil
}

func (s *sSettings) Set(ctx context.Context, key, value string) (err error) {
	var setting *entity.Settings
	err = dao.Settings.Ctx(ctx).Where(dao.Settings.Columns().Key, key).Scan(&setting)
	if err != nil {
		return
	}
	if setting != nil {
		// 更新时使用结构体方式
		setting.Value = value
		_, err = dao.Settings.Ctx(ctx).Data(setting).Where(dao.Settings.Columns().Key, key).Update()
	} else {
		// 插入时使用结构体
		_, err = dao.Settings.Ctx(ctx).Data(&entity.Settings{
			Key:   key,
			Value: value,
			Type:  "string",
		}).Insert()
	}
	return
}

func (s *sSettings) GetAll(ctx context.Context) (settings map[string]string, err error) {
	var list []*entity.Settings
	err = dao.Settings.Ctx(ctx).Scan(&list)
	if err != nil {
		return
	}
	settings = make(map[string]string)
	for _, item := range list {
		settings[item.Key] = item.Value
	}
	return
}

func (s *sSettings) BatchSet(ctx context.Context, data map[string]string) (err error) {
	for key, value := range data {
		err = s.Set(ctx, key, value)
		if err != nil {
			return
		}
	}
	return
}

