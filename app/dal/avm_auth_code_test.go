package dal

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAVMAuthCodeAtomicConsume(t *testing.T) {
	addr := os.Getenv("AVM_TEST_REDIS_ADDR")
	if !strings.HasPrefix(addr, "127.0.0.1:") {
		t.Skip("set AVM_TEST_REDIS_ADDR to an isolated localhost Redis")
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	defer client.Close()
	cache := &cacheImpl{user: client}
	ctx := context.Background()
	code := fmt.Sprintf("avm-google-test-%d", time.Now().UnixNano())
	defer client.Del(ctx, GetAVMAuthCodeKey(code))
	err := cache.SetAVMAuthCode(ctx, code, AVMAuthCodePayload{Provider: "google", Subject: "stable-google-sub"})
	if err != nil {
		t.Fatal(err)
	}
	var successes atomic.Int32
	var wait sync.WaitGroup
	for i := 0; i < 12; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			payload, err := cache.ConsumeAVMAuthCode(ctx, code)
			if err == nil {
				successes.Add(1)
				if payload.Subject != "stable-google-sub" || payload.WechatUID != "" {
					t.Errorf("unexpected identity: %+v", payload)
				}
			}
		}()
	}
	wait.Wait()
	if successes.Load() != 1 {
		t.Fatalf("expected one successful exchange, got %d", successes.Load())
	}
	if _, err := cache.ConsumeAVMAuthCode(ctx, code); err == nil {
		t.Fatal("code replay accepted")
	}
}
