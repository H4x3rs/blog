// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VisitLogDao is the data access object for the table visit_log.
type VisitLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  VisitLogColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// VisitLogColumns defines and stores column names for the table visit_log.
type VisitLogColumns struct {
	Id        string //
	Date      string // 访问日期
	Views     string // 访问量
	CreatedAt string //
	UpdatedAt string //
}

// visitLogColumns holds the columns for the table visit_log.
var visitLogColumns = VisitLogColumns{
	Id:        "id",
	Date:      "date",
	Views:     "views",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewVisitLogDao creates and returns a new DAO object for table data access.
func NewVisitLogDao(handlers ...gdb.ModelHandler) *VisitLogDao {
	return &VisitLogDao{
		group:    "default",
		table:    "visit_log",
		columns:  visitLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *VisitLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *VisitLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *VisitLogDao) Columns() VisitLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *VisitLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *VisitLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *VisitLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}


