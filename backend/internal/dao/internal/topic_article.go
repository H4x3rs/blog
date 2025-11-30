// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TopicArticleDao is the data access object for the table topic_article.
type TopicArticleDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  TopicArticleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// TopicArticleColumns defines and stores column names for the table topic_article.
type TopicArticleColumns struct {
	Id        string //
	TopicId   string // 专题ID
	ArticleId string // 文章ID
	SortOrder string // 排序
	CreatedAt string //
}

// topicArticleColumns holds the columns for the table topic_article.
var topicArticleColumns = TopicArticleColumns{
	Id:        "id",
	TopicId:   "topic_id",
	ArticleId: "article_id",
	SortOrder: "sort_order",
	CreatedAt: "created_at",
}

// NewTopicArticleDao creates and returns a new DAO object for table data access.
func NewTopicArticleDao(handlers ...gdb.ModelHandler) *TopicArticleDao {
	return &TopicArticleDao{
		group:    "default",
		table:    "topic_article",
		columns:  topicArticleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TopicArticleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TopicArticleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TopicArticleDao) Columns() TopicArticleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TopicArticleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TopicArticleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TopicArticleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
