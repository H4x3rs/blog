package upload

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type UploadFileReq struct {
	g.Meta `path:"/upload/file" method:"post" tags:"Upload" summary:"Upload file"`
}

type UploadImageReq struct {
	g.Meta `path:"/upload/image" method:"post" tags:"Upload" summary:"Upload image"`
}

type UploadRes struct {
	URL      string `json:"url"`      // 文件访问URL
	Filename string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小（字节）
}

// UploadFile 处理文件上传
func (c *ControllerV1) UploadFile(ctx context.Context, req *UploadFileReq) (res *UploadRes, err error) {
	r := g.RequestFromCtx(ctx)
	
	// 检查用户是否登录
	userID := r.GetCtxVar("user_id", 0).Int()
	if userID == 0 {
		return nil, gerror.New("请先登录")
	}

	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, gerror.New("未找到上传文件")
	}

	// 验证文件大小（默认最大 10MB）
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		return nil, gerror.New(fmt.Sprintf("文件大小不能超过 %dMB", maxSize/(1024*1024)))
	}

	// 获取文件扩展名
	ext := gfile.Ext(file.Filename)
	
	// 允许的文件类型
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".zip", ".rar"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, gerror.New("不支持的文件类型")
	}

	// 生成文件名：时间戳_随机数_原文件名
	timestamp := gtime.Now().Format("YmdHis")
	randomStr := fmt.Sprintf("%d", gtime.Now().UnixNano()%10000)
	// 清理文件名，移除特殊字符
	safeFilename := strings.ReplaceAll(file.Filename, " ", "_")
	safeFilename = strings.ReplaceAll(safeFilename, "..", "")
	filename := fmt.Sprintf("%s_%s_%s", timestamp, randomStr, safeFilename)
	
	// 根据文件类型确定存储目录
	var uploadDir string
	if isImage(ext) {
		uploadDir = "uploads/images"
	} else {
		uploadDir = "uploads/files"
	}

	// 创建上传目录（如果不存在）
	uploadPath := filepath.Join("public", uploadDir)
	if !gfile.Exists(uploadPath) {
		if err := gfile.Mkdir(uploadPath); err != nil {
			return nil, gerror.Wrap(err, "创建上传目录失败")
		}
	}

	// 保存文件
	savePath := filepath.Join(uploadPath, filename)
	_, err = file.Save(savePath)
	if err != nil {
		return nil, gerror.Wrap(err, "保存文件失败")
	}

	// 构建访问URL（使用相对路径，前端会自动处理）
	fileURL := fmt.Sprintf("/%s/%s", uploadDir, filename)

	return &UploadRes{
		URL:      fileURL,
		Filename: filename,
		Size:     file.Size,
	}, nil
}

// UploadImage 处理图片上传（专门用于编辑器）
func (c *ControllerV1) UploadImage(ctx context.Context, req *UploadImageReq) (res *UploadRes, err error) {
	r := g.RequestFromCtx(ctx)
	
	// 检查用户是否登录
	userID := r.GetCtxVar("user_id", 0).Int()
	if userID == 0 {
		return nil, gerror.New("请先登录")
	}

	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, gerror.New("未找到上传文件")
	}

	// 验证文件大小（图片最大 5MB）
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return nil, gerror.New("图片大小不能超过 5MB")
	}

	// 验证文件类型（只允许图片）
	ext := gfile.Ext(file.Filename)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, gerror.New("只支持图片格式：jpg, jpeg, png, gif, webp, svg")
	}

	// 生成文件名
	timestamp := gtime.Now().Format("YmdHis")
	randomStr := fmt.Sprintf("%d", gtime.Now().UnixNano()%10000)
	// 清理文件名，移除特殊字符
	safeFilename := strings.ReplaceAll(file.Filename, " ", "_")
	safeFilename = strings.ReplaceAll(safeFilename, "..", "")
	filename := fmt.Sprintf("%s_%s_%s", timestamp, randomStr, safeFilename)
	
	// 图片存储目录
	uploadDir := "uploads/images"
	uploadPath := filepath.Join("public", uploadDir)
	if !gfile.Exists(uploadPath) {
		if err := gfile.Mkdir(uploadPath); err != nil {
			return nil, gerror.Wrap(err, "创建上传目录失败")
		}
	}

	// 保存文件
	savePath := filepath.Join(uploadPath, filename)
	_, err = file.Save(savePath)
	if err != nil {
		return nil, gerror.Wrap(err, "保存文件失败")
	}

	// 构建访问URL（使用相对路径，前端会自动处理）
	fileURL := fmt.Sprintf("/%s/%s", uploadDir, filename)

	return &UploadRes{
		URL:      fileURL,
		Filename: filename,
		Size:     file.Size,
	}, nil
}

// 判断是否为图片文件
func isImage(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

