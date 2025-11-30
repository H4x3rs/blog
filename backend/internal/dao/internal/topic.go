// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TopicDao is the data access object for the table topic.
type TopicDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TopicColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TopicColumns defines and stores column names for the table topic.
type TopicColumns struct {
	Id          string //
	Name        string // 专题名称
	Slug        string // 别名
	Description string // 专题描述
	CoverImage  string // 封面图片
	Intro       string // 专题介绍
	IsFeatured  string // 是否置顶:0否,1是
	SortOrder   string // 排序
	CreatedAt   string //
	UpdatedAt   string //
}

// topicColumns holds the columns for the table topic.
var topicColumns = TopicColumns{
	Id:          "id",
	Name:        "name",
	Slug:        "slug",
	Description: "description",
	CoverImage:  "cover_image",
	Intro:       "intro",
	IsFeatured:  "is_featured",
	SortOrder:   "sort_order",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewTopicDao creates and returns a new DAO object for table data access.
func NewTopicDao(handlers ...gdb.ModelHandler) *TopicDao {
	return &TopicDao{
		group:    "default",
		table:    "topic",
		columns:  topicColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TopicDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TopicDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TopicDao) Columns() TopicColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TopicDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TopicDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TopicDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
