// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ArticleDao is the data access object for the table article.
type ArticleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ArticleColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ArticleColumns defines and stores column names for the table article.
type ArticleColumns struct {
	Id         string //
	Title      string // 文章标题
	Desc       string // 文章摘要
	Content    string // 文章内容
	CoverImage string // 封面图片
	CategoryId string // 分类ID
	Status      string // 状态:published,draft
	Views       string // 阅读数
	CreatedBy   string // 创建人员ID
	UpdatedBy   string // 更新人员ID
	PublishedBy string // 发布人员ID
	CreatedAt   string //
	UpdatedAt   string //
}

// articleColumns holds the columns for the table article.
var articleColumns = ArticleColumns{
	Id:         "id",
	Title:      "title",
	Desc:       "desc",
	Content:    "content",
	CoverImage: "cover_image",
	CategoryId: "category_id",
	Status:      "status",
	Views:       "views",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	PublishedBy: "published_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewArticleDao creates and returns a new DAO object for table data access.
func NewArticleDao(handlers ...gdb.ModelHandler) *ArticleDao {
	return &ArticleDao{
		group:    "default",
		table:    "article",
		columns:  articleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ArticleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ArticleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ArticleDao) Columns() ArticleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ArticleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ArticleDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ArticleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
