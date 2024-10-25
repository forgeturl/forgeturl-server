package connector_google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

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
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (g *ConnectorConfig) ConnectorReceiver(ctx context.Context, code string) (userInfo *AuthUserInfo, err error) {
	oauth2Config := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		Endpoint:     oauth2Google.Endpoint,
		RedirectURL:  "https://api.2049links.com/login/connector/google",
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return userInfo, err
	}

	client := oauth2Config.Client(context.TODO(), token)
	client.Timeout = 15 * time.Second
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return userInfo, err
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)

	respGoogleAuthUserInfo := &AuthUserInfo{}
	if err = json.Unmarshal(data, respGoogleAuthUserInfo); err != nil {
		return userInfo, fmt.Errorf("parse google oauth user info response failed: %v", err)
	}

	return respGoogleAuthUserInfo, nil
}
