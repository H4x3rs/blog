package chat

import (
	"context"
	"fmt"
	"strings"

	"blog/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// 对话请求
type ChatReq struct {
	g.Meta    `path:"/chat/message" method:"post" tags:"Chat" summary:"Send chat message"`
	Message   string `json:"message" v:"required#消息内容不能为空"`
	SessionId string `json:"sessionId" v:""` // 会话ID，用于多轮对话
}

// 对话响应
type ChatRes struct {
	Message   string `json:"message"`   // AI回复内容
	SessionId string `json:"sessionId"` // 会话ID
}

// 发送消息（流式）
func (c *ControllerV1) SendMessage(ctx context.Context, req *ChatReq) (res *ChatRes, err error) {
	// 从context中获取用户ID
	userId := g.RequestFromCtx(ctx).GetCtxVar("user_id").Int()
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 检查是否请求流式响应
	r := g.RequestFromCtx(ctx)
	acceptHeader := r.Header.Get("Accept")
	isStream := strings.Contains(acceptHeader, "text/event-stream")

	if isStream {
		// 流式响应
		return c.sendMessageStream(ctx, r, userId, req.Message, req.SessionId)
	}

	// 非流式响应（保持兼容）
	reply, sessionId, err := service.Chat.SendMessage(ctx, userId, req.Message, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &ChatRes{
		Message:   reply,
		SessionId: sessionId,
	}, nil
}

// 流式发送消息
func (c *ControllerV1) sendMessageStream(ctx context.Context, r *ghttp.Request, userId int, message string, sessionId string) (*ChatRes, error) {
	// 设置SSE响应头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("X-Accel-Buffering", "no") // 禁用nginx缓冲

	// 调用流式service
	sessionId, err := service.Chat.SendMessageStream(ctx, r.Response, userId, message, sessionId)
	if err != nil {
		// 发送错误事件
		r.Response.Writefln("event: error\ndata: %s\n\n", err.Error())
		return nil, err
	}

	// 发送完成事件
	doneData := fmt.Sprintf(`{"sessionId":"%s"}`, sessionId)
	r.Response.Writefln("event: done\ndata: %s\n\n", doneData)
	r.Response.Flush()

	return &ChatRes{
		SessionId: sessionId,
	}, nil
}

// 清除会话历史
type ClearSessionReq struct {
	g.Meta    `path:"/chat/clearSession" method:"post" tags:"Chat" summary:"Clear chat session"`
	SessionId string `json:"sessionId" v:"required#会话ID不能为空"`
}

type ClearSessionRes struct{}

// 清除会话
func (c *ControllerV1) ClearSession(ctx context.Context, req *ClearSessionReq) (res *ClearSessionRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("user_id").Int()
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	err = service.Chat.ClearSession(ctx, userId, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &ClearSessionRes{}, nil
}
