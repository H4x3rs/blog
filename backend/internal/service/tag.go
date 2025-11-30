package service

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/entity"
)

type sTag struct{}

var Tag = sTag{}

func (s *sTag) Create(ctx context.Context, in *entity.Tag) (id int, err error) {
	result, err := dao.Tag.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sTag) Update(ctx context.Context, in *entity.Tag) (err error) {
	_, err = dao.Tag.Ctx(ctx).Data(in).Where(dao.Tag.Columns().Id, in.Id).Update()
	return
}

func (s *sTag) GetList(ctx context.Context, name, slug string, page, size int) (list []*entity.Tag, total int, err error) {
	m := dao.Tag.Ctx(ctx)

	if name != "" {
		m = m.WhereLike(dao.Tag.Columns().Name, "%"+name+"%")
	}
	if slug != "" {
		m = m.WhereLike(dao.Tag.Columns().Slug, "%"+slug+"%")
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.Tag.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sTag) GetOne(ctx context.Context, id int) (tag *entity.Tag, err error) {
	err = dao.Tag.Ctx(ctx).Where(dao.Tag.Columns().Id, id).Scan(&tag)
	return
}

func (s *sTag) GetBySlug(ctx context.Context, slug string) (tag *entity.Tag, err error) {
	err = dao.Tag.Ctx(ctx).Where(dao.Tag.Columns().Slug, slug).Scan(&tag)
	return
}

func (s *sTag) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.Tag.Ctx(ctx).Where(dao.Tag.Columns().Id, id).Delete()
	return
}

