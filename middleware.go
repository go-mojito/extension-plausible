package plausible

// PlausibleMiddleware provides automatic visit tracking for all routes
func PlausibleMiddleware(ctx PlausibleContext, next func() error) error {
	ctx.PageView()
	return next()
}
