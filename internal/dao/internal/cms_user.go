// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmsUserDao is the data access object for the table cms_user.
type CmsUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CmsUserColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CmsUserColumns defines and stores column names for the table cms_user.
type CmsUserColumns struct {
	Id         string // 自增id
	Name       string // 用户名
	Email      string // 邮箱
	Mobile     string // 电话
	Pwd        string // 密码
	FeishuCode string // 飞书的唯一标识符
	CreateAt   string // 创建时间
	UpdateAt   string // 修改时间
}

// cmsUserColumns holds the columns for the table cms_user.
var cmsUserColumns = CmsUserColumns{
	Id:         "id",
	Name:       "name",
	Email:      "email",
	Mobile:     "mobile",
	Pwd:        "pwd",
	FeishuCode: "feishu_code",
	CreateAt:   "create_at",
	UpdateAt:   "update_at",
}

// NewCmsUserDao creates and returns a new DAO object for table data access.
func NewCmsUserDao(handlers ...gdb.ModelHandler) *CmsUserDao {
	return &CmsUserDao{
		group:    "default",
		table:    "cms_user",
		columns:  cmsUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CmsUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CmsUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CmsUserDao) Columns() CmsUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CmsUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CmsUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CmsUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
