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

// IsNil checks if a pointer is nil.
// This is useful for explicit nil checks in conditional logic or validation.
//
// Example - Validating required fields:
//
//	type CreateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func ValidateCreateUser(req CreateUserRequest) error {
//	    if wrapify.IsNil(req.Name) {
//	        return errors.New("name is required")
//	    }
//	    if wrapify.IsNil(req.Email) {
//	        return errors.New("email is required")
//	    }
//	    // Age is optional, no need to check
//	    return nil
//	}
//
// Example - Conditional processing:
//
//	type Config struct {
//	    Database *DatabaseConfig
//	    Cache    *CacheConfig
//	    Queue    *QueueConfig
//	}
//
//	func InitializeServices(cfg Config) {
//	    if !wrapify.IsNil(cfg.Database) {
//	        initDatabase(cfg.Database)
//	    }
//	    if !wrapify.IsNil(cfg.Cache) {
//	        initCache(cfg.Cache)
//	    }
//	    if !wrapify.IsNil(cfg.Queue) {
//	        initQueue(cfg.Queue)
//	    }
//	}
//
// Example - Checking optional API response fields:
//
//	type APIResponse struct {
//	    Data  *ResponseData
//	    Error *APIError
//	}
//
//	resp := callAPI()
//	if wrapify.IsNil(resp.Error) {
//	    // Success case
//	    processData(resp.Data)
//	} else {
//	    // Error case
//	    handleError(resp.Error)
//	}
func IsNil[T any](ptr *T) bool {
	return ptr == nil
}

// IsNotNil checks if a pointer is not nil.
// This is the inverse of IsNil and useful for positive condition checks.
//
// Example - Processing optional fields when present:
//
//	type User struct {
//	    ID      int
//	    Name    string
//	    Email   *string
//	    Phone   *string
//	    Address *Address
//	}
//
//	func DisplayUserInfo(user User) {
//	    fmt.Printf("User: %s\n", user.Name)
//
//	    if wrapify.IsNotNil(user.Email) {
//	        fmt.Printf("Email: %s\n", *user.Email)
//	    }
//	    if wrapify.IsNotNil(user.Phone) {
//	        fmt.Printf("Phone: %s\n", *user.Phone)
//	    }
//	    if wrapify.IsNotNil(user.Address) {
//	        fmt.Printf("Address: %v\n", *user.Address)
//	    }
//	}
//
// Example - Counting non-nil fields:
//
//	func CountProvidedFields(req UpdateUserRequest) int {
//	    count := 0
//	    if wrapify.IsNotNil(req.Name) {
//	        count++
//	    }
//	    if wrapify.IsNotNil(req.Email) {
//	        count++
//	    }
//	    if wrapify.IsNotNil(req.Age) {
//	        count++
//	    }
//	    return count
//	}
func IsNotNil[T any](ptr *T) bool {
	return ptr != nil
}

