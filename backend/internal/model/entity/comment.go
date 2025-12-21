// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Comment is the golang structure for table comment.
type Comment struct {
	Id        int         `json:"id"         orm:"id"          ` //
	ArticleId int         `json:"articleId"  orm:"article_id"  ` // 文章ID
	UserId    *int        `json:"userId"     orm:"user_id"     ` // 用户ID（可为空，支持匿名评论）
	ParentId  int         `json:"parentId"   orm:"parent_id"   ` // 父评论ID（0表示顶级评论）
	Content   string      `json:"content"    orm:"content"     ` // 评论内容
	Status    string      `json:"status"      orm:"status"       ` // 状态:approved已审核,pending待审核,rejected已拒绝
	CreatedAt *gtime.Time `json:"createdAt"  orm:"created_at"   ` //
	UpdatedAt *gtime.Time `json:"updatedAt"  orm:"updated_at"   ` //
}

// CommentWithUser 包含用户信息的评论
type CommentWithUser struct {
	*Comment
	User         *User   `json:"user,omitempty"`         // 用户信息
	Replies      []*CommentWithUser `json:"replies,omitempty"` // 回复列表
	ReplyCount   int     `json:"replyCount"`             // 回复数量
}


