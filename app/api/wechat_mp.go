package api

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"forgeturl-server/dal"
	"forgeturl-server/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
)

const (
	weChatMPDefaultAPIBase   = "https://api.weixin.qq.com"
	weChatMPDefaultOAuthBase = "https://open.weixin.qq.com"
	weChatMPSecretHeader     = "X-AVM-WeChat-Secret"
)

type weChatMPConfig struct {
	AppID          string
	AppSecret      string
	TemplateID     string
	CallbackURL    string
	AVMCallbackURL string
	MessageURL     string
	BridgeSecret   string
	APIBase        string
	OAuthBase      string
}

type weChatMPTokenStore interface {
	GetAVMWeChatMPAccessToken(context.Context, string) (string, error)
	SetAVMWeChatMPAccessToken(context.Context, string, string, time.Duration) error
	DeleteAVMWeChatMPAccessToken(context.Context, string) error
}

type weChatMPClient struct {
	config     weChatMPConfig
	httpClient *http.Client
	tokenStore weChatMPTokenStore
}

type weChatAPIError struct {
	Code    int
	Message string
}

func (e *weChatAPIError) Error() string {
	return fmt.Sprintf("wechat api error %d: %s", e.Code, e.Message)
}

type weChatMPBindExchangeReq struct {
	ResultCode string `json:"result_code"`
}

type weChatMPSendReq struct {
	OpenID           string `json:"openid"`
	WorkOrderName    string `json:"work_order_name"`
	ProjectName      string `json:"project_name"`
	TriggerCondition string `json:"trigger_condition"`
	TriggerSource    string `json:"trigger_source"`
	TriggerTime      string `json:"trigger_time"`
}

func weChatMPConfigValue(viperKey string, envKey string) string {
	if val := strings.TrimSpace(viper.C.GetString(viperKey)); val != "" {
		return val
	}
	return strings.TrimSpace(os.Getenv(envKey))
}

func loadWeChatMPConfig() weChatMPConfig {
	apiBase := strings.TrimRight(
		weChatMPConfigValue("base.WECHAT_MP_API_BASE", "WECHAT_MP_API_BASE"),
		"/",
	)
	if apiBase == "" {
		apiBase = weChatMPDefaultAPIBase
	}
	oauthBase := strings.TrimRight(
		weChatMPConfigValue("base.WECHAT_MP_OAUTH_BASE", "WECHAT_MP_OAUTH_BASE"),
		"/",
	)
	if oauthBase == "" {
		oauthBase = weChatMPDefaultOAuthBase
	}
	return weChatMPConfig{
		AppID:          weChatMPConfigValue("keys.WECHAT_MP_APP_ID", "WECHAT_MP_APP_ID"),
		AppSecret:      weChatMPConfigValue("keys.WECHAT_MP_APP_SECRET", "WECHAT_MP_APP_SECRET"),
		TemplateID:     weChatMPConfigValue("keys.WECHAT_MP_TEMPLATE_ID", "WECHAT_MP_TEMPLATE_ID"),
		CallbackURL:    weChatMPConfigValue("base.WECHAT_MP_CALLBACK_URL", "WECHAT_MP_CALLBACK_URL"),
		AVMCallbackURL: weChatMPConfigValue("base.WECHAT_MP_AVM_CALLBACK_URL", "WECHAT_MP_AVM_CALLBACK_URL"),
		MessageURL:     weChatMPConfigValue("base.WECHAT_MP_MESSAGE_URL", "WECHAT_MP_MESSAGE_URL"),
		BridgeSecret:   weChatMPConfigValue("keys.AVM_WECHAT_MP_BRIDGE_SECRET", "AVM_WECHAT_MP_BRIDGE_SECRET"),
		APIBase:        apiBase,
		OAuthBase:      oauthBase,
	}
}

func (cfg weChatMPConfig) bindingReady() bool {
	return cfg.AppID != "" &&
		cfg.AppSecret != "" &&
		validWeChatMPHTTPSURL(cfg.CallbackURL) &&
		validWeChatMPHTTPSURL(cfg.AVMCallbackURL) &&
		cfg.BridgeSecret != ""
}

