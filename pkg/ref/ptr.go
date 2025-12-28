package ref

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

// OrElse returns the pointer if it's not nil, otherwise returns the alternative pointer.
// This is useful for providing fallback pointer values.
//
// Example - Configuration fallback:
//
//	type Config struct {
//	    Timeout *time.Duration
//	}
//
//	userConfig := Config{} // Timeout is nil
//	defaultConfig := Config{
//	    Timeout: wrapify.Ptr(30 * time.Second),
//	}
//
//	timeout := wrapify.OrElse(userConfig.Timeout, defaultConfig.Timeout)
//	fmt.Println(*timeout) // 30s
//
// Example - HTTP client options:
//
//	type HTTPOptions struct {
//	    UserAgent *string
//	    Timeout   *time.Duration
//	}
//
//	opts := HTTPOptions{
//	    UserAgent: wrapify.Ptr("custom-agent"),
//	}
//
//	defaultOpts := HTTPOptions{
//	    UserAgent: wrapify.Ptr("default-agent"),
//	    Timeout:   wrapify.Ptr(10 * time.Second),
//	}
//
//	ua := wrapify.OrElse(opts.UserAgent, defaultOpts.UserAgent)
//	to := wrapify.OrElse(opts.Timeout, defaultOpts.Timeout)
//
// Example - Response data with cache:
//
//	liveData := fetchFromAPI()  // May return nil
//	cachedData := getFromCache() // Fallback
//
//	data := wrapify.OrElse(liveData, cachedData)
func OrElse[T any](ptr, alternative *T) *T {
	if ptr != nil {
		return ptr
	}
	return alternative
}

// OrElseGet returns the pointer if it's not nil, otherwise calls the supplier function.
// This is useful when the fallback value is expensive to compute and should only be
// computed if needed (lazy evaluation).
//
// Example - Lazy database query:
//
//	func GetUserName(userID int) *string {
//	    // Try cache first
//	    cached := getFromCache(userID)
//
//	    // Only query database if cache miss
//	    return wrapify.OrElseGet(cached, func() *string {
//	        return queryDatabase(userID) // Expensive operation
//	    })
//	}
//
// Example - Lazy configuration loading:
//
//	type Config struct {
//	    APIKey *string
//	}
//
//	envConfig := loadFromEnv() // Fast
//
//	config := wrapify.OrElseGet(envConfig.APIKey, func() *string {
//	    // Only load from file if env var not set
//	    cfg := loadFromFile() // Slow
//	    return cfg.APIKey
//	})
//
// Example - Lazy default value generation:
//
//	sessionID := wrapify.OrElseGet(existingSession, func() *string {
//	    // Only generate new UUID if no existing session
//	    return wrapify.Ptr(uuid.New().String())
//	})
//
// Example - Conditional API call:
//
//	cachedPrice := getCachedPrice(productID)
//
//	price := wrapify.OrElseGet(cachedPrice, func() *float64 {
//	    // Only call pricing API if cache miss
//	    return fetchPriceFromAPI(productID)
//	})
func OrElseGet[T any](ptr *T, supplier func() *T) *T {
	if ptr != nil {
		return ptr
	}
	return supplier()
}

// If executes the consumer function if the pointer is not nil.
// Returns true if the consumer was executed, false otherwise.
// This is useful for side effects on non-nil values.
//
// Example - Logging optional fields:
//
//	type User struct {
//	    Name  string
//	    Email *string
//	    Phone *string
//	}
//
//	user := User{
//	    Name:  "John",
//	    Email: wrapify.Ptr("john@example.com"),
//	}
//
//	wrapify.If(user.Email, func(email string) {
//	    log.Printf("Email: %s", email)
//	})
//
//	wrapify.If(user.Phone, func(phone string) {
//	    log.Printf("Phone: %s", phone) // Not executed
//	})
//
// Example - Sending notifications:
//
//	type Notification struct {
//	    Email *string
//	    SMS   *string
//	    Push  *string
//	}
//
//	notif := Notification{
//	    Email: wrapify.Ptr("user@example.com"),
//	    SMS:   wrapify.Ptr("+1234567890"),
//	}
//
//	wrapify.If(notif.Email, func(email string) {
//	    sendEmail(email, "Hello!")
//	})
//
//	wrapify.If(notif.SMS, func(phone string) {
//	    sendSMS(phone, "Hello!")
//	})
//
//	wrapify.If(notif.Push, func(token string) {
//	    sendPush(token, "Hello!") // Not executed
//	})
//
// Example - Processing optional metadata:
//
//	type Event struct {
//	    Name     string
//	    Metadata *map[string]string
//	}
//
//	event := Event{
//	    Name:     "user.login",
//	    Metadata: wrapify.Ptr(map[string]string{"ip": "127.0.0.1"}),
//	}
//
//	wrapify.If(event.Metadata, func(meta map[string]string) {
//	    for k, v := range meta {
//	        log.Printf("  %s: %s", k, v)
//	    }
//	})
func If[T any](ptr *T, consumer func(T)) bool {
	if ptr != nil {
		consumer(*ptr)
		return true
	}
	return false
}

