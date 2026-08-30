package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type IPLimiter struct {
	mu        sync.Mutex
	perSecond rate.Limit
	burst     int
	visitors  map[string]*visitor
}

func NewIPLimiter(perMinute, burst int) *IPLimiter {
	if perMinute < 1 {
		perMinute = 1
	}
	if burst < 1 {
		burst = 1
	}
	return &IPLimiter{perSecond: rate.Limit(float64(perMinute) / 60), burst: burst, visitors: map[string]*visitor{}}
}
func (l *IPLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	v := l.visitors[ip]
	if v == nil {
		if len(l.visitors) >= 10000 {
			for key, item := range l.visitors {
				if now.Sub(item.lastSeen) > 10*time.Minute {
					delete(l.visitors, key)
				}
			}
		}
		v = &visitor{limiter: rate.NewLimiter(l.perSecond, l.burst)}
		l.visitors[ip] = v
	}
	v.lastSeen = now
	return v.limiter.Allow()
}