func (cfg weChatMPConfig) sendingReady() bool {
	return cfg.bindingReady() &&
		cfg.TemplateID != "" &&
		(cfg.MessageURL == "" || validWeChatMPHTTPSURL(cfg.MessageURL))
}

func requireWeChatMPBridgeSecret(g *gin.Context, configuredSecret string) bool {
	provided := g.GetHeader(weChatMPSecretHeader)
	if configuredSecret == "" {
		g.JSON(http.StatusServiceUnavailable, gin.H{
			"code": 0,
			"msg":  "wechat mp bridge secret is not configured",
		})
		return false
	}
	if subtle.ConstantTimeCompare([]byte(provided), []byte(configuredSecret)) != 1 {
		g.JSON(http.StatusForbidden, gin.H{"code": 0, "msg": "invalid bridge secret"})
		return false
	}
	return true
}

func AVMWeChatMPBind() gin.HandlerFunc {
	return func(g *gin.Context) {
		cfg := loadWeChatMPConfig()
		if !cfg.bindingReady() {
			g.String(http.StatusServiceUnavailable, "微信服务号绑定尚未配置")
			return
		}
		bindCode := strings.TrimSpace(g.Query("bind_code"))
		expiresRaw := strings.TrimSpace(g.Query("expires_at"))
		signature := strings.TrimSpace(g.Query("signature"))
		expiresAt, err := strconv.ParseInt(expiresRaw, 10, 64)
		now := time.Now().Unix()
		if bindCode == "" ||
			len(bindCode) > 256 ||
			err != nil ||
			expiresAt <= now ||
			expiresAt > now+30*60 ||
			!validWeChatMPBindSignature(
				bindCode,
				expiresAt,
				signature,
				cfg.BridgeSecret,
			) {
			g.String(http.StatusBadRequest, "无效的绑定请求")
			return
		}
		state := middleware.NewUUID()
		if err := dal.C.SetAVMWeChatMPBindState(
			g.Request.Context(), state, bindCode,
		); err != nil {
			g.String(http.StatusInternalServerError, "暂时无法创建绑定请求")
			return
		}
		params := url.Values{
			"appid":         {cfg.AppID},
			"redirect_uri":  {cfg.CallbackURL},
			"response_type": {"code"},
			"scope":         {"snsapi_base"},
			"state":         {state},
		}
		authURL := cfg.OAuthBase + "/connect/oauth2/authorize?" +
			params.Encode() + "#wechat_redirect"
		g.Redirect(http.StatusFound, authURL)
	}
}

func AVMWeChatMPCallback() gin.HandlerFunc {
	return func(g *gin.Context) {
		cfg := loadWeChatMPConfig()
		if !cfg.bindingReady() {
			g.String(http.StatusServiceUnavailable, "微信服务号绑定尚未配置")
			return
		}
		code := strings.TrimSpace(g.Query("code"))
		state := strings.TrimSpace(g.Query("state"))
		if code == "" || state == "" {
			g.String(http.StatusBadRequest, "微信授权参数缺失")
			return
		}
		bindCode, err := dal.C.ConsumeAVMWeChatMPBindState(
			g.Request.Context(), state,
		)
		if err != nil {
			g.String(http.StatusUnauthorized, "绑定请求已失效，请返回电脑重新扫码")
			return
		}
		client := &weChatMPClient{
			config:     cfg,
			httpClient: &http.Client{Timeout: 15 * time.Second},
			tokenStore: dal.C,
		}
		openID, err := client.exchangeOAuthCode(g.Request.Context(), code)
		if err != nil {
			g.String(http.StatusBadGateway, "微信授权失败，请稍后重试")
			return
		}
		resultCode := middleware.NewUUID()
		err = dal.C.SetAVMWeChatMPBindResult(
			g.Request.Context(),
			resultCode,
			dal.AVMWeChatMPBindResultPayload{
				BindCode: bindCode,
				OpenID:   openID,
			},
		)
		if err != nil {
			g.String(http.StatusInternalServerError, "暂时无法保存绑定结果")
			return
		}
		callbackURL, err := url.Parse(cfg.AVMCallbackURL)
		if err != nil || callbackURL.Scheme != "https" || callbackURL.Host == "" {
			g.String(http.StatusInternalServerError, "绑定回调地址配置无效")
			return
		}
		query := callbackURL.Query()
		query.Set("result_code", resultCode)
		callbackURL.RawQuery = query.Encode()
		g.Redirect(http.StatusFound, callbackURL.String())
	}
}

