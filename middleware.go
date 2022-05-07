package plausible

import "github.com/go-mojito/mojito"

// PlausibleMiddleware provides automatic visit tracking for all routes
func PlausibleMiddleware(ctx Context, logger mojito.Logger, next func() error) error {
	go func(ctx Context) {
		if err := ctx.PageView(); err != nil {
			logger.Error(err)
		}
	}(ctx)
	return next()
}
