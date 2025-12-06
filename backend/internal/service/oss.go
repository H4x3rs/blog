package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	"blog/internal/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

type sOSS struct {
	client   *oss.Client
	bucket   *oss.Bucket
	endpoint string
	bucketName string
	domain   string // 自定义域名，如果配置了CDN域名
}

var OSS = sOSS{}

// InitOSS 初始化OSS客户端
func (s *sOSS) InitOSS(ctx context.Context) error {
	// 从配置读取OSS配置
	endpoint := utils.GetConfigString(ctx, "OSS_ENDPOINT", "oss.endpoint", "")
	accessKeyID := utils.GetConfigString(ctx, "OSS_ACCESS_KEY_ID", "oss.accessKeyId", "")
	accessKeySecret := utils.GetConfigString(ctx, "OSS_ACCESS_KEY_SECRET", "oss.accessKeySecret", "")
	bucketName := utils.GetConfigString(ctx, "OSS_BUCKET_NAME", "oss.bucketName", "")
	domain := utils.GetConfigString(ctx, "OSS_DOMAIN", "oss.domain", "")

	// 检查必要配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" {
		return fmt.Errorf("OSS配置不完整，请检查配置文件或环境变量")
	}

	// 创建OSS客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("获取OSS Bucket失败: %v", err)
	}

	s.client = client
	s.bucket = bucket
	s.endpoint = endpoint
	s.bucketName = bucketName
	s.domain = domain

	g.Log().Info(ctx, "OSS客户端初始化成功")
	return nil
}

// UploadFile 上传文件到OSS
// objectKey: OSS中的对象键（文件路径），如 "uploads/images/20240101_123456_photo.jpg"
// reader: 文件内容读取器
// contentType: 文件MIME类型，如 "image/jpeg"
func (s *sOSS) UploadFile(ctx context.Context, objectKey string, reader io.Reader, contentType string) error {
	if s.bucket == nil {
		if err := s.InitOSS(ctx); err != nil {
			return err
		}
	}

	// 设置上传选项
	options := []oss.Option{
		oss.ContentType(contentType),
	}

	// 上传文件
	err := s.bucket.PutObject(objectKey, reader, options...)
	if err != nil {
		return fmt.Errorf("上传文件到OSS失败: %v", err)
	}

	return nil
}

// GetFileURL 获取文件的访问URL
// objectKey: OSS中的对象键
// 如果配置了自定义域名，返回自定义域名URL；否则返回OSS默认URL
func (s *sOSS) GetFileURL(ctx context.Context, objectKey string) string {
	if s.bucket == nil {
		if err := s.InitOSS(ctx); err != nil {
			g.Log().Error(ctx, "OSS未初始化:", err)
			return ""
		}
	}

	// 如果配置了自定义域名（CDN域名），使用自定义域名
	if s.domain != "" {
		// 确保域名不以/结尾
		domain := strings.TrimSuffix(s.domain, "/")
		// 确保objectKey不以/开头
		objectKey = strings.TrimPrefix(objectKey, "/")
		return fmt.Sprintf("%s/%s", domain, objectKey)
	}

	// 使用OSS默认域名
	// 格式: https://bucket-name.endpoint/object-key
	// 例如: https://my-bucket.oss-cn-hangzhou.aliyuncs.com/uploads/images/photo.jpg
	endpoint := s.endpoint
	// 如果endpoint包含协议，移除协议部分（OSS SDK只需要域名部分）
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	// 确保endpoint不以/结尾
	endpoint = strings.TrimSuffix(endpoint, "/")
	// 确保objectKey不以/开头
	objectKey = strings.TrimPrefix(objectKey, "/")
	
	return fmt.Sprintf("https://%s.%s/%s", s.bucketName, endpoint, objectKey)
}

// DeleteFile 删除OSS中的文件
func (s *sOSS) DeleteFile(ctx context.Context, objectKey string) error {
	if s.bucket == nil {
		if err := s.InitOSS(ctx); err != nil {
			return err
		}
	}

	err := s.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("删除OSS文件失败: %v", err)
	}

	return nil
}

// FileExists 检查文件是否存在
func (s *sOSS) FileExists(ctx context.Context, objectKey string) (bool, error) {
	if s.bucket == nil {
		if err := s.InitOSS(ctx); err != nil {
			return false, err
		}
	}

	exists, err := s.bucket.IsObjectExist(objectKey)
	if err != nil {
		return false, fmt.Errorf("检查OSS文件是否存在失败: %v", err)
	}

	return exists, nil
}

// GetContentType 根据文件扩展名获取MIME类型
func GetContentType(filename string) string {
	ext := strings.ToLower(gfile.Ext(filename))
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".bmp":  "image/bmp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
	}

	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}

	return "application/octet-stream"
}

