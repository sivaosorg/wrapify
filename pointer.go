package wrapify

// Ptr returns a pointer to the given value.
// This is useful for creating pointers to literals or values inline,
// especially when working with optional fields in structs or API requests.
//
// Example - API Request with optional fields:
//
//	type UpdateUserRequest struct {
//	    Name  *string `json:"name,omitempty"`
//	    Email *string `json:"email,omitempty"`
//	    Age   *int    `json:"age,omitempty"`
//	}
//
//	// Without Ptr - need temporary variables
//	name := "John Doe"
//	email := "john@example.com"
//	req := UpdateUserRequest{
//	    Name:  &name,
//	    Email: &email,
//	}
//
//	// With Ptr - clean and concise
//	req := UpdateUserRequest{
//	    Name:  wrapify.Ptr("John Doe"),
//	    Email: wrapify.Ptr("john@example.com"),
//	    Age:   wrapify.Ptr(30),
//	}
//
// Example - Database insert with optional fields:
//
//	type User struct {
//	    Name  string
//	    Email *string // Can be NULL
//	    Bio   *string // Can be NULL
//	}
//
//	user := User{
//	    Name:  "Jane Smith",
//	    Email: wrapify.Ptr("jane@example.com"),
//	    Bio:   wrapify.Ptr("Software Engineer"),
//	}
//
// Example - Partial updates:
//
//	// Only update name and email, leave other fields unchanged
//	update := UpdateUserRequest{
//	    Name:  wrapify.Ptr("New Name"),
//	    Email: wrapify.Ptr("new@example.com"),
//	    // Age is nil, won't be updated
//	}
func Ptr[T any](v T) *T {
	return &v
}

