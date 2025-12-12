package wrapify

// NewPagination creates a new instance of the `pagination` struct.
//
// This function initializes a `pagination` struct with its default values.
//
// Returns:
//   - A pointer to a newly created `pagination` instance.
func NewPagination() *pagination {
	p := &pagination{}
	return p
}

// NewMeta creates a new instance of the `meta` struct.
//
// This function initializes a `meta` struct with its default values,
// including an empty `CustomFields` map.
//
// Returns:
//   - A pointer to a newly created `meta` instance with initialized fields.
func NewMeta() *meta {
	m := &meta{
		customFields: map[string]any{},
	}
	return m
}

// NewHeader creates a new instance of the `header` struct.
//
// This function initializes a `header` struct with its default values.
//
// Returns:
//   - A pointer to a newly created `header` instance.
func NewHeader() *header {
	h := &header{}
	return h
}

// New creates a new instance of the `wrapper` struct.
//
// This function initializes a `wrapper` struct with its default values,
// including an empty map for the `Debug` field.
//
// Returns:
//   - A pointer to a newly created `wrapper` instance with initialized fields.
func New() *wrapper {
	w := &wrapper{
		meta: defaultMetaValues(),
	}
	return w
}

// Get returns buffer from pool
//
// If pool is empty, it creates a new buffer
// of predefined size
//
// Returns:
//   - A byte slice buffer
func (bp *BufferPool) Get() []byte {
	select {
	case buf := <-bp.buffers:
		return buf
	default:
		return make([]byte, bp.size)
	}
}

// Put returns buffer to pool
//
// # If pool is full, the buffer is discarded
//
// Parameters:
//   - buf: A byte slice buffer to be returned to the pool
//
// Returns:
//   - None
func (bp *BufferPool) Put(buf []byte) {
	if buf == nil {
		return
	}
	select {
	case bp.buffers <- buf:
	default:
		// Pool is full, discard
	}
}
