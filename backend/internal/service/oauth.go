package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"blog/internal/dao"
	"blog/internal/model/entity"
	"blog/internal/utils"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

type sOAuth struct{}

var OAuth = sOAuth{}

// OAuth 用户信息
type OAuthUserInfo struct {
	Token    string
	Username string
	Nickname string
	Email    string
	Avatar   string
	Provider string
	OpenID   string
}

// 获取 OAuth 授权 URL
func (s *sOAuth) GetAuthURL(ctx context.Context, provider string) (string, error) {
	state := generateState()

	// 保存 state 到 session 或缓存（这里简化处理）
	// 实际应该保存到 Redis 或 Session 中，用于验证回调

	var authURL string
	frontendURL := utils.GetConfigString(ctx, "OAUTH_FRONTEND_URL", "oauth.frontendUrl", "http://localhost:5173")
	redirectURI := frontendURL + "/oauth/" + provider + "/callback"

	switch provider {
	case "github":
		clientID := utils.GetConfigString(ctx, "GITHUB_CLIENT_ID", "oauth.github.clientId", "")
		authURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
			clientID, url.QueryEscape(redirectURI), state)
	case "google":
		clientID := utils.GetConfigString(ctx, "GOOGLE_CLIENT_ID", "oauth.google.clientId", "")
		authURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20profile%%20email&state=%s",
			clientID, url.QueryEscape(redirectURI), state)
	case "wechat":
		appID := utils.GetConfigString(ctx, "WECHAT_APP_ID", "oauth.wechat.appId", "")
		authURL = fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
			appID, url.QueryEscape(redirectURI), state)
	default:
		return "", fmt.Errorf("不支持的登录方式: %s", provider)
	}

	return authURL, nil
}

// 处理 OAuth 回调
func (s *sOAuth) HandleCallback(ctx context.Context, provider, code, state string) (*OAuthUserInfo, error) {
	// 验证 state（实际应该从 session 或缓存中验证）
	// 这里简化处理，跳过验证

	var accessToken string
	var userInfo *OAuthUserInfo
	var err error

	switch provider {
	case "github":
		accessToken, err = s.getGitHubAccessToken(ctx, code)
		if err != nil {
			return nil, err
		}
		userInfo, err = s.getGitHubUserInfo(ctx, accessToken)
	case "google":
		accessToken, err = s.getGoogleAccessToken(ctx, code)
		if err != nil {
			return nil, err
		}
		userInfo, err = s.getGoogleUserInfo(ctx, accessToken)
	case "wechat":
		accessToken, openID, err := s.getWeChatAccessToken(ctx, code)
		if err != nil {
			return nil, err
		}
		userInfo, err = s.getWeChatUserInfo(ctx, accessToken, openID)
	default:
		return nil, fmt.Errorf("不支持的登录方式: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	// 查找或创建用户
	user, err := s.findOrCreateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.Id, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}
	userInfo.Token = token
	userInfo.Username = user.Username
	userInfo.Nickname = user.Nickname
	userInfo.Avatar = user.Avatar

	return userInfo, nil
}

// 获取 GitHub Access Token
func (s *sOAuth) getGitHubAccessToken(ctx context.Context, code string) (string, error) {
	clientID := utils.GetConfigString(ctx, "GITHUB_CLIENT_ID", "oauth.github.clientId", "")
	clientSecret := utils.GetConfigString(ctx, "GITHUB_CLIENT_SECRET", "oauth.github.clientSecret", "")
	frontendURL := utils.GetConfigString(ctx, "OAUTH_FRONTEND_URL", "oauth.frontendUrl", "http://localhost:5173")
	redirectURI := frontendURL + "/oauth/github/callback"

	client := gclient.New()
	client.SetHeader("Accept", "application/json")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	resp, err := client.Post(ctx, "https://github.com/login/oauth/access_token", data.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// GitHub 返回的是 form-urlencoded 格式
	if strings.Contains(string(body), "access_token=") {
		params, err := url.ParseQuery(string(body))
		if err != nil {
			return "", err
		}
		return params.Get("access_token"), nil
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("GitHub OAuth错误: %s", result.Error)
	}

	return result.AccessToken, nil
}

// 获取 GitHub 用户信息
func (s *sOAuth) getGitHubUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	client := gclient.New()
	client.SetHeader("Authorization", "Bearer "+accessToken)
	client.SetHeader("Accept", "application/json")

	resp, err := client.Get(ctx, "https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// 如果没有邮箱，尝试获取邮箱
	if githubUser.Email == "" {
		resp, err := client.Get(ctx, "https://api.github.com/user/emails")
		if err == nil {
			defer resp.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if json.NewDecoder(resp.Body).Decode(&emails) == nil {
				for _, email := range emails {
					if email.Primary {
						githubUser.Email = email.Email
						break
					}
				}
			}
		}
	}

	return &OAuthUserInfo{
		Username: fmt.Sprintf("github_%d", githubUser.ID),
		Nickname: githubUser.Name,
		Email:    githubUser.Email,
		Avatar:   githubUser.AvatarURL,
		Provider: "github",
		OpenID:   fmt.Sprintf("%d", githubUser.ID),
	}, nil
}

// 获取 Google Access Token
func (s *sOAuth) getGoogleAccessToken(ctx context.Context, code string) (string, error) {
	clientID := utils.GetConfigString(ctx, "GOOGLE_CLIENT_ID", "oauth.google.clientId", "")
	clientSecret := utils.GetConfigString(ctx, "GOOGLE_CLIENT_SECRET", "oauth.google.clientSecret", "")
	frontendURL := utils.GetConfigString(ctx, "OAUTH_FRONTEND_URL", "oauth.frontendUrl", "http://localhost:5173")
	redirectURI := frontendURL + "/oauth/google/callback"

	client := gclient.New()
	client.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	resp, err := client.Post(ctx, "https://oauth2.googleapis.com/token", data.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("Google OAuth错误: %s", result.Error)
	}

	return result.AccessToken, nil
}

// 获取 Google 用户信息
func (s *sOAuth) getGoogleUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	client := gclient.New()
	resp, err := client.Get(ctx, fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &OAuthUserInfo{
		Username: fmt.Sprintf("google_%s", googleUser.ID),
		Nickname: googleUser.Name,
		Email:    googleUser.Email,
		Avatar:   googleUser.Picture,
		Provider: "google",
		OpenID:   googleUser.ID,
	}, nil
}

// 获取微信 Access Token
func (s *sOAuth) getWeChatAccessToken(ctx context.Context, code string) (string, string, error) {
	appID := utils.GetConfigString(ctx, "WECHAT_APP_ID", "oauth.wechat.appId", "")
	appSecret := utils.GetConfigString(ctx, "WECHAT_APP_SECRET", "oauth.wechat.appSecret", "")

	client := gclient.New()
	resp, err := client.Get(ctx, fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, appSecret, code))
	if err != nil {
		return "", "", err
	}
	defer resp.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
		UnionID      string `json:"unionid"`
		Errcode      int    `json:"errcode"`
		Errmsg       string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if result.Errcode != 0 {
		return "", "", fmt.Errorf("微信OAuth错误: %s", result.Errmsg)
	}

	return result.AccessToken, result.OpenID, nil
}

// 获取微信用户信息
func (s *sOAuth) getWeChatUserInfo(ctx context.Context, accessToken, openID string) (*OAuthUserInfo, error) {
	client := gclient.New()
	resp, err := client.Get(ctx, fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openID))
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var wechatUser struct {
		OpenID     string `json:"openid"`
		Nickname   string `json:"nickname"`
		HeadImgURL string `json:"headimgurl"`
		UnionID    string `json:"unionid"`
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wechatUser); err != nil {
		return nil, err
	}

	if wechatUser.Errcode != 0 {
		return nil, fmt.Errorf("微信API错误: %s", wechatUser.Errmsg)
	}

	return &OAuthUserInfo{
		Username: fmt.Sprintf("wechat_%s", wechatUser.OpenID),
		Nickname: wechatUser.Nickname,
		Email:    "",
		Avatar:   wechatUser.HeadImgURL,
		Provider: "wechat",
		OpenID:   wechatUser.OpenID,
	}, nil
}