// IfElse executes onPresent if pointer is not nil, otherwise executes onAbsent.
// This is useful for branching logic based on pointer presence.
//
// Example - User greeting with fallback:
//
//	user := User{
//	    Name:          "John",
//	    PreferredName: wrapify.Ptr("Johnny"),
//	}
//
//	wrapify.IfElse(user.PreferredName,
//	    func(name string) {
//	        fmt.Printf("Hello, %s!", name)
//	    },
//	    func() {
//	        fmt.Printf("Hello, %s!", user.Name)
//	    },
//	)
//
// Example - Configuration with default behavior:
//
//	config := Config{
//	    CustomHandler: wrapify.Ptr(myHandler),
//	}
//
//	wrapify.IfElse(config.CustomHandler,
//	    func(handler Handler) {
//	        handler.Handle(request)
//	    },
//	    func() {
//	        defaultHandler.Handle(request)
//	    },
//	)
//
// Example - Caching strategy:
//
//	cached := getFromCache(key)
//
//	wrapify.IfElse(cached,
//	    func(data Data) {
//	        log.Println("Cache hit")
//	        processData(data)
//	    },
//	    func() {
//	        log.Println("Cache miss, fetching...")
//	        data := fetchFromDB(key)
//	        saveToCache(key, data)
//	        processData(data)
//	    },
//	)
func IfElse[T any](ptr *T, onPresent func(T), onAbsent func()) {
	if ptr != nil {
		onPresent(*ptr)
	} else {
		onAbsent()
	}
}

// Must dereferences the pointer and panics if it's nil.
// This should only be used when you're certain the pointer is not nil,
// such as after validation or in test code.
//
// Example - After validation:
//
//	func ProcessUser(req CreateUserRequest) error {
//	    if wrapify.IsNil(req.Name) {
//	        return errors.New("name is required")
//	    }
//	    if wrapify.IsNil(req.Email) {
//	        return errors.New("email is required")
//	    }
//
//	    // Safe to use Must after validation
//	    name := wrapify.Must(req.Name)
//	    email := wrapify.Must(req.Email)
//
//	    createUser(name, email)
//	    return nil
//	}
//
// Example - Test code:
//
//	func TestUserCreation(t *testing.T) {
//	    user := createTestUser()
//
//	    // In tests, we expect these to be non-nil
//	    name := wrapify.Must(user.Name)
//	    email := wrapify.Must(user.Email)
//
//	    assert.Equal(t, "Test User", name)
//	    assert.Equal(t, "test@example.com", email)
//	}
//
// Example - Configuration that must exist:
//
//	config := loadConfig()
//
//	// These are required, panic if missing
//	dbHost := wrapify.Must(config.Database.Host)
//	dbPort := wrapify.Must(config.Database.Port)
//	apiKey := wrapify.Must(config.API.Key)
//
//	connectDB(dbHost, dbPort, apiKey)
//
// Warning: Only use Must when you're absolutely certain the pointer is not nil.
// For production code, prefer using If, IfElse, or explicit nil checks.
func Must[T any](ptr *T) T {
	if ptr == nil {
		panic("wrapify: attempted to dereference nil pointer")
	}
	return *ptr
}

// MustPtr is like Must but returns a panic-safe pointer dereference.
// Panics with a custom message if the pointer is nil.
//
// Example - With custom error message:
//
//	func GetConfig() Config {
//	    cfg := loadConfigFromFile()
//
//	    return Config{
//	        DBHost: wrapify.MustPtr(cfg.DBHost, "database host is required"),
//	        DBPort: wrapify.MustPtr(cfg.DBPort, "database port is required"),
//	        APIKey: wrapify.MustPtr(cfg.APIKey, "API key is required"),
//	    }
//	}
//
// Example - Critical validation:
//
//	func ProcessPayment(req PaymentRequest) error {
//	    amount := wrapify.MustPtr(req.Amount, "payment amount is required")
//	    currency := wrapify.MustPtr(req.Currency, "currency is required")
//	    cardToken := wrapify.MustPtr(req.CardToken, "card token is required")
//
//	    return processPayment(amount, currency, cardToken)
//	}
func MustPtr[T any](ptr *T, message string) T {
	if ptr == nil {
		panic("wrapify: " + message)
	}
	return *ptr
}