// Equal checks if two pointers point to equal values.
// Returns true if both are nil, or both are non-nil and their values are equal.
// Returns false if one is nil and the other is not.
//
// Example - Comparing optional configuration:
//
//	oldConfig := Config{
//	    Timeout: wrapify.Ptr(30 * time.Second),
//	    MaxRetries: wrapify.Ptr(3),
//	}
//
//	newConfig := Config{
//	    Timeout: wrapify.Ptr(30 * time.Second),
//	    MaxRetries: wrapify.Ptr(5),
//	}
//
//	if wrapify.Equal(oldConfig.Timeout, newConfig.Timeout) {
//	    fmt.Println("Timeout unchanged")
//	}
//	if !wrapify.Equal(oldConfig.MaxRetries, newConfig.MaxRetries) {
//	    fmt.Println("MaxRetries changed")
//	}
//
// Example - Detecting changes in optional fields:
//
//	type UpdateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func HasChanges(current User, update UpdateUserRequest) bool {
//	    if !wrapify.Equal(wrapify.Ptr(current.Name), update.Name) {
//	        return true
//	    }
//	    if !wrapify.Equal(current.Email, update.Email) {
//	        return true
//	    }
//	    if !wrapify.Equal(current.Age, update.Age) {
//	        return true
//	    }
//	    return false
//	}
//
// Example - Comparing API responses:
//
//	oldResp := &APIResponse{Status: "success"}
//	newResp := &APIResponse{Status: "success"}
//
//	if wrapify.Equal(wrapify.Ptr(oldResp.Status), wrapify.Ptr(newResp.Status)) {
//	    fmt.Println("Status is the same")
//	}
//
// Note: This function requires type T to be comparable.
// For non-comparable types (slices, maps, functions), use custom comparison logic.
func Equal[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Copy creates a new pointer with a copy of the value.
// If the input pointer is nil, returns nil.
// This is useful for creating independent copies of pointer values.
//
// Example - Deep copying configuration:
//
//	originalConfig := Config{
//	    Host: wrapify.Ptr("localhost"),
//	    Port: wrapify.Ptr(8080),
//	}
//
//	// Create independent copy
//	configCopy := Config{
//	    Host: wrapify.Copy(originalConfig.Host),
//	    Port: wrapify.Copy(originalConfig.Port),
//	}
//
//	// Modifying copy doesn't affect original
//	*configCopy.Host = "example.com"
//	fmt.Println(*originalConfig.Host) // Still "localhost"
//
// Example - Cloning user data:
//
//	type User struct {
//	    Name  *string
//	    Email *string
//	    Age   *int
//	}
//
//	func CloneUser(user User) User {
//	    return User{
//	        Name:  wrapify.Copy(user.Name),
//	        Email: wrapify.Copy(user.Email),
//	        Age:   wrapify.Copy(user.Age),
//	    }
//	}
//
//	original := User{
//	    Name:  wrapify.Ptr("John"),
//	    Email: wrapify.Ptr("john@example.com"),
//	}
//
//	clone := CloneUser(original)
//	*clone.Name = "Jane" // Doesn't affect original
//
// Example - Creating independent request copies:
//
//	originalReq := UpdateUserRequest{
//	    Name: wrapify.Ptr("Original Name"),
//	    Age:  wrapify.Ptr(25),
//	}
//
//	// Create copy for retry with modifications
//	retryReq := UpdateUserRequest{
//	    Name: wrapify.Copy(originalReq.Name),
//	    Age:  wrapify.Ptr(26), // Different value
//	}
func Copy[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	v := *ptr
	return &v
}

// CoalescePtr returns the first non-nil pointer from the arguments.
// If all pointers are nil, returns nil.
// This is useful for fallback chains and default value selection.
//
// Example - Configuration fallback chain:
//
//	type Config struct {
//	    Host *string
//	    Port *int
//	}
//
//	userConfig := Config{Port: wrapify.Ptr(3000)}
//	envConfig := Config{Host: wrapify.Ptr("env-host")}
//	defaultConfig := Config{
//	    Host: wrapify.Ptr("localhost"),
//	    Port: wrapify.Ptr(8080),
//	}
//
//	// Use first non-nil value: user > env > default
//	finalHost := wrapify.CoalescePtr(
//	    userConfig.Host,
//	    envConfig.Host,
//	    defaultConfig.Host,
//	) // Returns envConfig.Host
//
//	finalPort := wrapify.CoalescePtr(
//	    userConfig.Port,
//	    envConfig.Port,
//	    defaultConfig.Port,
//	) // Returns userConfig.Port
//
// Example - API response with multiple possible sources:
//
//	type Response struct {
//	    PrimaryData   *Data
//	    SecondaryData *Data
//	    CachedData    *Data
//	}
//
//	resp := getResponse()
//	data := wrapify.CoalescePtr(
//	    resp.PrimaryData,
//	    resp.SecondaryData,
//	    resp.CachedData,
//	) // Use first available data source
//
// Example - User preference with fallbacks:
//
//	type UserSettings struct {
//	    Theme    *string
//	    Language *string
//	}
//
//	userPref := UserSettings{Language: wrapify.Ptr("en")}
//	orgPref := UserSettings{Theme: wrapify.Ptr("dark")}
//	systemPref := UserSettings{
//	    Theme:    wrapify.Ptr("light"),
//	    Language: wrapify.Ptr("en-US"),
//	}
//
//	theme := wrapify.CoalescePtr(
//	    userPref.Theme,
//	    orgPref.Theme,
//	    systemPref.Theme,
//	) // "dark"
//
//	language := wrapify.CoalescePtr(
//	    userPref.Language,
//	    orgPref.Language,
//	    systemPref.Language,
//	) // "en"
func CoalescePtr[T any](ptrs ...*T) *T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return ptr
		}
	}
	return nil
}

