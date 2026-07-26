package dal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"forgeturl-server/api/common"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	gredis "github.com/sunmi-OS/gocore/v2/db/redis"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/lib/ratelimiter"
)

const (
	// LoginTimeout 登录过期时间
	LoginTimeout                    = time.Hour * 24 * 180
	AVMAuthCodeTimeout              = time.Minute * 5
	AVMWeChatMPBindStateTimeout     = time.Minute * 10
	AVMWeChatMPBindResultTimeout    = time.Minute * 5
	RdsTokenPrefix                  = "auth:tk"
	RdsAVMAuthCodePrefix            = "auth:avm:code"
	RdsAVMWeChatMPBindStatePrefix   = "auth:avm:wechat-mp:state"
	RdsAVMWeChatMPBindResultPrefix  = "auth:avm:wechat-mp:result"
	RdsAVMWeChatMPAccessTokenPrefix = "auth:avm:wechat-mp:access-token"
)

type AVMAuthCodePayload struct {
	Provider    string `json:"provider"`
	WechatUID   string `json:"wechat_uid"`
	ForgetURLID int64  `json:"forgeturl_id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
}

type AVMWeChatMPBindResultPayload struct {
	BindCode string `json:"bind_code"`
	OpenID   string `json:"openid"`
}

type cacheImpl struct {
	user *redis.Client
	lock *redis.Client
}

var C *cacheImpl

var lockRateLimiter *ratelimiter.RedisRateLimiter

func initRedis() {
	C = &cacheImpl{}
	initList := []struct {
		client     **redis.Client
		redisNoKey string
	}{
		{client: &C.user, redisNoKey: "redisServer.userCache"},
		{client: &C.lock, redisNoKey: "redisServer.lock"},
	}

	// 初始化client
	for _, v := range initList {
		err := gredis.NewOrUpdateRedis(v.redisNoKey)
		if err != nil {
			panic(fmt.Errorf("init cache(%v) failed, err: %v", v.redisNoKey, err))
		}
		*v.client = gredis.GetRedis(v.redisNoKey)
		if *v.client == nil {
			panic(fmt.Errorf("get cache(%v) failed", v.redisNoKey))
		}

		glog.InfoF("init redis %v success", v.redisNoKey)
	}

	var err error
	rate := "1000-S"
	lockRateLimiter, err = ratelimiter.NewRedisRateLimiter(C.lock, ratelimiter.RedisConfig{
		Rate:   rate,
		Prefix: "redisRateLimiter",
	})
	if err != nil {
		panic(err)
	}
}

func (c *cacheImpl) GetXToken(ctx context.Context, key string) int64 {
	if key == "" {
		return 0
	}
	tk := GetTokenKey(key)
	val, err := c.user.Get(ctx, tk).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0
		}
		glog.WarnC(ctx, "get x-token failed, key: %s, err: %v", key, err)
		return 0
	}
	return cast.ToInt64(val)
}

func (c *cacheImpl) SetXToken(ctx context.Context, key string, uid int64) error {
	if key == "" || uid == 0 {
		return common.ErrInternalServerError("invalid key or uid")
	}
	tk := GetTokenKey(key)
	err := c.user.Set(ctx, tk, uid, LoginTimeout).Err()
	if err != nil {
		return common.ErrInternalServerError(fmt.Sprintf("set x-token failed, key: %s, uid: %d, err: %v", key, uid, err))
	}
	glog.InfoC(ctx, "set x-token, key: %s, uid: %d", key, uid)
	return nil
}

// DelXToken 删除token
func (c *cacheImpl) DelXToken(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	tk := GetTokenKey(key)
	err := c.user.Del(ctx, tk).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		glog.WarnC(ctx, "del x-token failed, key: %s, err: %v", tk, err)
		return common.ErrInternalServerError("del x-token failed")
	}
	glog.InfoC(ctx, "del x-token, key: %s", key)
	return nil
}

func GetTokenKey(key string) string {
	return RdsTokenPrefix + ":" + key
}

func (c *cacheImpl) SetAVMAuthCode(ctx context.Context, code string, payload AVMAuthCodePayload) error {
	if code == "" || payload.WechatUID == "" {
		return common.ErrInternalServerError("invalid avm auth code payload")
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return common.ErrInternalServerError(fmt.Sprintf("marshal avm auth code failed: %v", err))
	}
	if err := c.user.Set(ctx, GetAVMAuthCodeKey(code), string(buf), AVMAuthCodeTimeout).Err(); err != nil {
		return common.ErrInternalServerError(fmt.Sprintf("set avm auth code failed, err: %v", err))
	}
	return nil
}

func (c *cacheImpl) ConsumeAVMAuthCode(ctx context.Context, code string) (*AVMAuthCodePayload, error) {
	if code == "" {
		return nil, common.ErrBadRequest("missing auth_code")
	}
	key := GetAVMAuthCodeKey(code)
	val, err := c.user.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.ErrNotAuthenticated("invalid or expired auth_code")
		}
		return nil, common.ErrInternalServerError(fmt.Sprintf("get avm auth code failed, err: %v", err))
	}
	if err := c.user.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return nil, common.ErrInternalServerError(fmt.Sprintf("delete avm auth code failed, err: %v", err))
	}
	payload := &AVMAuthCodePayload{}
	if err := json.Unmarshal([]byte(val), payload); err != nil {
		return nil, common.ErrInternalServerError(fmt.Sprintf("unmarshal avm auth code failed, err: %v", err))
	}
	return payload, nil
}

func GetAVMAuthCodeKey(code string) string {
	return RdsAVMAuthCodePrefix + ":" + code
}

func (c *cacheImpl) SetAVMWeChatMPBindState(
	ctx context.Context, state string, bindCode string,
) error {
	if state == "" || bindCode == "" {
		return common.ErrBadRequest("missing bind state or bind code")
	}
	err := c.user.Set(
		ctx,
		GetAVMWeChatMPBindStateKey(state),
		bindCode,
		AVMWeChatMPBindStateTimeout,
	).Err()
	if err != nil {
		return common.ErrInternalServerError(
			fmt.Sprintf("set wechat mp bind state failed: %v", err),
		)
	}
	return nil
}

func (c *cacheImpl) ConsumeAVMWeChatMPBindState(
	ctx context.Context, state string,
) (string, error) {
	if state == "" {
		return "", common.ErrBadRequest("missing bind state")
	}
	val, err := c.user.GetDel(ctx, GetAVMWeChatMPBindStateKey(state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", common.ErrNotAuthenticated("invalid or expired bind state")
		}
		return "", common.ErrInternalServerError(
			fmt.Sprintf("consume wechat mp bind state failed: %v", err),
		)
	}
	return val, nil
}

func (c *cacheImpl) SetAVMWeChatMPBindResult(
	ctx context.Context, resultCode string, payload AVMWeChatMPBindResultPayload,
) error {
	if resultCode == "" || payload.BindCode == "" || payload.OpenID == "" {
		return common.ErrBadRequest("invalid wechat mp bind result")
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return common.ErrInternalServerError(
			fmt.Sprintf("marshal wechat mp bind result failed: %v", err),
		)
	}
	err = c.user.Set(
		ctx,
		GetAVMWeChatMPBindResultKey(resultCode),
		string(buf),
		AVMWeChatMPBindResultTimeout,
	).Err()
	if err != nil {
		return common.ErrInternalServerError(
			fmt.Sprintf("set wechat mp bind result failed: %v", err),
		)
	}
	return nil
}

func (c *cacheImpl) ConsumeAVMWeChatMPBindResult(
	ctx context.Context, resultCode string,
) (*AVMWeChatMPBindResultPayload, error) {
	if resultCode == "" {
		return nil, common.ErrBadRequest("missing result_code")
	}
	val, err := c.user.GetDel(
		ctx,
		GetAVMWeChatMPBindResultKey(resultCode),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.ErrNotAuthenticated("invalid or expired result_code")
		}
		return nil, common.ErrInternalServerError(
			fmt.Sprintf("consume wechat mp bind result failed: %v", err),
		)
	}
	payload := &AVMWeChatMPBindResultPayload{}
	if err := json.Unmarshal([]byte(val), payload); err != nil {
		return nil, common.ErrInternalServerError(
			fmt.Sprintf("unmarshal wechat mp bind result failed: %v", err),
		)
	}
	return payload, nil
}

func (c *cacheImpl) GetAVMWeChatMPAccessToken(
	ctx context.Context, appID string,
) (string, error) {
	if appID == "" {
		return "", common.ErrBadRequest("missing wechat mp app id")
	}
	val, err := c.user.Get(ctx, GetAVMWeChatMPAccessTokenKey(appID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", common.ErrInternalServerError(
			fmt.Sprintf("get wechat mp access token failed: %v", err),
		)
	}
	return val, nil
}

func (c *cacheImpl) SetAVMWeChatMPAccessToken(
	ctx context.Context, appID string, token string, ttl time.Duration,
) error {
	if appID == "" || token == "" || ttl <= 0 {
		return common.ErrBadRequest("invalid wechat mp access token")
	}
	err := c.user.Set(
		ctx,
		GetAVMWeChatMPAccessTokenKey(appID),
		token,
		ttl,
	).Err()
	if err != nil {
		return common.ErrInternalServerError(
			fmt.Sprintf("set wechat mp access token failed: %v", err),
		)
	}
	return nil
}

func (c *cacheImpl) DeleteAVMWeChatMPAccessToken(
	ctx context.Context, appID string,
) error {
	if appID == "" {
		return nil
	}
	err := c.user.Del(ctx, GetAVMWeChatMPAccessTokenKey(appID)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return common.ErrInternalServerError(
			fmt.Sprintf("delete wechat mp access token failed: %v", err),
		)
	}
	return nil
}

func GetAVMWeChatMPBindStateKey(state string) string {
	return RdsAVMWeChatMPBindStatePrefix + ":" + state
}

func GetAVMWeChatMPBindResultKey(resultCode string) string {
	return RdsAVMWeChatMPBindResultPrefix + ":" + resultCode
}

func GetAVMWeChatMPAccessTokenKey(appID string) string {
	return RdsAVMWeChatMPAccessTokenPrefix + ":" + appID
}