// FlatMap applies a function that returns a pointer to the dereferenced value.
// If the input pointer is nil, returns nil without calling the function.
// This is useful for chaining operations that may return nil.
//
// Example - Nested optional access:
//
//	type Address struct {
//	    City    string
//	    ZipCode *string
//	}
//
//	type User struct {
//	    Name    string
//	    Address *Address
//	}
//
//	user := User{
//	    Name:    "John",
//	    Address: wrapify.Ptr(Address{City: "NYC"}),
//	}
//
//	zipCode := wrapify.FlatMap(user.Address, func(addr Address) *string {
//	    return addr.ZipCode
//	})
//	// Returns nil if Address is nil OR ZipCode is nil
//
// Example - Database lookup chain:
//
//	userID := wrapify.Ptr(123)
//
//	user := wrapify.FlatMap(userID, func(id int) *User {
//	    return findUserByID(id) // May return nil
//	})
//
//	address := wrapify.FlatMap(user, func(u User) *Address {
//	    return u.Address // May be nil
//	})
//
// Example - API response chain:
//
//	type APIResponse struct {
//	    Data *ResponseData
//	}
//
//	type ResponseData struct {
//	    User *User
//	}
//
//	resp := callAPI()
//
//	user := wrapify.FlatMap(resp.Data, func(data ResponseData) *User {
//	    return data.User
//	})
//
// Example - Safe navigation:
//
//	type Config struct {
//	    Database *DatabaseConfig
//	}
//
//	type DatabaseConfig struct {
//	    Connection *ConnectionConfig
//	}
//
//	config := loadConfig()
//
//	connString := wrapify.FlatMap(config.Database, func(db DatabaseConfig) *string {
//	    return wrapify.FlatMap(db.Connection, func(conn ConnectionConfig) *string {
//	        return wrapify.Ptr(conn.String())
//	    })
//	})
func FlatMap[T any, R any](ptr *T, fn func(T) *R) *R {
	if ptr == nil {
		return nil
	}
	return fn(*ptr)
}

// Validate returns the pointer if it passes validation, otherwise returns nil.
// The validator function should return an error if validation fails.
// This is useful for filtering values based on complex validation rules.
//
// Example - Email validation:
//
//	func isValidEmail(email string) error {
//	    if !strings.Contains(email, "@") {
//	        return errors.New("invalid email format")
//	    }
//	    return nil
//	}
//
//	email := wrapify.Ptr("user@example.com")
//	validEmail := wrapify.Validate(email, isValidEmail)
//	// Returns email pointer
//
//	badEmail := wrapify.Ptr("invalid-email")
//	invalidEmail := wrapify.Validate(badEmail, isValidEmail)
//	// Returns nil
//
// Example - Age validation:
//
//	func isAdult(age int) error {
//	    if age < 18 {
//	        return errors.New("must be 18 or older")
//	    }
//	    return nil
//	}
//
//	age := wrapify.Ptr(25)
//	validAge := wrapify.Validate(age, isAdult) // Returns age pointer
//
//	minorAge := wrapify.Ptr(15)
//	invalidAge := wrapify.Validate(minorAge, isAdult) // Returns nil
//
// Example - Password strength validation:
//
//	func isStrongPassword(pwd string) error {
//	    if len(pwd) < 8 {
//	        return errors.New("password too short")
//	    }
//	    return nil
//	}
//
//	password := wrapify.Ptr("StrongP@ss123")
//	validPassword := wrapify.Validate(password, isStrongPassword)
//
//	if wrapify.IsNotNil(validPassword) {
//	    // Password is valid, proceed
//	    hashAndStore(*validPassword)
//	}
//
// Example - URL validation:
//
//	func isValidURL(url string) error {
//	    _, err := url.Parse(url)
//	    return err
//	}
//
//	inputURL := wrapify.Ptr("https://example.com")
//	validURL := wrapify.Validate(inputURL, isValidURL)
func Validate[T any](ptr *T, validator func(T) error) *T {
	if ptr == nil {
		return nil
	}
	if err := validator(*ptr); err != nil {
		return nil
	}
	return ptr
}

