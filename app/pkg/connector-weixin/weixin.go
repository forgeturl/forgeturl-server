package connector_weixin

import "github.com/sunmi-OS/gocore/v2/conf/viper"

// 网站应用微信登录开发指南 https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
// 授权后接口调用（UnionID） https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Authorized_Interface_Calling_UnionID.html

var Connector *ConnectorConfig

func Init() {
	Connector = &ConnectorConfig{
		ClientID:     viper.GetEnvConfig("WEIXIN_CLIENT_ID").String(),
		ClientSecret: viper.GetEnvConfig("WEIXIN_CLIENT_SECRET").String(),
	}
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