func AVMWeChatMPBindExchange() gin.HandlerFunc {
	return func(g *gin.Context) {
		cfg := loadWeChatMPConfig()
		if !requireWeChatMPBridgeSecret(g, cfg.BridgeSecret) {
			return
		}
		req := &weChatMPBindExchangeReq{}
		if err := g.ShouldBindJSON(req); err != nil ||
			strings.TrimSpace(req.ResultCode) == "" {
			g.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "missing result_code"})
			return
		}
		payload, err := dal.C.ConsumeAVMWeChatMPBindResult(
			g.Request.Context(), strings.TrimSpace(req.ResultCode),
		)
		if err != nil {
			g.JSON(http.StatusUnauthorized, gin.H{"code": 0, "msg": err.Error()})
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": gin.H{
				"bind_code": payload.BindCode,
				"openid":    payload.OpenID,
			},
		})
	}
}

func AVMWeChatMPSend() gin.HandlerFunc {
	return func(g *gin.Context) {
		cfg := loadWeChatMPConfig()
		if !cfg.sendingReady() {
			g.JSON(http.StatusServiceUnavailable, gin.H{
				"code": 0,
				"msg":  "wechat mp sending is not configured",
			})
			return
		}
		if !requireWeChatMPBridgeSecret(g, cfg.BridgeSecret) {
			return
		}
		req := &weChatMPSendReq{}
		if err := g.ShouldBindJSON(req); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "invalid request"})
			return
		}
		req.OpenID = strings.TrimSpace(req.OpenID)
		if req.OpenID == "" || len(req.OpenID) > 256 {
			g.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "invalid openid"})
			return
		}
		client := &weChatMPClient{
			config:     cfg,
			httpClient: &http.Client{Timeout: 15 * time.Second},
			tokenStore: dal.C,
		}
		msgID, err := client.sendTemplate(g.Request.Context(), *req)
		if err != nil {
			status := http.StatusBadGateway
			result := gin.H{"code": 0, "msg": err.Error()}
			var apiErr *weChatAPIError
			if errors.As(err, &apiErr) {
				result["error_code"] = apiErr.Code
			}
			g.JSON(status, result)
			return
		}
		g.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": gin.H{"msgid": msgID},
		})
	}
}

