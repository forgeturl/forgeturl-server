package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type fakeWeChatMPTokenStore struct {
	mu      sync.Mutex
	token   string
	ttl     time.Duration
	deletes int
}

func (s *fakeWeChatMPTokenStore) GetAVMWeChatMPAccessToken(
	_ context.Context, _ string,
) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.token, nil
}

func (s *fakeWeChatMPTokenStore) SetAVMWeChatMPAccessToken(
	_ context.Context, _ string, token string, ttl time.Duration,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
	s.ttl = ttl
	return nil
}

func (s *fakeWeChatMPTokenStore) DeleteAVMWeChatMPAccessToken(
	_ context.Context, _ string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
	s.deletes++
	return nil
}

func TestWeChatMPSendRefreshesInvalidTokenAndMapsTemplateFields(t *testing.T) {
	var stableTokenCalls int
	var sendCalls int
	var sentPayload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter, request *http.Request,
	) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/cgi-bin/stable_token":
			stableTokenCalls++
			var payload map[string]any
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["appid"] != "mp-app" || payload["secret"] != "mp-secret" {
				t.Fatalf("unexpected token payload: %#v", payload)
			}
			_, _ = writer.Write([]byte(
				`{"access_token":"fresh-token","expires_in":7200}`,
			))
		case "/cgi-bin/message/template/send":
			sendCalls++
			if err := json.NewDecoder(request.Body).Decode(&sentPayload); err != nil {
				t.Fatal(err)
			}
			if request.URL.Query().Get("access_token") == "stale-token" {
				_, _ = writer.Write([]byte(`{"errcode":40014,"errmsg":"invalid token"}`))
				return
			}
			_, _ = writer.Write([]byte(`{"errcode":0,"errmsg":"ok","msgid":12345}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	store := &fakeWeChatMPTokenStore{token: "stale-token"}
	client := &weChatMPClient{
		config: weChatMPConfig{
			AppID:      "mp-app",
			AppSecret:  "mp-secret",
			TemplateID: "template-id",
			MessageURL: "https://lixiaoyaoai.com/hotspots",
			APIBase:    server.URL,
		},
		httpClient: server.Client(),
		tokenStore: store,
	}
	msgID, err := client.sendTemplate(context.Background(), weChatMPSendReq{
		OpenID:           "openid-1",
		WorkOrderName:    "这是一个超过二十个汉字需要被安全截断的热点标题示例",
		ProjectName:      "AI 创业",
		TriggerCondition: "评分88分，阈值75分",
		TriggerSource:    "微博",
		TriggerTime:      "2026-07-27 12:30:00",
	})
	if err != nil {
		t.Fatal(err)
	}
	if msgID != 12345 || stableTokenCalls != 1 || sendCalls != 2 {
		t.Fatalf(
			"unexpected result: msgID=%d tokenCalls=%d sendCalls=%d",
			msgID,
			stableTokenCalls,
			sendCalls,
		)
	}
	if store.deletes != 1 || store.token != "fresh-token" ||
		store.ttl != 6900*time.Second {
		t.Fatalf("unexpected token store state: %#v", store)
	}
	if sentPayload["touser"] != "openid-1" ||
		sentPayload["template_id"] != "template-id" ||
		sentPayload["url"] != "https://lixiaoyaoai.com/hotspots" {
		t.Fatalf("unexpected template payload: %#v", sentPayload)
	}
	data := sentPayload["data"].(map[string]any)
	thing4 := data["thing4"].(map[string]any)["value"].(string)
	if got := len([]rune(thing4)); got != 20 {
		t.Fatalf("thing4 length = %d, want 20", got)
	}
	if data["thing10"].(map[string]any)["value"] != "AI 创业" {
		t.Fatalf("unexpected project name: %#v", data["thing10"])
	}
}

func TestWeChatMPOAuthCodeExchange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter, request *http.Request,
	) {
		if request.URL.Path != "/sns/oauth2/access_token" {
			http.NotFound(writer, request)
			return
		}
		query := request.URL.Query()
		if query.Get("appid") != "mp-app" ||
			query.Get("secret") != "mp-secret" ||
			query.Get("code") != "oauth-code" ||
			query.Get("grant_type") != "authorization_code" {
			t.Fatalf("unexpected oauth query: %s", request.URL.RawQuery)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"openid":"service-openid","scope":"snsapi_base"}`))
	}))
	defer server.Close()

	client := &weChatMPClient{
		config: weChatMPConfig{
			AppID:     "mp-app",
			AppSecret: "mp-secret",
			APIBase:   server.URL,
		},
		httpClient: server.Client(),
	}
	openID, err := client.exchangeOAuthCode(context.Background(), "oauth-code")
	if err != nil {
		t.Fatal(err)
	}
	if openID != "service-openid" {
		t.Fatalf("openid = %q", openID)
	}
}

func TestWeChatMPBindSignature(t *testing.T) {
	if !validWeChatMPBindSignature(
		"bind-code",
		1_800_000_000,
		"4f7ccd238d5cafe77255ac6141f4ef262f01dfd7c19337b028431e4c0c6e6312",
		"bridge-secret",
	) {
		t.Fatal("valid signature was rejected")
	}
	if validWeChatMPBindSignature(
		"tampered",
		1_800_000_000,
		"4f7ccd238d5cafe77255ac6141f4ef262f01dfd7c19337b028431e4c0c6e6312",
		"bridge-secret",
	) {
		t.Fatal("tampered bind code was accepted")
	}
}