// Coalesce returns the first non-nil pointer's value from the arguments.
// If all pointers are nil, returns the zero value of type T.
// This is similar to CoalescePtr but returns the dereferenced value.
//
// Example - Configuration with multiple fallbacks:
//
//	userTimeout := wrapify.Ptr(5 * time.Second)
//	var envTimeout *time.Duration // nil
//	defaultTimeout := wrapify.Ptr(30 * time.Second)
//
//	timeout := wrapify.Coalesce(
//	    userTimeout,
//	    envTimeout,
//	    defaultTimeout,
//	) // Returns 5 * time.Second
//
// Example - String with fallbacks:
//
//	var userTitle *string // nil
//	var orgTitle *string  // nil
//	defaultTitle := wrapify.Ptr("Untitled")
//
//	title := wrapify.Coalesce(
//	    userTitle,
//	    orgTitle,
//	    defaultTitle,
//	) // Returns "Untitled"
//
// Example - Database field with fallbacks:
//
//	type User struct {
//	    PreferredName *string
//	    DisplayName   *string
//	    Username      string
//	}
//
//	user := User{
//	    DisplayName: wrapify.Ptr("John Doe"),
//	    Username:    "john123",
//	}
//
//	name := wrapify.Coalesce(
//	    user.PreferredName,
//	    user.DisplayName,
//	    wrapify.Ptr(user.Username),
//	) // Returns "John Doe"
//
// Example - Numeric with zero value fallback:
//
//	var customLimit *int     // nil
//	var userLimit *int       // nil
//	var defaultLimit *int    // nil
//
//	limit := wrapify.Coalesce(
//	    customLimit,
//	    userLimit,
//	    defaultLimit,
//	) // Returns 0 (zero value for int)
func Coalesce[T any](ptrs ...*T) T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return *ptr
		}
	}
	var zero T
	return zero
}