func (c *weChatMPClient) exchangeOAuthCode(
	ctx context.Context, code string,
) (string, error) {
	params := url.Values{
		"appid":      {c.config.AppID},
		"secret":     {c.config.AppSecret},
		"code":       {code},
		"grant_type": {"authorization_code"},
	}
	endpoint := c.config.APIBase + "/sns/oauth2/access_token?" + params.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var payload struct {
		OpenID  string `json:"openid"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := decodeWeChatResponse(response, &payload); err != nil {
		return "", err
	}
	if payload.ErrCode != 0 {
		return "", &weChatAPIError{Code: payload.ErrCode, Message: payload.ErrMsg}
	}
	if strings.TrimSpace(payload.OpenID) == "" {
		return "", errors.New("wechat oauth response did not include openid")
	}
	return strings.TrimSpace(payload.OpenID), nil
}

func (c *weChatMPClient) accessToken(
	ctx context.Context, forceRefresh bool,
) (string, error) {
	if !forceRefresh {
		cached, err := c.tokenStore.GetAVMWeChatMPAccessToken(ctx, c.config.AppID)
		if err != nil {
			return "", err
		}
		if cached != "" {
			return cached, nil
		}
	}
	body, err := json.Marshal(map[string]any{
		"grant_type":    "client_credential",
		"appid":         c.config.AppID,
		"secret":        c.config.AppSecret,
		"force_refresh": forceRefresh,
	})
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.config.APIBase+"/cgi-bin/stable_token",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var payload struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := decodeWeChatResponse(response, &payload); err != nil {
		return "", err
	}
	if payload.ErrCode != 0 {
		return "", &weChatAPIError{Code: payload.ErrCode, Message: payload.ErrMsg}
	}
	if payload.AccessToken == "" {
		return "", errors.New("wechat access token response was empty")
	}
	ttlSeconds := payload.ExpiresIn - 300
	if ttlSeconds < 60 {
		ttlSeconds = 60
	}
	if err := c.tokenStore.SetAVMWeChatMPAccessToken(
		ctx,
		c.config.AppID,
		payload.AccessToken,
		time.Duration(ttlSeconds)*time.Second,
	); err != nil {
		return "", err
	}
	return payload.AccessToken, nil
}

func (c *weChatMPClient) sendTemplate(
	ctx context.Context, req weChatMPSendReq,
) (int64, error) {
	token, err := c.accessToken(ctx, false)
	if err != nil {
		return 0, err
	}
	msgID, err := c.sendTemplateWithToken(ctx, token, req)
	var apiErr *weChatAPIError
	if !errors.As(err, &apiErr) || !isInvalidWeChatAccessToken(apiErr.Code) {
		return msgID, err
	}
	if deleteErr := c.tokenStore.DeleteAVMWeChatMPAccessToken(
		ctx, c.config.AppID,
	); deleteErr != nil {
		return 0, deleteErr
	}
	token, err = c.accessToken(ctx, true)
	if err != nil {
		return 0, err
	}
	return c.sendTemplateWithToken(ctx, token, req)
}

func (c *weChatMPClient) sendTemplateWithToken(
	ctx context.Context, token string, req weChatMPSendReq,
) (int64, error) {
	data := map[string]map[string]string{
		"thing4":  {"value": truncateWeChatTemplateValue(req.WorkOrderName, 20)},
		"thing10": {"value": truncateWeChatTemplateValue(req.ProjectName, 20)},
		"thing5":  {"value": truncateWeChatTemplateValue(req.TriggerCondition, 20)},
		"thing8":  {"value": truncateWeChatTemplateValue(req.TriggerSource, 20)},
		"time14":  {"value": truncateWeChatTemplateValue(req.TriggerTime, 20)},
	}
	payload := map[string]any{
		"touser":      req.OpenID,
		"template_id": c.config.TemplateID,
		"data":        data,
	}
	if c.config.MessageURL != "" {
		payload["url"] = c.config.MessageURL
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	endpoint := c.config.APIBase + "/cgi-bin/message/template/send?access_token=" +
		url.QueryEscape(token)
	request, err := http.NewRequestWithContext(
		ctx, http.MethodPost, endpoint, bytes.NewReader(body),
	)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MsgID   int64  `json:"msgid"`
	}
	if err := decodeWeChatResponse(response, &result); err != nil {
		return 0, err
	}
	if result.ErrCode != 0 {
		return 0, &weChatAPIError{Code: result.ErrCode, Message: result.ErrMsg}
	}
	return result.MsgID, nil
}

func decodeWeChatResponse(response *http.Response, target any) error {
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("wechat api returned HTTP %d", response.StatusCode)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode wechat api response failed: %w", err)
	}
	return nil
}

func truncateWeChatTemplateValue(value string, maxRunes int) string {
	value = strings.TrimSpace(value)
	if maxRunes <= 0 || utf8.RuneCountInString(value) <= maxRunes {
		return value
	}
	runes := []rune(value)
	return string(runes[:maxRunes])
}

func isInvalidWeChatAccessToken(code int) bool {
	switch code {
	case 40001, 40014, 42001:
		return true
	default:
		return false
	}
}

func validWeChatMPBindSignature(
	bindCode string, expiresAt int64, signature string, secret string,
) bool {
	if bindCode == "" || signature == "" || secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(fmt.Sprintf("%s.%d", bindCode, expiresAt)))
	expected := fmt.Sprintf("%x", mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(strings.ToLower(signature)))
}

func validWeChatMPHTTPSURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	return err == nil &&
		parsed.Scheme == "https" &&
		parsed.Host != "" &&
		parsed.User == nil
}