// 查找或创建用户
func (s *sOAuth) findOrCreateUser(ctx context.Context, oauthInfo *OAuthUserInfo) (*entity.User, error) {
	// 根据 provider 和 openID 查找用户
	// 这里简化处理，实际应该在 user 表中添加 provider 和 open_id 字段
	username := oauthInfo.Username

	var user *entity.User
	err := dao.User.Ctx(ctx).Where(dao.User.Columns().Username, username).Scan(&user)
	if err == nil && user != nil {
		// 用户已存在，更新头像等信息
		if oauthInfo.Avatar != "" && user.Avatar != oauthInfo.Avatar {
			dao.User.Ctx(ctx).Data(g.Map{
				"avatar": oauthInfo.Avatar,
			}).Where(dao.User.Columns().Id, user.Id).Update()
			user.Avatar = oauthInfo.Avatar
		}
		return user, nil
	}

	// 创建新用户
	nickname := oauthInfo.Nickname
	if nickname == "" {
		nickname = username
	}

	result, err := dao.User.Ctx(ctx).Data(g.Map{
		"username": username,
		"password": "", // OAuth 用户不需要密码
		"nickname": nickname,
		"email":    oauthInfo.Email,
		"avatar":   oauthInfo.Avatar,
		"status":   1,
	}).Insert()
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	user = &entity.User{
		Id:       int(id),
		Username: username,
		Nickname: nickname,
		Email:    oauthInfo.Email,
		Avatar:   oauthInfo.Avatar,
		Status:   1,
	}

	// 为新用户分配默认user角色
	err = s.assignDefaultUserRole(ctx, int(id))
	if err != nil {
		// 记录错误但不影响用户创建
		g.Log().Error(ctx, "分配默认角色失败:", err)
	}

	return user, nil
}

// 分配默认user角色
func (s *sOAuth) assignDefaultUserRole(ctx context.Context, userId int) error {
	// 查找user角色
	userRole, err := Role.GetByKey(ctx, "user")
	if err != nil || userRole == nil {
		// 如果user角色不存在，创建它
		userRole = &entity.Role{
			Name:   "用户",
			Key:    "user",
			Desc:   "普通用户角色",
			Status: 1,
		}
		roleId, err := Role.Create(ctx, userRole)
		if err != nil {
			return err
		}
		userRole.Id = roleId
	}

	// 分配角色给用户
	return User.SetRoles(ctx, userId, []int{userRole.Id})
}

// 生成随机 state
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
