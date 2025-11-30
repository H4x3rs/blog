package service

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/do"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sCategory struct{}

var Category = sCategory{}

func (s *sCategory) Create(ctx context.Context, in *entity.Category) (id int, err error) {
	result, err := dao.Category.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sCategory) Update(ctx context.Context, in *entity.Category) (err error) {
	// 使用 do.Category 构建更新数据
	// 注意：do.Category 的字段类型是 any，可以接受字符串值，包括空字符串
	updateData := do.Category{
		Name:        in.Name,
		Slug:        in.Slug,
		Icon:        in.Icon,        // 明确设置，即使是空字符串
		Description: in.Description, // 明确设置，即使是空字符串
	}

	g.Log().Debugf(ctx, "更新分类 ID=%d, 数据=%+v", in.Id, updateData)

	// 使用 Fields() 明确指定要更新的字段
	// 不使用 OmitEmpty()，确保所有字段都被更新，包括空字符串
	result, err := dao.Category.Ctx(ctx).
		Data(updateData).
		Fields(dao.Category.Columns().Name, dao.Category.Columns().Slug, dao.Category.Columns().Icon, dao.Category.Columns().Description).
		Where(dao.Category.Columns().Id, in.Id).
		Update()

	if err != nil {
		g.Log().Errorf(ctx, "更新分类失败: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	g.Log().Debugf(ctx, "更新分类成功，影响行数: %d", rowsAffected)

	return
}

func (s *sCategory) GetList(ctx context.Context, name, slug string, page, size int) (list []*entity.Category, total int, err error) {
	m := dao.Category.Ctx(ctx)

	// 筛选条件
	if name != "" {
		m = m.WhereLike(dao.Category.Columns().Name, "%"+name+"%")
	}
	if slug != "" {
		m = m.WhereLike(dao.Category.Columns().Slug, "%"+slug+"%")
	}

	// 获取总数
	total, err = m.Count()
	if err != nil {
		return
	}

	// 分页查询
	err = m.OrderDesc(dao.Category.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sCategory) GetOne(ctx context.Context, id int) (category *entity.Category, err error) {
	err = dao.Category.Ctx(ctx).Where(dao.Category.Columns().Id, id).Scan(&category)
	return
}

func (s *sCategory) GetBySlug(ctx context.Context, slug string) (category *entity.Category, err error) {
	err = dao.Category.Ctx(ctx).Where(dao.Category.Columns().Slug, slug).Scan(&category)
	return
}

func (s *sCategory) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.Category.Ctx(ctx).Where(dao.Category.Columns().Id, id).Delete()
	return
}
