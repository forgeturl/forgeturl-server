package connector_provider

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/markbates/goth/providers/google"
)

func TestGoogleAuthorizationUsesOfficialOriginWithServerProxies(t *testing.T) {
	originalEndpoint := google.Endpoint
	t.Cleanup(func() { google.Endpoint = originalEndpoint })

	provider := newGoogleProvider("test-client", "test-secret", "https://forgeturl.example/auth/callback/google",
		"accounts-proxy.example", "token-proxy.example", "profile-proxy.example")
	session, err := provider.BeginAuth("test-state")
	if err != nil {
		t.Fatal(err)
	}
	authURL, err := session.GetAuthURL()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" || parsed.Host != "accounts.google.com" {
		t.Fatalf("browser authorization must stay on Google's origin, got %s://%s", parsed.Scheme, parsed.Host)
	}
	for key, want := range map[string]string{
		"client_id": "test-client", "state": "test-state", "response_type": "code",
		"redirect_uri": "https://forgeturl.example/auth/callback/google", "scope": "openid email profile",
	} {
		if got := parsed.Query().Get(key); got != want {
			t.Errorf("authorization %s = %q, want %q", key, got, want)
		}
	}
}

func TestGoogleServerProxiesStillCompleteLogin(t *testing.T) {
	originalEndpoint := google.Endpoint
	t.Cleanup(func() { google.Endpoint = originalEndpoint })

	provider := newGoogleProvider("test-client", "test-secret", "https://forgeturl.example/auth/callback/google",
		"accounts-proxy.example", "token-proxy.example", "profile-proxy.example")
	var requests []string
	provider.HTTPClient = &http.Client{Transport: googleRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.String())
		var body string
		switch req.URL.String() {
		case "https://token-proxy.example/token":
			if err := req.ParseForm(); err != nil {
				return nil, err
			}
			if req.Form.Get("code") != "test-code" || req.Form.Get("client_id") != "test-client" || req.Form.Get("client_secret") != "test-secret" {
				return nil, fmt.Errorf("token exchange lost its authorization parameters")
			}
			body = `{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`
		case "https://profile-proxy.example/v1/userinfo":
			if req.Header.Get("Authorization") != "Bearer test-token" {
				return nil, fmt.Errorf("profile request lost its access token")
			}
			body = `{"sub":"test-google-subject","name":"Test User","email":"test@example.invalid"}`
		default:
			return nil, fmt.Errorf("unexpected Google request: %s", req.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}
	session, err := provider.BeginAuth("test-state")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = session.Authorize(provider, url.Values{"code": {"test-code"}}); err != nil {
		t.Fatal(err)
	}
	user, err := provider.FetchUser(session)
	if err != nil {
		t.Fatal(err)
	}
	if user.UserID != "test-google-subject" || user.Email != "test@example.invalid" {
		t.Fatalf("unexpected authenticated Google identity: id=%q, email=%q", user.UserID, user.Email)
	}
	want := []string{"POST https://token-proxy.example/token", "GET https://profile-proxy.example/v1/userinfo"}
	if !reflect.DeepEqual(requests, want) {
		t.Fatalf("server requests = %v, want %v", requests, want)
	}
}

type googleRoundTripFunc func(*http.Request) (*http.Response, error)

func (f googleRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
