package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"github.com/markbates/goth"
	"strings"
)

// Keep the AVM exchange contract separate from ordinary forgeturl sessions.
// Only the authenticated provider subject is an account identifier, never email.
func avmIdentity(user goth.User, profile *model.User) (dal.AVMAuthCodePayload, error) {
	if user.Provider != "wechat" && user.Provider != "google" {
		return dal.AVMAuthCodePayload{}, common.ErrBadRequest("unsupported AVM login provider")
	}
	if user.UserID == "" || user.UserID != strings.TrimSpace(user.UserID) || len(user.UserID) > 255 {
		return dal.AVMAuthCodePayload{}, common.ErrBadRequest("invalid AVM login subject")
	}
	result := dal.AVMAuthCodePayload{
		Provider: user.Provider, Subject: user.UserID, ForgetURLID: profile.ID,
		DisplayName: profile.DisplayName, Username: profile.Username,
		Avatar: profile.Avatar, Email: profile.Email,
	}
	if user.Provider == "wechat" {
		result.WechatUID = user.UserID
	}
	return result, nil
}