// MapOr applies a function to the pointer value if not nil, otherwise returns default.
// This combines Map and DerefOr in a single operation.
//
// Example - String formatting with default:
//
//	name := wrapify.Ptr("john")
//
//	formatted := wrapify.MapOr(name,
//	    func(n string) string {
//	        return strings.ToUpper(n)
//	    },
//	    "ANONYMOUS",
//	)
//	// Returns "JOHN"
//
//	var nilName *string
//	formatted2 := wrapify.MapOr(nilName,
//	    func(n string) string {
//	        return strings.ToUpper(n)
//	    },
//	    "ANONYMOUS",
//	)
//	// Returns "ANONYMOUS"
//
// Example - Price calculation with default:
//
//	discount := wrapify.Ptr(0.1) // 10% discount
//
//	finalPrice := wrapify.MapOr(discount,
//	    func(d float64) float64 {
//	        return 100.0 * (1 - d)
//	    },
//	    100.0, // No discount
//	)
//	// Returns 90.0
//
// Example - Date formatting:
//
//	createdAt := wrapify.Ptr(time.Now())
//
//	formatted := wrapify.MapOr(createdAt,
//	    func(t time.Time) string {
//	        return t.Format("2006-01-02")
//	    },
//	    "N/A",
//	)
func MapOr[T any, R any](ptr *T, fn func(T) R, defaultValue R) R {
	if ptr == nil {
		return defaultValue
	}
	return fn(*ptr)
}

// All checks if all pointers in the slice are non-nil.
// Returns true if all pointers are non-nil, false otherwise.
// This is useful for validating that all required fields are present.
//
// Example - Validating required fields:
//
//	type CreateUserRequest struct {
//	    Name     *string
//	    Email    *string
//	    Password *string
//	}
//
//	req := CreateUserRequest{
//	    Name:     wrapify.Ptr("John"),
//	    Email:    wrapify.Ptr("john@example.com"),
//	    Password: wrapify.Ptr("secret"),
//	}
//
//	if wrapify.All(req.Name, req.Email, req.Password) {
//	    createUser(*req.Name, *req.Email, *req.Password)
//	} else {
//	    return errors.New("all fields are required")
//	}
//
// Example - Configuration validation:
//
//	type DatabaseConfig struct {
//	    Host     *string
//	    Port     *int
//	    Database *string
//	    User     *string
//	    Password *string
//	}
//
//	func ValidateConfig(cfg DatabaseConfig) error {
//	    if !wrapify.All(cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password) {
//	        return errors.New("all database config fields are required")
//	    }
//	    return nil
//	}
//
// Example - Payment validation:
//
//	type Payment struct {
//	    Amount   *float64
//	    Currency *string
//	    CardToken *string
//	}
//
//	payment := Payment{
//	    Amount:   wrapify.Ptr(100.0),
//	    Currency: wrapify.Ptr("USD"),
//	    // CardToken is nil
//	}
//
//	if !wrapify.All(payment.Amount, payment.Currency, payment.CardToken) {
//	    return errors.New("incomplete payment information")
//	}
func All[T any](ptrs ...*T) bool {
	for _, ptr := range ptrs {
		if ptr == nil {
			return false
		}
	}
	return true
}

// Any checks if any pointer in the slice is non-nil.
// Returns true if at least one pointer is non-nil, false otherwise.
// This is useful for checking if any optional field is provided.
//
// Example - Checking if any contact method provided:
//
//	type ContactInfo struct {
//	    Email *string
//	    Phone *string
//	    SMS   *string
//	}
//
//	contact := ContactInfo{
//	    Phone: wrapify.Ptr("+1234567890"),
//	}
//
//	if wrapify.Any(contact.Email, contact.Phone, contact.SMS) {
//	    fmt.Println("At least one contact method provided")
//	} else {
//	    return errors.New("at least one contact method is required")
//	}
//
// Example - Search filter validation:
//
//	type SearchFilter struct {
//	    Query     *string
//	    Category  *string
//	    MinPrice  *float64
//	    MaxPrice  *float64
//	}
//
//	filter := SearchFilter{} // All nil
//
//	if !wrapify.Any(filter.Query, filter.Category, filter.MinPrice, filter.MaxPrice) {
//	    return errors.New("at least one search criteria is required")
//	}
//
// Example - Update validation:
//
//	type UpdateUserRequest struct {
//	    Name  *string
//	    Email *string
//	    Bio   *string
//	}
//
//	req := UpdateUserRequest{} // All nil
//
//	if !wrapify.Any(req.Name, req.Email, req.Bio) {
//	    return errors.New("at least one field must be updated")
//	}
func Any[T any](ptrs ...*T) bool {
	for _, ptr := range ptrs {
		if ptr != nil {
			return true
		}
	}
	return false
}
