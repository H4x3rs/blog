package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"

	"blog/internal/controller/article"
	"blog/internal/controller/category"
	"blog/internal/controller/chat"
	"blog/internal/controller/hello"
	"blog/internal/controller/menu"
	"blog/internal/controller/oauth"
	"blog/internal/controller/permission"
	"blog/internal/controller/role"
	"blog/internal/controller/settings"
	"blog/internal/controller/tag"
	"blog/internal/controller/topic"
	"blog/internal/controller/upload"
	"blog/internal/controller/user"
	"blog/internal/middleware"
	"blog/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 初始化OSS服务（如果配置了OSS）
			if err := service.OSS.InitOSS(ctx); err != nil {
				g.Log().Warning(ctx, "OSS初始化失败（如果使用OSS存储，请检查配置）:", err)
			} else {
				g.Log().Info(ctx, "OSS服务初始化成功")
			}

			// 配置静态文件服务（用于访问上传的文件）
			// 注意：如果使用OSS存储，这部分代码可以保留作为备用，或者删除
			// 使用绝对路径，相对于 backend 目录
			publicPath := "public"
			if !gfile.Exists(publicPath) {
				if err := gfile.Mkdir(publicPath); err != nil {
					g.Log().Warning(ctx, "创建 public 目录失败:", err)
				}
			}
			uploadsPath := gfile.Join(publicPath, "uploads")
			if !gfile.Exists(uploadsPath) {
				if err := gfile.Mkdir(uploadsPath); err != nil {
					g.Log().Warning(ctx, "创建 uploads 目录失败:", err)
				}
			}
			imagesPath := gfile.Join(uploadsPath, "images")
			if !gfile.Exists(imagesPath) {
				if err := gfile.Mkdir(imagesPath); err != nil {
					g.Log().Warning(ctx, "创建 images 目录失败:", err)
				}
			}
			filesPath := gfile.Join(uploadsPath, "files")
			if !gfile.Exists(filesPath) {
				if err := gfile.Mkdir(filesPath); err != nil {
					g.Log().Warning(ctx, "创建 files 目录失败:", err)
				}
			}

			s.SetServerRoot(publicPath)
			s.AddStaticPath("/uploads", uploadsPath)

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(func(r *ghttp.Request) {
					r.Response.CORSDefault()
					r.Middleware.Next()
				})
				// 添加认证中间件（可选，某些接口需要认证）
				group.Middleware(middleware.Auth)
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
					article.NewV1(),
					category.NewV1(),
					user.NewV1(),
					tag.NewV1(),
					topic.NewV1(),
					role.NewV1(),
					permission.NewV1(),
					menu.NewV1(),
					settings.NewV1(),
					oauth.NewV1(),
					upload.NewV1(),
					chat.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