// ToPtr converts a value to a pointer, but only if the value is not the zero value.
// If the value is the zero value, returns nil.
// This is useful for omitting zero values in API requests or database updates.
//
// Example - API request omitting zero values:
//
//	type UpdateUserRequest struct {
//	    Name  *string `json:"name,omitempty"`
//	    Age   *int    `json:"age,omitempty"`
//	    Score *int    `json:"score,omitempty"`
//	}
//
//	name := "John"
//	age := 0      // Zero value, should be omitted
//	score := 100
//
//	req := UpdateUserRequest{
//	    Name:  wrapify.ToPtr(name),   // &"John"
//	    Age:   wrapify.ToPtr(age),    // nil (zero value)
//	    Score: wrapify.ToPtr(score),  // &100
//	}
//
//	// JSON output: {"name":"John","score":100}
//	// Age is omitted because it's nil
//
// Example - Database partial update:
//
//	func UpdateUser(id int, name string, age int, email string) {
//	    update := map[string]interface{}{}
//
//	    if name := wrapify.ToPtr(name); name != nil {
//	        update["name"] = *name
//	    }
//	    if age := wrapify.ToPtr(age); age != nil {
//	        update["age"] = *age
//	    }
//	    if email := wrapify.ToPtr(email); email != nil {
//	        update["email"] = *email
//	    }
//
//	    // Only update non-zero fields
//	    db.Model(&User{}).Where("id = ?", id).Updates(update)
//	}
//
// Example - Form data validation:
//
//	type SearchFilter struct {
//	    Query     *string
//	    MinPrice  *float64
//	    MaxPrice  *float64
//	    InStock   *bool
//	}
//
//	query := "laptop"
//	minPrice := 0.0    // Zero, will be nil
//	maxPrice := 1000.0
//	inStock := false   // Zero value for bool, will be nil
//
//	filter := SearchFilter{
//	    Query:    wrapify.ToPtr(query),    // &"laptop"
//	    MinPrice: wrapify.ToPtr(minPrice), // nil
//	    MaxPrice: wrapify.ToPtr(maxPrice), // &1000.0
//	    InStock:  wrapify.ToPtr(inStock),  // nil
//	}
//
// Note: This function requires type T to be comparable to detect zero values.
func ToPtr[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// Map applies a function to the value pointed to by the pointer.
// If the pointer is nil, returns nil.
// This is useful for transforming pointer values without explicit nil checks.
//
// Example - String transformations:
//
//	name := wrapify.Ptr("john doe")
//
//	// Convert to uppercase
//	upperName := wrapify.Map(name, strings.ToUpper)
//	fmt.Println(*upperName) // "JOHN DOE"
//
//	// Nil input returns nil
//	var nilName *string
//	result := wrapify.Map(nilName, strings.ToUpper)
//	fmt.Println(result) // nil
//
// Example - Numeric transformations:
//
//	price := wrapify.Ptr(100.0)
//
//	// Apply tax
//	priceWithTax := wrapify.Map(price, func(p float64) float64 {
//	    return p * 1.1 // 10% tax
//	})
//	fmt.Println(*priceWithTax) // 110.0
//
// Example - Type conversions:
//
//	intValue := wrapify.Ptr(42)
//
//	// Convert int to string
//	stringValue := wrapify.Map(intValue, func(i int) string {
//	    return fmt.Sprintf("value: %d", i)
//	})
//	fmt.Println(*stringValue) // "value: 42"
//
// Example - Processing optional user input:
//
//	type UserInput struct {
//	    Email *string
//	}
//
//	input := UserInput{
//	    Email: wrapify.Ptr("  USER@EXAMPLE.COM  "),
//	}
//
//	// Normalize email
//	normalizedEmail := wrapify.Map(input.Email, func(email string) string {
//	    return strings.ToLower(strings.TrimSpace(email))
//	})
//	fmt.Println(*normalizedEmail) // "user@example.com"
//
// Example - Chain transformations:
//
//	age := wrapify.Ptr(25)
//	ageGroup := wrapify.Map(age, func(a int) string {
//	    if a < 18 {
//	        return "minor"
//	    } else if a < 65 {
//	        return "adult"
//	    }
//	    return "senior"
//	})
//	fmt.Println(*ageGroup) // "adult"
func Map[T any, R any](ptr *T, fn func(T) R) *R {
	if ptr == nil {
		return nil
	}
	result := fn(*ptr)
	return &result
}

// Filter returns the pointer if the predicate function returns true, otherwise returns nil.
// If the input pointer is nil, returns nil.
// This is useful for conditional filtering of pointer values.
//
// Example - Validating optional fields:
//
//	age := wrapify.Ptr(15)
//
//	// Only keep age if >= 18
//	validAge := wrapify.Filter(age, func(a int) bool {
//	    return a >= 18
//	})
//	fmt.Println(validAge) // nil (15 < 18)
//
//	adult := wrapify.Ptr(25)
//	validAdult := wrapify.Filter(adult, func(a int) bool {
//	    return a >= 18
//	})
//	fmt.Println(*validAdult) // 25
//
// Example - Email validation:
//
//	email := wrapify.Ptr("invalid-email")
//
//	validEmail := wrapify.Filter(email, func(e string) bool {
//	    return strings.Contains(e, "@")
//	})
//	fmt.Println(validEmail) // nil
//
//	goodEmail := wrapify.Ptr("user@example.com")
//	validGoodEmail := wrapify.Filter(goodEmail, func(e string) bool {
//	    return strings.Contains(e, "@")
//	})
//	fmt.Println(*validGoodEmail) // "user@example.com"
//
// Example - Price range filtering:
//
//	type Product struct {
//	    Name  string
//	    Price *float64
//	}
//
//	products := []Product{
//	    {Name: "Cheap", Price: wrapify.Ptr(10.0)},
//	    {Name: "Expensive", Price: wrapify.Ptr(1000.0)},
//	    {Name: "Affordable", Price: wrapify.Ptr(50.0)},
//	}
//
//	for _, p := range products {
//	    // Only products under $100
//	    affordablePrice := wrapify.Filter(p.Price, func(price float64) bool {
//	        return price < 100
//	    })
//	    if wrapify.IsNotNil(affordablePrice) {
//	        fmt.Printf("%s: $%.2f\n", p.Name, *affordablePrice)
//	    }
//	}
//	// Output:
//	// Cheap: $10.00
//	// Affordable: $50.00
//
// Example - Filtering with multiple conditions:
//
//	username := wrapify.Ptr("ab")
//
//	validUsername := wrapify.Filter(username, func(u string) bool {
//	    return len(u) >= 3 && len(u) <= 20
//	})
//	fmt.Println(validUsername) // nil (too short)
func Filter[T any](ptr *T, predicate func(T) bool) *T {
	if ptr == nil {
		return nil
	}
	if predicate(*ptr) {
		return ptr
	}
	return nil
}
