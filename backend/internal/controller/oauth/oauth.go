package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// OAuth 登录请求
type OAuthLoginReq struct {
	g.Meta `path:"/oauth/:provider/login" method:"get" tags:"OAuth" summary:"OAuth login"`
	Provider string `json:"provider" in:"path" v:"required|in:github,google,wechat"`
}

// OAuth 登录 - 重定向到第三方登录页面
func (c *ControllerV1) Login(ctx context.Context, req *OAuthLoginReq) (res *OAuthLoginRes, err error) {
	redirectURL, err := service.OAuth.GetAuthURL(ctx, req.Provider)
	if err != nil {
		return nil, err
	}
	return &OAuthLoginRes{RedirectURL: redirectURL}, nil
}

// OAuth 回调请求
type OAuthCallbackReq struct {
	g.Meta `path:"/oauth/:provider/callback" method:"get" tags:"OAuth" summary:"OAuth callback"`
	Provider string `json:"provider" in:"path" v:"required|in:github,google,wechat"`
	Code     string `json:"code" in:"query"`
	State    string `json:"state" in:"query"`
}

// OAuth 回调 - 处理第三方登录回调，返回JSON（包含JWT token）
func (c *ControllerV1) Callback(ctx context.Context, req *OAuthCallbackReq) (res *OAuthCallbackRes, err error) {
	userInfo, err := service.OAuth.HandleCallback(ctx, req.Provider, req.Code, req.State)
	if err != nil {
		return nil, err
	}
	return &OAuthCallbackRes{
		Token:    userInfo.Token,
		Username: userInfo.Username,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
	}, nil
}

// OAuth 登录响应
type OAuthLoginRes struct {
	RedirectURL string `json:"redirectUrl"`
}

// OAuth 回调响应
type OAuthCallbackRes struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// GitHub 用户信息结构
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// Google 用户信息结构
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// 微信用户信息结构
type WeChatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}

// 获取 GitHub 用户信息
func GetGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	client := gclient.New()
	client.SetHeader("Authorization", "Bearer "+accessToken)
	
	resp, err := client.Get(context.Background(), "https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

// 获取 Google 用户信息
func GetGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	client := gclient.New()
	resp, err := client.Get(context.Background(), fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

// 获取微信用户信息
func GetWeChatUserInfo(accessToken, openID string) (*WeChatUserInfo, error) {
	client := gclient.New()
	resp, err := client.Get(context.Background(), fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openID))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OpenID     string `json:"openid"`
		Nickname   string `json:"nickname"`
		HeadImgURL string `json:"headimgurl"`
		UnionID    string `json:"unionid"`
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Errcode != 0 {
		return nil, fmt.Errorf("微信API错误: %s", result.Errmsg)
	}

	return &WeChatUserInfo{
		OpenID:     result.OpenID,
		Nickname:   result.Nickname,
		HeadImgURL: result.HeadImgURL,
		UnionID:    result.UnionID,
	}, nil
}

