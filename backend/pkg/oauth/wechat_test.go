package oauth

import (
	"net/url"
	"testing"
)

func TestWechatOAuthGetAuthURLUsesWebsiteQRConnect(t *testing.T) {
	oauth := NewWechatOAuth(WechatOAuthConfig{
		Enabled:     true,
		AppID:       "wx_test_app",
		RedirectURI: "https://www.zioran.com/auth/wechat/callback",
	})

	authURL := oauth.GetAuthURL("login")
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("parse auth url: %v", err)
	}

	if parsed.Scheme != "https" || parsed.Host != "open.weixin.qq.com" || parsed.Path != "/connect/qrconnect" {
		t.Fatalf("unexpected auth endpoint: %s", authURL)
	}
	query := parsed.Query()
	if query.Get("appid") != "wx_test_app" {
		t.Fatalf("unexpected appid: %s", query.Get("appid"))
	}
	if query.Get("redirect_uri") != "https://www.zioran.com/auth/wechat/callback" {
		t.Fatalf("unexpected redirect_uri: %s", query.Get("redirect_uri"))
	}
	if query.Get("scope") != "snsapi_login" {
		t.Fatalf("unexpected scope: %s", query.Get("scope"))
	}
	if query.Get("state") != "login" {
		t.Fatalf("unexpected state: %s", query.Get("state"))
	}
}

func TestWechatOAuthFrontendCallbackURL(t *testing.T) {
	oauth := NewWechatOAuth(WechatOAuthConfig{
		FrontendRedirectURI: "https://www.zioran.com/auth/wechat/callback?from=backend",
	})

	callbackURL := oauth.GetFrontendCallbackURL("code123", "login")
	parsed, err := url.Parse(callbackURL)
	if err != nil {
		t.Fatalf("parse callback url: %v", err)
	}
	query := parsed.Query()
	if query.Get("from") != "backend" || query.Get("code") != "code123" || query.Get("state") != "login" {
		t.Fatalf("unexpected callback url: %s", callbackURL)
	}
}
