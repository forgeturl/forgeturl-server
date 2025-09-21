package dal

import (
	"context"
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
	LoginTimeout = time.Hour * 24 * 180
)

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
	val, err := c.user.Get(ctx, key).Result()
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
	err := c.user.Set(ctx, key, uid, LoginTimeout).Err()
	if err != nil {
		return common.ErrInternalServerError(fmt.Sprintf("set x-token failed, key: %s, uid: %d, err: %v", key, uid, err))
	}
	return nil
}

// DelXToken 删除token
func (c *cacheImpl) DelXToken(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	err := c.user.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		glog.WarnC(ctx, "del x-token failed, key: %s, err: %v", key, err)
		return common.ErrInternalServerError("del x-token failed")
	}
	return nil
}
