package api

import (
	"forgeturl-server/dal/model"
	"github.com/markbates/goth"
	"testing"
)

func TestAVMIdentityProviders(t *testing.T) {
	profile := &model.User{ID: 42, DisplayName: "Example", Email: "same@example.test"}
	for _, provider := range []string{"wechat", "google"} {
		identity, err := avmIdentity(goth.User{Provider: provider, UserID: "stable-id"}, profile)
		if err != nil || identity.Subject != "stable-id" || identity.Provider != provider || identity.ForgetURLID != 42 {
			t.Fatalf("unexpected %s identity: %+v, %v", provider, identity, err)
		}
		if provider == "wechat" && identity.WechatUID != "stable-id" {
			t.Fatal("must preserve legacy WeChat UID")
		}
		if provider == "google" && identity.WechatUID != "" {
			t.Fatal("Google must not impersonate a WeChat identity")
		}
	}
}

func TestAVMIdentityRejectsUntrustedSubjects(t *testing.T) {
	for _, user := range []goth.User{
		{Provider: "github", UserID: "subject"},
		{Provider: "google", Email: "same@example.test"},
		{Provider: "google", UserID: " subject "},
	} {
		if _, err := avmIdentity(user, &model.User{}); err == nil {
			t.Fatal("must reject unsupported provider or missing subject")
		}
	}
}
