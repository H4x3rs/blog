package settings

import (
	"context"

	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type GetSettingsReq struct {
	g.Meta `path:"/settings/getSettings" method:"post" tags:"Settings" summary:"Get settings"`
}
type GetSettingsRes struct {
	SiteName    string `json:"siteName"`
	SiteDesc    string `json:"siteDesc"`
	Keywords    string `json:"keywords"`
	Copyright   string `json:"copyright"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	IcpNumber   string `json:"icpNumber"`
	AboutAvatar string `json:"aboutAvatar"`
	AboutSubtitle string `json:"aboutSubtitle"`
	AboutContent string `json:"aboutContent"`
	AboutEmail   string `json:"aboutEmail"`
	AboutGithub  string `json:"aboutGithub"`
}

type UpdateSettingsReq struct {
	g.Meta        `path:"/settings/updateSettings" method:"post" tags:"Settings" summary:"Update settings"`
	SiteName      string `json:"siteName"`
	SiteDesc      string `json:"siteDesc"`
	Keywords      string `json:"keywords"`
	Copyright     string `json:"copyright"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	IcpNumber     string `json:"icpNumber"`
	AboutAvatar   string `json:"aboutAvatar"`
	AboutSubtitle string `json:"aboutSubtitle"`
	AboutContent  string `json:"aboutContent"`
	AboutEmail    string `json:"aboutEmail"`
	AboutGithub   string `json:"aboutGithub"`
}
type UpdateSettingsRes struct{}

// 前台获取 Banner 标题和副标题
type GetBannerReq struct {
	g.Meta `path:"/settings/getBanner" method:"post" tags:"Settings" summary:"Get banner title and subtitle for frontend"`
}
type GetBannerRes struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	SiteName string `json:"siteName"` // 网站名称，用于浏览器标题
	IcpNumber string `json:"icpNumber"` // ICP备案号
}

// 前台获取关于我页面信息
type GetAboutReq struct {
	g.Meta `path:"/settings/getAbout" method:"post" tags:"Settings" summary:"Get about page info for frontend"`
}
type GetAboutRes struct {
	Avatar   string `json:"avatar"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
	Email    string `json:"email"`
	Github   string `json:"github"`
}

func (c *ControllerV1) GetBanner(ctx context.Context, req *GetBannerReq) (res *GetBannerRes, err error) {
	title, err := service.Settings.Get(ctx, "title")
	if err != nil || title == "" {
		title = "分享编程心得\n记录技术成长" // 默认值
	}
	subtitle, err := service.Settings.Get(ctx, "subtitle")
	if err != nil || subtitle == "" {
		subtitle = "分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。" // 默认值
	}
	siteName, err := service.Settings.Get(ctx, "siteName")
	if err != nil || siteName == "" {
		siteName = "Blog System" // 默认值
	}
	icpNumber, _ := service.Settings.Get(ctx, "icpNumber")
	return &GetBannerRes{
		Title:     title,
		Subtitle:  subtitle,
		SiteName:  siteName,
		IcpNumber: icpNumber,
	}, nil
}

func (c *ControllerV1) GetAbout(ctx context.Context, req *GetAboutReq) (res *GetAboutRes, err error) {
	avatar, _ := service.Settings.Get(ctx, "aboutAvatar")
	subtitle, _ := service.Settings.Get(ctx, "aboutSubtitle")
	content, _ := service.Settings.Get(ctx, "aboutContent")
	email, _ := service.Settings.Get(ctx, "aboutEmail")
	github, _ := service.Settings.Get(ctx, "aboutGithub")

	// 设置默认值
	if avatar == "" {
		avatar = "https://picsum.photos/id/64/200/200"
	}

	if subtitle == "" {
		subtitle = "Full Stack Developer / Open Source Enthusiast"
	}
	if content == "" {
		content = "我是 Admin，一名热爱技术的全栈开发者。\n\n## 技术栈\n\n- **后端**: Go, GoFrame, MySQL, Redis\n- **前端**: Vue 3, Element Plus, Vite\n- **工具**: Docker, Git, Linux\n\n## 关于这个博客\n\n这个博客系统是我用来记录技术成长、分享编程心得的小天地。在这里我会分享：\n\n- Go 语言的高性能并发模型\n- Vue 3 的前端开发实践\n- 云原生技术的应用经验\n- 系统架构设计思路\n\n欢迎交流学习！"
	}
	if email == "" {
		email = "admin@example.com"
	}

	return &GetAboutRes{
		Avatar:   avatar,
		Subtitle: subtitle,
		Content:  content,
		Email:    email,
		Github:   github,
	}, nil
}

func (c *ControllerV1) GetSettings(ctx context.Context, req *GetSettingsReq) (res *GetSettingsRes, err error) {
	settings, err := service.Settings.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &GetSettingsRes{
		SiteName:      settings["siteName"],
		SiteDesc:      settings["siteDesc"],
		Keywords:      settings["keywords"],
		Copyright:     settings["copyright"],
		Title:         settings["title"],
		Subtitle:      settings["subtitle"],
		IcpNumber:     settings["icpNumber"],
		AboutAvatar:   settings["aboutAvatar"],
		AboutSubtitle: settings["aboutSubtitle"],
		AboutContent:  settings["aboutContent"],
		AboutEmail:    settings["aboutEmail"],
		AboutGithub:   settings["aboutGithub"],
	}, nil
}

func (c *ControllerV1) UpdateSettings(ctx context.Context, req *UpdateSettingsReq) (res *UpdateSettingsRes, err error) {
	data := map[string]string{
		"siteName":      req.SiteName,
		"siteDesc":      req.SiteDesc,
		"keywords":      req.Keywords,
		"copyright":     req.Copyright,
		"title":         req.Title,
		"subtitle":      req.Subtitle,
		"icpNumber":     req.IcpNumber,
		"aboutAvatar":   req.AboutAvatar,
		"aboutSubtitle": req.AboutSubtitle,
		"aboutContent":  req.AboutContent,
		"aboutEmail":    req.AboutEmail,
		"aboutGithub":   req.AboutGithub,
	}
	err = service.Settings.BatchSet(ctx, data)
	return
}

