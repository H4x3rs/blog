package operation_log

import (
	"context"

	"blog/internal/model/entity"
	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type GetListReq struct {
	g.Meta        `path:"/operationLog/getList" method:"post" tags:"OperationLog" summary:"Get operation log list"`
	UserId        int    `json:"userId" v:""`
	OperationType string `json:"operationType" v:""`
	Module        string `json:"module" v:""`
	Page          int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size          int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.OperationLog `json:"list"`
	Total int                    `json:"total"`
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.OperationLog.GetList(ctx, req.UserId, req.OperationType, req.Module, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}


