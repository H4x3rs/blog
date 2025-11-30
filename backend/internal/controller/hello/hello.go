package hello

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type HelloReq struct {
	g.Meta `path:"/hello" method:"get"`
}

type HelloRes struct {
	Message string `json:"message"`
}

func (c *ControllerV1) Hello(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	return &HelloRes{
		Message: "Hello from GoFrame Backend!",
	}, nil
}

