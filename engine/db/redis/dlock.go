package redis

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var cache *redis.Client

const redisMutexLockExpTime = 15

// TryGetDistributedLock 分布式锁获取
// requestId 用于标识请求客户端，可以是随机字符串，需确保唯一
func TryGetDistributedLock(c context.Context, lockKey, requestId string, isNegative bool) bool {
	if isNegative { // 多次尝试获取
		retry := 1
		for {
			cmd := cache.Do(c, "SET", lockKey, requestId, "EX", redisMutexLockExpTime, "NX")
			// 获取锁成功
			if cmd != nil && cmd.Err() == nil {
				return true
			}
			// 尝试多次没获取成功
			if retry > 10 {
				return false
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			retry += 1
		}
	} else { // 只尝试一次
		cmd := cache.Do(c, "SET", lockKey, requestId, "EX", redisMutexLockExpTime, "NX")
		// 获取锁成功
		if cmd != nil && cmd.Err() == nil {
			return true
		}

		return false
	}
}

// ReleaseDistributedLock 释放锁，通过比较requestId，用于确保只释放自己的锁，使用lua脚本保证操作的原子型
func ReleaseDistributedLock(c context.Context, lockKey, requestId string) (bool, error) {
	luaScript := `
    if redis.call("get",KEYS[1]) == ARGV[1]
    then
        return redis.call("del",KEYS[1])
    else
        return 0
    end`

	cmd := cache.Do(c, "eval", luaScript, 1, lockKey, requestId)

	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	do, err := cmd.Int()
	if err != nil {
		return false, err
	}

	if do == 1 {
		return true, err
	}

	return false, err

}

// DLock 获取分布式锁，成功需要释放锁
func DLock(c context.Context, lockKey, requestId string) (releaseLock func(), err error) {

	lockKey = fmt.Sprintf("dlockKey:%s", lockKey)
	requestId = fmt.Sprintf("dlockReqID:%s", requestId)
	lockOk := TryGetDistributedLock(c, lockKey, requestId, true)
	if !lockOk {
		return nil, errors.New("get dlock failed")
	}
	releaseLock = func() {
		_, _ = ReleaseDistributedLock(c, lockKey, requestId)
	}
	return releaseLock, nil
}
