package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WechatOAuthConfig holds WeChat OAuth configuration.
type WechatOAuthConfig struct {
	Enabled             bool   `yaml:"enabled"`
	AppID               string `yaml:"app_id"`
	AppSecret           string `yaml:"app_secret"`
	RedirectURI         string `yaml:"redirect_uri"`
	FrontendRedirectURI string `yaml:"frontend_redirect_uri"`
}

// WechatUser represents a WeChat user profile.
type WechatUser struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
}

// WechatOAuth implements WeChat OAuth 2.0 login.
type WechatOAuth struct {
	cfg    WechatOAuthConfig
	client *http.Client
}

func NewWechatOAuth(cfg WechatOAuthConfig) *WechatOAuth {
	return &WechatOAuth{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetAuthURL generates the WeChat authorization URL.
func (w *WechatOAuth) GetAuthURL(state string) string {
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("scope", "snsapi_login")
	values.Set("state", state)
	if !w.cfg.Enabled {
		values.Set("appid", "MOCK_APPID")
		values.Set("redirect_uri", "MOCK_REDIRECT")
		return "https://open.weixin.qq.com/connect/qrconnect?" + values.Encode()
	}
	values.Set("appid", w.cfg.AppID)
	values.Set("redirect_uri", w.cfg.RedirectURI)
	return "https://open.weixin.qq.com/connect/qrconnect?" + values.Encode()
}

// GetFrontendCallbackURL builds the frontend callback URL for deployments where
// WeChat redirects to the backend first.
func (w *WechatOAuth) GetFrontendCallbackURL(code, state string) string {
	if w.cfg.FrontendRedirectURI == "" {
		return ""
	}
	u, err := url.Parse(w.cfg.FrontendRedirectURI)
	if err != nil {
		return ""
	}
	q := u.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// GetUserInfo exchanges code for access_token and fetches user info.
func (w *WechatOAuth) GetUserInfo(code string) (*WechatUser, error) {
	if !w.cfg.Enabled {
		return &WechatUser{OpenID: "mock_openid", Nickname: "wx_mock_user", Avatar: ""}, nil
	}

	// Exchange code for access_token
	tokenValues := url.Values{}
	tokenValues.Set("appid", w.cfg.AppID)
	tokenValues.Set("secret", w.cfg.AppSecret)
	tokenValues.Set("code", code)
	tokenValues.Set("grant_type", "authorization_code")
	tokenURL := "https://api.weixin.qq.com/sns/oauth2/access_token?" + tokenValues.Encode()
	resp, err := w.client.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("wechat oauth token: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("wechat oauth token: status %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("wechat oauth token: invalid response")
	}
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat oauth: %d %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	// Fetch user info
	infoValues := url.Values{}
	infoValues.Set("access_token", tokenResp.AccessToken)
	infoValues.Set("openid", tokenResp.OpenID)
	infoValues.Set("lang", "zh_CN")
	infoURL := "https://api.weixin.qq.com/sns/userinfo?" + infoValues.Encode()
	resp2, err := w.client.Get(infoURL)
	if err != nil {
		return nil, fmt.Errorf("wechat oauth userinfo: %w", err)
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	if resp2.StatusCode < http.StatusOK || resp2.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("wechat oauth userinfo: status %d", resp2.StatusCode)
	}

	var user WechatUser
	if err := json.Unmarshal(body2, &user); err != nil {
		return nil, fmt.Errorf("wechat oauth userinfo: invalid response")
	}
	if user.OpenID == "" {
		return nil, fmt.Errorf("wechat oauth: empty openid")
	}
	return &user, nil
}