// Deref safely dereferences a pointer and returns its value.
// If the pointer is nil, it returns the zero value of type T.
// This is useful for safely accessing pointer values without explicit nil checks.
//
// Example - Safe access to optional config fields:
//
//	type Config struct {
//	    Host    *string
//	    Port    *int
//	    Timeout *int
//	}
//
//	config := Config{
//	    Host: wrapify.Ptr("localhost"),
//	    // Port and Timeout are nil
//	}
//
//	// Without Deref - manual nil checks
//	var port int
//	if config.Port != nil {
//	    port = *config.Port
//	} else {
//	    port = 0 // zero value
//	}
//
//	// With Deref - clean and safe
//	host := wrapify.Deref(config.Host)    // "localhost"
//	port := wrapify.Deref(config.Port)    // 0 (zero value for int)
//	timeout := wrapify.Deref(config.Timeout) // 0
//
// Example - Processing API response with optional fields:
//
//	type APIResponse struct {
//	    Status  string
//	    Message *string
//	    Data    *ResponseData
//	}
//
//	var resp APIResponse
//	json.Unmarshal(data, &resp)
//
//	message := wrapify.Deref(resp.Message) // "" if nil
//	if message != "" {
//	    fmt.Println("Message:", message)
//	}
//
// Example - Database query results:
//
//	type User struct {
//	    ID    int
//	    Name  string
//	    Email *string // Can be NULL in database
//	    Age   *int    // Can be NULL in database
//	}
//
//	var user User
//	db.QueryRow("SELECT id, name, email, age FROM users WHERE id = ?", 1).Scan(
//	    &user.ID, &user.Name, &user.Email, &user.Age,
//	)
//
//	email := wrapify.Deref(user.Email) // "" if NULL
//	age := wrapify.Deref(user.Age)     // 0 if NULL
//
// Note: If you need a custom default value instead of the zero value,
// use DerefOr instead.
func Deref[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

// DerefOr safely dereferences a pointer and returns its value.
// If the pointer is nil, it returns the provided default value instead of the zero value.
// This is useful when you need specific defaults for nil pointers rather than zero values.
//
// Example - Configuration with custom defaults:
//
//	type ServerConfig struct {
//	    Host     *string
//	    Port     *int
//	    Timeout  *time.Duration
//	    MaxConns *int
//	}
//
//	config := ServerConfig{
//	    Host: wrapify.Ptr("localhost"),
//	    // Port, Timeout, MaxConns are nil
//	}
//
//	// Apply defaults for nil values
//	host := wrapify.DerefOr(config.Host, "0.0.0.0")
//	port := wrapify.DerefOr(config.Port, 8080)
//	timeout := wrapify.DerefOr(config.Timeout, 30*time.Second)
//	maxConns := wrapify.DerefOr(config.MaxConns, 100)
//
//	fmt.Printf("Server: %s:%d (timeout: %v, max conns: %d)\n",
//	    host, port, timeout, maxConns)
//	// Output: Server: localhost:8080 (timeout: 30s, max conns: 100)
//
// Example - HTTP client options with defaults:
//
//	type HTTPOptions struct {
//	    Timeout    *time.Duration
//	    MaxRetries *int
//	    UserAgent  *string
//	}
//
//	func NewHTTPClient(opts HTTPOptions) *http.Client {
//	    return &http.Client{
//	        Timeout: wrapify.DerefOr(opts.Timeout, 30*time.Second),
//	        Transport: &http.Transport{
//	            MaxIdleConns: wrapify.DerefOr(opts.MaxRetries, 3),
//	        },
//	    }
//	}
//
//	// Use with partial options
//	client := NewHTTPClient(HTTPOptions{
//	    Timeout: wrapify.Ptr(10 * time.Second),
//	    // MaxRetries will use default: 3
//	})
//
// Example - Database query with default pagination:
//
//	type QueryOptions struct {
//	    Limit  *int
//	    Offset *int
//	    Sort   *string
//	}
//
//	func QueryUsers(opts QueryOptions) []User {
//	    limit := wrapify.DerefOr(opts.Limit, 10)      // Default: 10
//	    offset := wrapify.DerefOr(opts.Offset, 0)     // Default: 0
//	    sort := wrapify.DerefOr(opts.Sort, "id ASC")  // Default: "id ASC"
//
//	    query := fmt.Sprintf("SELECT * FROM users ORDER BY %s LIMIT %d OFFSET %d",
//	        sort, limit, offset)
//	    // Execute query...
//	}
//
//	// Query with custom limit only
//	users := QueryUsers(QueryOptions{
//	    Limit: wrapify.Ptr(20),
//	    // Offset and Sort use defaults
//	})
//
// Example - Feature flags with defaults:
//
//	type FeatureFlags struct {
//	    EnableCache   *bool
//	    EnableMetrics *bool
//	    EnableTracing *bool
//	}
//
//	flags := FeatureFlags{
//	    EnableCache: wrapify.Ptr(true),
//	    // EnableMetrics and EnableTracing are nil
//	}
//
//	cache := wrapify.DerefOr(flags.EnableCache, false)     // true
//	metrics := wrapify.DerefOr(flags.EnableMetrics, true)  // true (default)
//	tracing := wrapify.DerefOr(flags.EnableTracing, false) // false (default)
//
// Example - API rate limiting with defaults:
//
//	type RateLimitConfig struct {
//	    RequestsPerSecond *int
//	    BurstSize         *int
//	    Timeout           *time.Duration
//	}
//
//	func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
//	    return &RateLimiter{
//	        rps:     wrapify.DerefOr(cfg.RequestsPerSecond, 100),
//	        burst:   wrapify.DerefOr(cfg.BurstSize, 10),
//	        timeout: wrapify.DerefOr(cfg.Timeout, 5*time.Second),
//	    }
//	}
//
// Note: The difference between Deref and DerefOr:
//   - Deref(ptr) returns zero value if nil: Deref[int](nil) = 0
//   - DerefOr(ptr, 100) returns custom default if nil: DerefOr[int](nil, 100) = 100
func DerefOr[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
