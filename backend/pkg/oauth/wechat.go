package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// WechatOAuthConfig holds WeChat OAuth configuration.
type WechatOAuthConfig struct {
	Enabled     bool   `yaml:"enabled"`
	AppID       string `yaml:"app_id"`
	AppSecret   string `yaml:"app_secret"`
	RedirectURI string `yaml:"redirect_uri"`
}

// WechatUser represents a WeChat user profile.
type WechatUser struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"headimgurl"`
}

// WechatOAuth implements WeChat OAuth 2.0 login.
type WechatOAuth struct {
	cfg WechatOAuthConfig
}

func NewWechatOAuth(cfg WechatOAuthConfig) *WechatOAuth {
	return &WechatOAuth{cfg: cfg}
}

// GetAuthURL generates the WeChat authorization URL.
func (w *WechatOAuth) GetAuthURL(state string) string {
	if !w.cfg.Enabled {
		return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=MOCK_APPID&redirect_uri=MOCK_REDIRECT&response_type=code&scope=snsapi_userinfo&state=" + state
	}
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		w.cfg.AppID, url.QueryEscape(w.cfg.RedirectURI), state)
}

// GetUserInfo exchanges code for access_token and fetches user info.
func (w *WechatOAuth) GetUserInfo(code string) (*WechatUser, error) {
	if !w.cfg.Enabled {
		return &WechatUser{OpenID: "mock_openid", Nickname: "wx_mock_user", Avatar: ""}, nil
	}

	// Exchange code for access_token
	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		w.cfg.AppID, w.cfg.AppSecret, code)
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("wechat oauth token: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	json.Unmarshal(body, &tokenResp)
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat oauth: %d %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	// Fetch user info
	infoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		tokenResp.AccessToken, tokenResp.OpenID)
	resp2, err := http.Get(infoURL)
	if err != nil {
		return nil, fmt.Errorf("wechat oauth userinfo: %w", err)
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)

	var user WechatUser
	json.Unmarshal(body2, &user)
	if user.OpenID == "" {
		return nil, fmt.Errorf("wechat oauth: empty openid")
	}
	return &user, nil
}
