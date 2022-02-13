package rate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

const (
	expireDefault = 30 * time.Minute
	cleanupPeriod = 1 * time.Hour
)

// lCache keeps limiters for each custom key (e.g. IP address or other unique key to rate on)
var lCache = cache.New(expireDefault, cleanupPeriod)

type LimiterFunc = func(*gin.Context) (*rate.Limiter, time.Duration)

// ByIP creates a limiter that allows only `numRequests` every `per` duration based on caller's IP
func ByIP(numRequests int, per time.Duration) gin.HandlerFunc {
	var getIp = func(c *gin.Context) string {
		return c.ClientIP()
	}
	return NewRateLimiter(getIp, limiter(numRequests, per), abortWithTooManyRequests)
}

func ByCustomKey(numRequests int, per time.Duration, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return NewRateLimiter(keyFunc, limiter(numRequests, per), abortWithTooManyRequests)
}

// limiter guesses the correct expire time for a limiter - make it greater rather than smaller
// to be on the safe side
func limiter(numRequests int, per time.Duration) LimiterFunc {
	var expire = 2 * per
	if expire < expireDefault {
		expire = expireDefault
	}
	return func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(per), numRequests), expire
	}
}

// abortWithTooManyRequests aborts with http status code TooManyReuqests (429)
func abortWithTooManyRequests(c *gin.Context) {
	c.AbortWithStatus(http.StatusTooManyRequests)
}

// NewRateLimiter is here for your complete control but there are sufficient defaults such as OnIP
func NewRateLimiter(
	keyFunc func(*gin.Context) string,
	createLimiter LimiterFunc,
	abort func(*gin.Context),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := keyFunc(c)
		limiter, ok := lCache.Get(k)
		if !ok {
			var expire time.Duration
			limiter, expire = createLimiter(c)
			lCache.Set(k, limiter, expire)
		}
		ok = limiter.(*rate.Limiter).Allow()
		if !ok {
			abort(c)
			return
		}
		c.Next()
	}
}
