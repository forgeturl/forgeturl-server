package connector_google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"forgeturl-server/api/common"

	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"golang.org/x/oauth2"
	oauth2Google "golang.org/x/oauth2/google"
)

var Connector *ConnectorConfig

func Init() {
	Connector = &ConnectorConfig{
		ClientID:     viper.GetEnvConfig("GOOGLE_CLIENT_ID").String(),
		ClientSecret: viper.GetEnvConfig("GOOGLE_CLIENT_SECRET").String(),
	}
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthUserInfo struct {
	Sub           string `json:"sub"`  // 用户ID，长度255个字符
	Name          string `json:"name"` // 展示名称
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"` // 头像地址
	Email         string `json:"email"`   // 邮箱地址
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (g *ConnectorConfig) ConnectorSender(ctx context.Context, receiverURL string) (redirectURL string) {
	oauth2Config := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		Endpoint:     oauth2Google.Endpoint,
		RedirectURL:  receiverURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
	}
	return oauth2Config.AuthCodeURL("state")
}

func (g *ConnectorConfig) ConnectorReceiver(ctx context.Context, code string) (userInfo *AuthUserInfo, err error) {
	oauth2Config := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		Endpoint:     oauth2Google.Endpoint,
		RedirectURL:  "https://api.forgeturl.com/login/connector/google",
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
	}

	client := oauth2Config.Client(ctx, token)
	client.Timeout = 15 * time.Second
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
	}

	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)

	respGoogleAuthUserInfo := &AuthUserInfo{}
	if err = json.Unmarshal(data, respGoogleAuthUserInfo); err != nil {
		return nil, common.ErrNotAuthenticated(fmt.Sprintf("parse google oauth user info response failed: %v", err))
	}

	// register id = respGoogleAuthUserInfo.Sub
	// dispalyName = respGoogleAuthUserInfo.Name
	return respGoogleAuthUserInfo, nil
}
