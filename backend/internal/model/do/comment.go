// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Comment is the golang structure of table comment for DAO operations like Where/Data.
type Comment struct {
	g.Meta    `orm:"table:comment, do:true"`
	Id        any         //
	ArticleId any         // 文章ID
	UserId    any         // 用户ID（可为空，支持匿名评论）
	ParentId  any         // 父评论ID（0表示顶级评论）
	Content   any         // 评论内容
	Status    any         // 状态:approved已审核,pending待审核,rejected已拒绝
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}


