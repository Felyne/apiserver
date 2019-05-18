package middleware

import (
	"apiserver/handler"
	"apiserver/pkg/errno"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	defaultBacklogTimeout = time.Second * 30
)

// Throttle is a middleware that limits number of currently processed requests
// at a time.
func Throttle(limit int) gin.HandlerFunc {
	return ThrottleBacklog(limit, 0, defaultBacklogTimeout)
}

// ThrottleBacklog is a middleware that limits number of currently processed
// requests at a time and provides a backlog for holding a finite number of
// pending requests.
func ThrottleBacklog(limit int, backlogLimit int, backlogTimeout time.Duration) gin.HandlerFunc {
	if limit < 1 {
		panic("chi/middleware: Throttle expects limit > 0")
	}

	if backlogLimit < 0 {
		panic("chi/middleware: Throttle expects backlogLimit to be positive")
	}

	t := throttler{
		tokens:         make(chan token, limit),
		backlogTokens:  make(chan token, limit+backlogLimit),
		backlogTimeout: backlogTimeout,
	}

	// Filling tokens.
	for i := 0; i < limit+backlogLimit; i++ {
		if i < limit {
			t.tokens <- token{}
		}
		t.backlogTokens <- token{}
	}

	return gin.HandlerFunc(t.RateLimit)
}

type token struct{}

type throttler struct {
	tokens         chan token
	backlogTokens  chan token
	backlogTimeout time.Duration
}

func (t *throttler) RateLimit(c *gin.Context) {
	ctx := c.Request.Context()
	select {
	case <-ctx.Done():
		handler.SendResponse(c, errno.ErrContextCanceled, nil)
		c.Abort()
	case btok := <-t.backlogTokens:
		timer := time.NewTimer(t.backlogTimeout)
		defer func() {
			t.backlogTokens <- btok
		}()
		select {
		case <-ctx.Done():
			handler.SendResponse(c, errno.ErrContextCanceled, nil)
			c.Abort()
		case <-timer.C:
			handler.SendResponse(c, errno.ErrTimedOut, nil)
			c.Abort()
		case tok := <-t.tokens:
			defer func() {
				t.tokens <- tok
			}()
			c.Next()
		}
	default:
		handler.SendResponse(c, errno.ErrCapacityExceeded, nil)
		c.Abort()
	}
}
