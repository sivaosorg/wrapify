package wrapify

// WrapMiddleware defines an interface for middleware that processes
// the response wrapper. It allows for custom processing of the response
// before it is returned to the client. This can include adding headers,
// modifying the response body, or logging additional information.
// The middleware should implement the Process method, which takes a
// wrapper and returns a modified wrapper. This allows for chaining
// multiple middleware components together, enabling a flexible and
// extensible architecture for handling API responses.
type WrapMiddleware interface {
	Process(w *wrapper) *wrapper
}

// SecurityMiddleware is a middleware that adds security-related headers
// to the response. It implements the ResponseMiddleware interface
// and provides a method to process the response wrapper by adding
// security headers such as X-Content-Type-Options, X-Frame-Options,
// and X-XSS-Protection. This middleware is useful for enhancing the
// security of the API responses by preventing common web vulnerabilities
// such as MIME type sniffing, click-jacking, and cross-site scripting (XSS).
type SecurityMiddleware struct{}

func (s *SecurityMiddleware) Process(w *wrapper) *wrapper {
	return w.WithDebuggingKV("security_headers", map[string]string{
		HeaderXContentTypeOpts: "nosniff",
		HeaderXFrameOptions:    "DENY",
		HeaderXXSSProtection:   "1; mode=block",
	})
}
