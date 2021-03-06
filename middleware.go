package plausible

import "github.com/go-mojito/mojito"

// Middleware provides automatic visit tracking for all routes
func Middleware(ctx Context, logger mojito.Logger, next func() error) error {
	go func(ctx Context) {
		// Skip unwanted URLs
		for _, filter := range filters {
			if filter.MatchString(ctx.Request().GetRequest().URL.Path) {
				return
			}
		}
		if err := ctx.PageView(); err != nil {
			logger.Error(err)
		}
	}(ctx)
	return next()
}
