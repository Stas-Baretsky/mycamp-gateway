package bucketlimiter

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	DECINE        int = 0
	APPROVE       int = 1
	CLEAROLDTOKEN int = 2
)

type RateLimiterStorage interface {
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration)
	Read(ctx context.Context, key string) (int64, error)
}

type LimiterSettings struct {
	Speed      uint
	Limit      uint
	CallWeight uint
	TokenLife  time.Duration
}

type Limiter interface {
	Allow(key string) bool
	Stop()
}

type BucketLimiter struct {
	settings *LimiterSettings
	storage  *RateLimiterStorage
	wg       sync.WaitGroup
	m        sync.Mutex
	stopped  bool
}

func MakeBucketLimiter() BucketLimiter {
	return BucketLimiter{
		stopped: false,
	}
}

func (bl *BucketLimiter) Init(
	storage *RateLimiterStorage,
	settings *LimiterSettings,
) {
	bl.settings = settings
	bl.storage = storage

	bl.AsyncBucketUpdate()
	bl.controlTokenLife()
}

func (bl *BucketLimiter) Stop() {
	bl.m.Lock()
	bl.stopped = true
	bl.m.Unlock()
	bl.wg.Wait()
}

func (bl *BucketLimiter) AsyncBucketUpdate() {
	bl.wg.Add(1)
	go func() {
		defer bl.wg.Done()

		ticker := time.NewTicker(time.Second / time.Duration(bl.speed))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				bl.m.Lock()
				if bl.stopped {
					bl.m.Unlock()
					return
				}
				for token, cap := range bl.cap {
					if cap < bl.settings.Limit {
						bl.cap[token] = cap + 1
					}
				}
				bl.m.Unlock()
			}
		}
	}()
}

func (bl *BucketLimiter) Allow(token string) bool {
	bl.m.Lock()
	defer bl.m.Unlock()

	//Init new client
	if _, ok := bl.cap[token]; !ok {
		bl.cap[token] = bl.limit
		bl.lastTokenCalls[token] = time.Now()
		fmt.Printf("🆕 Новый клиент %s инициализирован с %d токенами\n",
			token, bl.limit)
	}

	if bl.cap[token] < bl.callWeight {
		bl.logDecision(token, DECINE)
		return false
	}

	bl.logDecision(token, APPROVE)
	bl.cap[token] -= bl.callWeight
	bl.lastTokenCalls[token] = time.Now()

	return true

}

func (bl *BucketLimiter) controlTokenLife() {
	bl.wg.Add(1)

	go func() {
		defer bl.wg.Done()
		bl.clearOldToken()
	}()
}

func (bl *BucketLimiter) clearOldToken() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bl.m.Lock()
			if bl.stopped {
				bl.m.Unlock()
				return
			}
			for token, callTime := range bl.lastTokenCalls {
				if time.Since(callTime) > bl.tokenLife {
					bl.logDecision(token, CLEAROLDTOKEN)
					delete(bl.cap, token)
					delete(bl.lastTokenCalls, token)
				}
			}
			bl.m.Unlock()
		}
	}
}

func randomWorks(bl *BucketLimiter, token string) {
	bl.wg.Add(1)
	go func() {
		defer bl.wg.Done()
		for {
			if bl.stopped {
				return
			}
			time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
			bl.Allow(token)
		}
	}()
}

func (bl *BucketLimiter) logDecision(token string, message int) {
	currentCap := bl.cap[token]
	fmt.Println("____________________________\n")
	if message == APPROVE {
		fmt.Printf("✅ РАЗРЕШЕНО для %s\n", token)
	}
	if message == DECINE {
		fmt.Printf("❌ ОТКАЗАНО для %s\n", token)
	}
	if message == CLEAROLDTOKEN {
		fmt.Printf("🧹 Очистка просроченных токенов для %s\n",
			token)
	}
	fmt.Printf("Токенов: %d / %d\n", currentCap, bl.limit)
	fmt.Println("____________________________\n")
}
