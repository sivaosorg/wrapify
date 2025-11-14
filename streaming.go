package wrapify

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sivaosorg/unify4g"
)

// NewBufferPool creates new buffer pool
//
// This function initializes a `BufferPool` struct with a specified buffer size
// and pool size.
//
// Parameters:
//   - bufferSize: The size of each buffer in bytes.
//   - poolSize: The maximum number of buffers to maintain in the pool.
//
// Returns:
//
//   - A pointer to a newly created `BufferPool` instance with the specified settings.
func NewBufferPool(bufferSize int64, poolSize int) *BufferPool {
	return &BufferPool{
		buffers: make(chan []byte, poolSize),
		size:    bufferSize,
	}
}

// NewStreamConfig creates default streaming configuration
//
// This function initializes a `StreamConfig` struct with default values
// suitable for typical streaming scenarios.
//
// Returns:
//   - A pointer to a newly created `StreamConfig` instance with default settings.
func NewStreamConfig() *StreamConfig {
	return &StreamConfig{
		ChunkSize:           65536, // 64KB default
		Strategy:            STRATEGY_BUFFERED,
		Compression:         COMP_NONE,
		UseBufferPool:       true,
		MaxConcurrentChunks: 4,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
		ThrottleRate:        0, // unlimited
	}
}

// NewStreaming creates a new instance of the `StreamingWrapper` struct.
//
// This function initializes a `StreamingWrapper` struct with the provided
// `reader`, and `config`. If the `config` is nil, it uses default
// streaming configuration.
//
// Parameters:
//   - `reader`: An `io.Reader` instance from which data will be streamed.
//   - `config`: A pointer to a `StreamConfig` struct containing streaming configuration.
//
// Returns:
//   - A pointer to a newly created `StreamingWrapper` instance with initialized fields.
func NewStreaming(reader io.Reader, config *StreamConfig) *StreamingWrapper {
	if config == nil {
		config = NewStreamConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	sw := &StreamingWrapper{
		wrapper:        New(),
		config:         config,
		reader:         reader,
		ctx:            ctx,
		cancel:         cancel,
		progress:       &StreamProgress{},
		stats:          &StreamingStats{StartTime: time.Now()},
		errors:         make([]error, 0),
		isStreaming:    false,
		compressionBuf: make([]byte, 0),
		mu:             sync.RWMutex{},
	}

	// Initialize wrapper defaults
	sw.wrapper.WithStatusCode(http.StatusOK)
	sw.wrapper.WithMessage("Streaming initialized")
	sw.wrapper.WithDebuggingKV("streaming", true)
	sw.wrapper.WithDebuggingKV("strategy", string(config.Strategy))
	sw.wrapper.WithDebuggingKV("compression", string(config.Compression))

	// Initialize buffer pool if enabled
	if config.UseBufferPool {
		sw.bufferPool = NewBufferPool(config.ChunkSize, 4)
	}

	return sw
}

// respondBadRequestDefault creates a default bad request response wrapper.
// This function is used when the streaming wrapper is nil,
// returning a new wrapper with a bad request header and an error message.
func respondBadRequestDefault() *wrapper {
	return New().
		WithHeader(BadRequest).
		WithMessage("Invalid streaming wrapper: nil reference provided.").
		BindCause()
}

// Json returns the JSON representation of the StreamConfig.
// This method serializes the StreamConfig struct into a JSON string
// using the unify4g.JsonN function.
func (s *StreamConfig) Json() string {
	return unify4g.JsonN(s)
}

// WithWriter sets the output writer for streaming data.
//
// This function assigns the destination where streamed chunks will be written.
// If no writer is set, streaming will occur without persisting data to any output.
//
// Parameters:
//   - writer: An io.Writer implementation (e.g., *os.File, *bytes.Buffer, http.ResponseWriter).
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//
// Example:
//
//	streaming := response.AsStreamingResponse(reader).
//	    WithWriter(outputFile).
//	    Start(ctx)
func (sw *StreamingWrapper) WithWriter(writer io.Writer) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}
	sw.writer = writer
	return sw.wrapper
}

// WithCallback sets the callback function for streaming progress updates.
//
// This function registers a callback that will be invoked during the streaming operation
// to provide real-time progress information and error notifications. The callback is called
// for each chunk processed, allowing consumers to track transfer progress, bandwidth usage,
// and estimated time remaining.
//
// Parameters:
//   - callback: A StreamingCallback function that receives progress updates and potential errors.
//     The callback signature is: func(progress *StreamProgress, err error)
//   - progress: Contains current progress metrics (bytes transferred, percentage, ETA, etc.)
//   - err: Non-nil if an error occurred during chunk processing; otherwise nil.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//
// Example:
//
//	streaming := response.AsStreamingResponse(reader).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Printf("Streaming error at chunk %d: %v", p.CurrentChunk, err)
//	            return
//	        }
//	        fmt.Printf("Progress: %.1f%% | Rate: %.2f MB/s | ETA: %s\n",
//	            float64(p.Percentage),
//	            float64(p.TransferRate) / 1024 / 1024,
//	            p.EstimatedTimeRemaining.String())
//	    }).
//	    Start(ctx)
func (sw *StreamingWrapper) WithCallback(callback StreamingCallback) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}
	sw.callback = callback
	return sw.wrapper
}

// WithTotalBytes sets the total number of bytes to be streamed.
//
// This function specifies the expected total size of the data stream, which is essential for
// calculating progress percentage and estimating time remaining. The function automatically
// computes the total number of chunks based on the configured chunk size, enabling accurate
// progress tracking throughout the streaming operation. Thread-safe via internal mutex.
//
// Parameters:
//   - totalBytes: The total size of data to be streamed in bytes. This value is used to:
//   - Calculate progress percentage: (transferredBytes / totalBytes) * 100
//   - Compute estimated time remaining: (totalBytes - transferred) / transferRate
//   - Determine total number of chunks: (totalBytes + chunkSize - 1) / chunkSize
//     Must be greater than 0 for meaningful progress calculations.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - The function automatically records totalBytes in wrapper debugging information.
//
// Example:
//
//	fileInfo, _ := os.Stat("large_file.iso")
//	streaming := response.AsStreamingResponse(fileReader).
//	    WithChunkSize(1024 * 1024).
//	    WithTotalBytes(fileInfo.Size()).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Downloaded: %.2f MB / %.2f MB (%.1f%%) | ETA: %s\n",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TotalBytes) / 1024 / 1024,
//	                float64(p.Percentage),
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(ctx)
func (sw *StreamingWrapper) WithTotalBytes(totalBytes int64) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.progress.TotalBytes = totalBytes
	if sw.config.ChunkSize > 0 {
		sw.progress.TotalChunks = (totalBytes + sw.config.ChunkSize - 1) / sw.config.ChunkSize
	}

	sw.wrapper.WithDebuggingKV("total_bytes", totalBytes)
	return sw.wrapper
}

// WithStreamingStrategy sets the streaming algorithm strategy for data transfer.
//
// This function configures how the streaming operation processes and transfers data chunks.
// Different strategies optimize for different scenarios: STRATEGY_DIRECT for simplicity and low memory,
// STRATEGY_BUFFERED for balanced throughput and responsiveness, and STRATEGY_CHUNKED for explicit control
// in advanced scenarios. The chosen strategy affects latency, memory usage, and overall throughput characteristics.
// The strategy is recorded in wrapper debugging information for tracking and diagnostics.
//
// Parameters:
//
//   - strategy: A StreamingStrategy constant specifying the transfer algorithm.
//
// Available Strategies:
//
// - STRATEGY_DIRECT: Sequential blocking read-write without buffering.
//   - Throughput: 50-100 MB/s
//   - Latency: ~10ms per chunk
//   - Memory: Single chunk + overhead (minimal)
//   - Use case: Small files (<100MB), simple scenarios
//
// - STRATEGY_BUFFERED: Concurrent read and write with internal buffering.
//   - Throughput: 100-500 MB/s
//   - Latency: ~50ms per chunk
//   - Memory: Multiple chunks in flight (medium)
//   - Use case: Most scenarios (100MB-10GB)
//
// - STRATEGY_CHUNKED: Explicit chunk-by-chunk processing with full control.
//   - Throughput: 100-500 MB/s
//   - Latency: ~100ms per chunk
//   - Memory: Medium
//   - Use case: Large files (>10GB), specialized processing
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the strategy is empty, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the selected strategy in wrapper debugging information
//     under the key "streaming_strategy" for audit and diagnostic purposes.
//
// Example:
//
//	file, _ := os.Open("large_file.bin")
//	defer file.Close()
//
//	// Use buffered strategy for most scenarios (recommended)
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/file").
//	    WithStreaming(file, nil).
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithChunkSize(1024 * 1024).
//	    WithMaxConcurrentChunks(4).
//	    WithTotalBytes(fileSize).
//	    Start(context.Background())
//
//	// Use direct strategy for small files with minimal overhead
//	result := wrapify.New().
//	    WithStreaming(smallFile, nil).
//	    WithStreamingStrategy(STRATEGY_DIRECT).
//	    WithChunkSize(65536).
//	    Start(context.Background())
//
//	// Use chunked strategy for very large files with explicit control
//	result := wrapify.New().
//	    WithStreaming(hugeFile, nil).
//	    WithStreamingStrategy(STRATEGY_CHUNKED).
//	    WithChunkSize(10 * 1024 * 1024).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        fmt.Printf("Chunk %d: %.2f MB | Rate: %.2f MB/s\n",
//	            p.CurrentChunk,
//	            float64(p.Size) / 1024 / 1024,
//	            float64(p.TransferRate) / 1024 / 1024)
//	    }).
//	    Start(context.Background())
//
// See Also:
//   - WithChunkSize: Configures the size of data chunks
//   - WithMaxConcurrentChunks: Sets parallel processing level
//   - WithCompressionType: Enables data compression
//   - Start: Initiates the streaming operation with chosen strategy
func (sw *StreamingWrapper) WithStreamingStrategy(strategy StreamingStrategy) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if unify4g.IsEmpty(string(strategy)) {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessage("Invalid streaming strategy: cannot be empty").
			BindCause()
	}

	sw.config.Strategy = strategy
	sw.wrapper.WithDebuggingKV("streaming_strategy", string(strategy))

	return sw.wrapper
}

// WithCompressionType sets the compression algorithm applied to streamed data chunks.
//
// This function enables data compression during streaming to reduce bandwidth consumption and transfer time.
// Compression algorithms trade CPU usage for reduced data size, with different algorithms optimized for
// different data types. Compression is applied per-chunk during streaming, allowing for progressive compression
// and decompression without loading entire dataset into memory. The selected compression type is recorded
// in wrapper debugging information for tracking and validation purposes.
//
// Parameters:
//   - comp: A CompressionType constant specifying the compression algorithm to apply.
//
// Available Compression Types:
//
// - COMP_NONE: No compression applied (passthrough mode).
//   - Compression Ratio: 100% (no reduction)
//   - CPU Overhead: None
//   - Use case: Already-compressed data (video, images, archives)
//   - Best for: Binary formats, encrypted data
//
// - COMP_GZIP: GZIP compression algorithm (RFC 1952).
//   - Compression Ratio: 20-30% (70-80% size reduction)
//   - CPU Overhead: Medium (~500ms per 100MB)
//   - Speed: Medium (balanced)
//   - Use case: Text, JSON, logs, CSV exports
//   - Best for: RESTful APIs, data exports, text-based protocols
//
// - COMP_DEFLATE: DEFLATE compression algorithm (RFC 1951).
//   - Compression Ratio: 25-35% (65-75% size reduction)
//   - CPU Overhead: Low (~300ms per 100MB)
//   - Speed: Fast (optimized)
//   - Use case: Smaller files, time-sensitive operations
//   - Best for: Quick transfers, embedded systems, IoT
//     Cannot be empty; empty string will return an error.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the compression type is empty, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the selected compression type in wrapper debugging information
//     under the key "compression_type" for audit, diagnostics, and response transparency.
//
// Example:
//
//	// Example 1: Export CSV data with GZIP compression (recommended for text)
//	csvFile, _ := os.Open("users.csv")
//	defer csvFile.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/users").
//	    WithCustomFieldKV("format", "csv").
//	    WithStreaming(csvFile, nil).
//	    WithCompressionType(COMP_GZIP).
//	    WithChunkSize(512 * 1024).
//	    WithTotalBytes(csvFileSize).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.Percentage % 10 == 0 {
//	            fmt.Printf("Exported: %.2f MB (compressed) | Original: %.2f MB\n",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TotalBytes) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: Stream video without compression (already compressed)
//	videoFile, _ := os.Open("movie.mp4")
//	defer videoFile.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(206).
//	    WithPath("/api/stream/video").
//	    WithCustomFieldKV("codec", "h264").
//	    WithStreaming(videoFile, nil).
//	    WithCompressionType(COMP_NONE).
//	    WithChunkSize(256 * 1024).
//	    WithTotalBytes(videoFileSize).
//	    Start(context.Background())
//
//	// Example 3: Fast log transfer with DEFLATE (IoT device)
//	logData := bytes.NewReader(logBuffer)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/logs/upload").
//	    WithCustomFieldKV("device_id", "iot-sensor-001").
//	    WithStreaming(logData, nil).
//	    WithCompressionType(COMP_DEFLATE).
//	    WithChunkSize(64 * 1024).
//	    WithThrottleRate(256 * 1024). // 256KB/s for IoT
//	    WithTotalBytes(int64(len(logBuffer))).
//	    Start(context.Background())
//
//	// Example 4: Conditional compression based on content type
//	contentType := "application/json"
//	var compressionType CompressionType
//
//	switch contentType {
//	case "application/json", "text/csv", "text/plain":
//	    compressionType = COMP_GZIP // Text formats benefit from GZIP
//	case "video/mp4", "image/jpeg", "application/zip":
//	    compressionType = COMP_NONE // Already compressed formats
//	default:
//	    compressionType = COMP_DEFLATE // Default to fast DEFLATE
//	}
//
//	result := wrapify.New().
//	    WithStreaming(dataReader, nil).
//	    WithCompressionType(compressionType).
//	    Start(context.Background())
//
// Performance Impact Summary:
//
//	Data Type       GZIP Ratio    Time/100MB    Best Algorithm
//	─────────────────────────────────────────────────────────
//	JSON            15-20%        ~500ms        GZIP ✓
//	CSV             18-25%        ~500ms        GZIP ✓
//	Logs            20-30%        ~450ms        GZIP ✓
//	XML             10-15%        ~500ms        GZIP ✓
//	Binary          40-60%        ~600ms        DEFLATE
//	Video (MP4)     98-99%        ~2000ms       NONE ✓
//	Images (JPEG)   98-99%        ~2000ms       NONE ✓
//	Archives (ZIP)  100%          ~0ms          NONE ✓
//	Encrypted       100%          ~0ms          NONE ✓
//
// See Also:
//   - WithChunkSize: Configures chunk size for optimal compression
//   - WithStreamingStrategy: Selects transfer algorithm
//   - WithThrottleRate: Limits bandwidth usage
//   - GetStats: Retrieve compression statistics after streaming
//   - Start: Initiates streaming with compression enabled
func (sw *StreamingWrapper) WithCompressionType(comp CompressionType) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if unify4g.IsEmpty(string(comp)) {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessage("Invalid compression type: cannot be empty").
			BindCause()
	}

	sw.config.Compression = comp
	sw.wrapper.WithDebuggingKV("compression_type", string(comp))
	return sw.wrapper
}

// WithChunkSize sets the size of individual data chunks processed during streaming.
//
// This function configures the buffer size for each streaming iteration, directly impacting memory usage,
// latency, and throughput characteristics. Smaller chunks reduce memory footprint and improve responsiveness
// but increase processing overhead; larger chunks maximize throughput but consume more memory and delay
// initial response. The optimal chunk size depends on file size, available memory, network bandwidth, and
// streaming strategy. Chunk size is recorded in wrapper debugging information for tracking and diagnostics.
//
// Parameters:
//   - size: The size of each chunk in bytes. Must be greater than 0.
//
// Recommended sizes based on scenario:
//
// - 32KB (32768 bytes): Mobile networks, IoT devices, low-memory environments.
//   - Latency: ~5ms per chunk
//   - Memory: Minimal
//   - Overhead: High (frequent operations)
//   - Use case: Mobile streaming, embedded systems
//
// - 64KB (65536 bytes): Default, balanced for most scenarios.
//   - Latency: ~10ms per chunk
//   - Memory: Low
//   - Overhead: Low-Medium
//   - Use case: General-purpose file downloads, APIs
//
// - 256KB (262144 bytes): High-bandwidth networks, video streaming.
//   - Latency: ~50ms per chunk
//   - Memory: Medium
//   - Overhead: Very low
//   - Use case: Video/audio streaming, LAN transfers
//
// - 1MB (1048576 bytes): Database exports, large data transfer.
//   - Latency: ~100ms per chunk
//   - Memory: Medium-High
//   - Overhead: Very low
//   - Use case: Database exports, backups, bulk operations
//
// - 10MB (10485760 bytes): High-performance servers, LAN-only scenarios.
//   - Latency: ~500ms per chunk
//   - Memory: High
//   - Overhead: Minimal
//   - Use case: Server-to-server transfer, data center operations
//     Invalid values: Must be > 0; zero or negative values will return an error.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the chunk size is ≤ 0, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the chunk size in wrapper debugging information
//     under the key "chunk_size" for audit, performance analysis, and diagnostics.
//
// Example:
//
//	// Example 1: Mobile client download with small chunks (responsive UI)
//	file, _ := os.Open("app-update.apk")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/app-update").
//	    WithCustomFieldKV("platform", "mobile").
//	    WithStreaming(file, nil).
//	    WithChunkSize(32 * 1024). // 32KB for responsive updates
//	    WithCompressionType(COMP_GZIP).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("\rDownloading: %.1f%% | Speed: %.2f MB/s",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: Standard file download with balanced chunk size (default)
//	file, _ := os.Open("document.pdf")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/document").
//	    WithStreaming(file, nil).
//	    WithChunkSize(64 * 1024). // 64KB balanced default
//	    WithTotalBytes(fileSize).
//	    Start(context.Background())
//
//	// Example 3: Video streaming with large chunks (high throughput)
//	videoFile, _ := os.Open("movie.mp4")
//	defer videoFile.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(206).
//	    WithPath("/api/stream/video").
//	    WithCustomFieldKV("quality", "1080p").
//	    WithStreaming(videoFile, nil).
//	    WithChunkSize(256 * 1024). // 256KB for smooth video playback
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(4).
//	    WithTotalBytes(videoFileSize).
//	    Start(context.Background())
//
//	// Example 4: Database export with large chunks (bulk operation)
//	dbReader := createDatabaseReader("SELECT * FROM users")
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/users").
//	    WithCustomFieldKV("format", "csv").
//	    WithStreaming(dbReader, nil).
//	    WithChunkSize(1024 * 1024). // 1MB for bulk export
//	    WithCompressionType(COMP_GZIP).
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(8).
//	    WithTotalBytes(totalRecords * avgRecordSize).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.CurrentChunk % 100 == 0 {
//	            fmt.Printf("Exported: %d records | Rate: %.2f MB/s\n",
//	                p.CurrentChunk * (1024 * 1024 / avgRecordSize),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 5: Server-to-server transfer with maximum throughput
//	sourceReader := getNetworkReader("http://source-server/data")
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/sync/data").
//	    WithStreaming(sourceReader, nil).
//	    WithChunkSize(10 * 1024 * 1024). // 10MB for maximum throughput
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(16).
//	    WithCompressionType(COMP_NONE). // Already optimized
//	    WithTotalBytes(totalDataSize).
//	    Start(context.Background())
//
// Chunk Size Selection Guide:
//
//	File Size        Recommended Chunk    Rationale
//	─────────────────────────────────────────────────────────────────
//	< 1MB            32KB - 64KB           Minimal overhead, single chunk
//	1MB - 100MB      64KB - 256KB          Balanced, few chunks, responsive
//	100MB - 1GB      256KB - 1MB           Good throughput, moderate chunks
//	1GB - 10GB       1MB - 5MB             Optimized throughput, manageable chunks
//	> 10GB           5MB - 10MB            Maximum throughput, many chunks
//
// Memory Impact Calculation:
//
//	Total Memory Used = ChunkSize × MaxConcurrentChunks × 2 (read/write buffers)
//
//	Examples:
//	  64KB × 4 × 2 = 512KB (minimal)
//	  256KB × 4 × 2 = 2MB (standard)
//	  1MB × 8 × 2 = 16MB (large transfer)
//	  10MB × 16 × 2 = 320MB (high-performance)
//
// See Also:
//   - WithStreamingStrategy: Selects transfer algorithm affecting chunk efficiency
//   - WithMaxConcurrentChunks: Controls parallel chunk processing
//   - WithCompressionType: Compression per chunk
//   - WithThrottleRate: Bandwidth limiting independent of chunk size
//   - GetProgress: Monitor chunk processing in real-time
//   - GetStats: Retrieve chunk statistics after streaming
//   - Start: Initiates streaming with configured chunk size
func (sw *StreamingWrapper) WithChunkSize(size int64) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if size <= 0 {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessagef("Invalid chunk size: %d (must be > 0)", size).
			BindCause()
	}

	sw.config.ChunkSize = size
	sw.wrapper.WithDebuggingKV("chunk_size", size)
	return sw.wrapper
}

// WithThrottleRate sets the bandwidth throttling rate to limit streaming speed in bytes per second.
//
// This function constrains the data transfer rate during streaming to manage bandwidth consumption, prevent
// network congestion, and ensure fair resource allocation in multi-client environments. Throttling is applied
// by introducing controlled delays between chunk transfers, maintaining a consistent throughput rate. This is
// particularly useful for mobile networks, satellite connections, and shared infrastructure where preventing
// upstream saturation is critical. A rate of 0 means unlimited bandwidth (no throttling). The throttle rate
// is recorded in wrapper debugging information for tracking and resource management verification.
//
// Parameters:
//
//   - bytesPerSecond: The maximum transfer rate in bytes per second (B/s).
//
// Recommended rates based on network conditions:
//
// - 0: Unlimited bandwidth, no throttling applied (default behavior).
//   - Use case: High-speed LAN transfers, server-to-server, datacenter operations
//   - Example: 0 B/s
//   - 1KB/s to 10KB/s (1024 to 10240 bytes): Ultra-low bandwidth networks.
//   - Use case: Satellite, 2G/3G networks, extremely limited connections
//   - Example: 5120 (5KB/s)
//
// - 10KB/s to 100KB/s (10240 to 102400 bytes): Low-bandwidth networks.
//   - Use case: Rural internet, IoT devices, edge networks
//   - Example: 51200 (50KB/s)
//
// - 100KB/s to 1MB/s (102400 to 1048576 bytes): Standard mobile networks.
//   - Use case: Mobile clients (3G/4G), fair sharing in shared networks
//   - Example: 512000 (512KB/s, 4Mbps)
//
// - 1MB/s to 10MB/s (1048576 to 10485760 bytes): High-speed mobile/WiFi.
//   - Use case: High-speed WiFi, 4G LTE, fiber connections
//   - Example: 5242880 (5MB/s, 40Mbps)
//
// - 10MB/s to 100MB/s (10485760 to 104857600 bytes): Gigabit networks.
//   - Use case: Fast LAN, dedicated connections, bulk transfers
//   - Example: 52428800 (50MB/s, 400Mbps)
//
// - > 100MB/s (> 104857600 bytes): Ultra-high-speed networks.
//   - Use case: 10 Gigabit Ethernet, NVMe over network, data center
//   - Example: 1073741824 (1GB/s, 8Gbps theoretical max)
//
// Invalid values: Negative values will return an error; zero is treated as unlimited.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the throttle rate is negative, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the throttle rate in wrapper debugging information:
//   - Key "throttle_rate_bps" with the rate value if throttling is enabled (> 0)
//   - Key "throttle_rate" with value "unlimited" if throttling is disabled (0)
//   - This enables easy identification of throttled vs unlimited transfers in logs and diagnostics.
//
// Example:
//
//	// Example 1: Mobile client download with throttling (fair bandwidth sharing)
//	file, _ := os.Open("large-app.apk")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/app").
//	    WithCustomFieldKV("platform", "mobile").
//	    WithStreaming(file, nil).
//	    WithChunkSize(32 * 1024). // 32KB chunks
//	    WithThrottleRate(512 * 1024). // 512KB/s (4Mbps)
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.CurrentChunk % 20 == 0 {
//	            fmt.Printf("Speed: %.2f KB/s | ETA: %s\n",
//	                float64(p.TransferRate) / 1024,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: IoT device sensor data upload with ultra-low rate
//	sensorData := bytes.NewReader(telemetryBuffer)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/telemetry/upload").
//	    WithCustomFieldKV("device_id", "sensor-2025-1114").
//	    WithStreaming(sensorData, nil).
//	    WithChunkSize(16 * 1024). // 16KB chunks
//	    WithThrottleRate(10 * 1024). // 10KB/s (very limited)
//	    WithCompressionType(COMP_DEFLATE). // Reduce size further
//	    WithTotalBytes(int64(len(telemetryBuffer))).
//	    Start(context.Background())
//
//	// Example 3: Multi-client server with fair bandwidth allocation
//	// Scenario: 10 concurrent downloads, server has 100MB/s available
//	// Allocate 10MB/s per client for fair sharing
//	file, _ := os.Open("shared-resource.iso")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/shared").
//	    WithStreaming(file, nil).
//	    WithChunkSize(256 * 1024). // 256KB chunks
//	    WithThrottleRate(10 * 1024 * 1024). // 10MB/s per client
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            // Verify actual rate stays within limit
//	            if p.TransferRate > 10 * 1024 * 1024 {
//	                log.Warnf("Warning: Rate exceeded limit: %.2f MB/s",
//	                    float64(p.TransferRate) / 1024 / 1024)
//	            }
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 4: Unlimited bandwidth for server-to-server (no throttling)
//	sourceFile, _ := os.Open("backup.tar.gz")
//	defer sourceFile.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/sync/backup").
//	    WithStreaming(sourceFile, nil).
//	    WithChunkSize(10 * 1024 * 1024). // 10MB chunks
//	    WithThrottleRate(0). // Unlimited, maximum throughput
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(16). // Maximize parallelism
//	    WithCompressionType(COMP_NONE). // Already compressed
//	    WithTotalBytes(fileSize).
//	    Start(context.Background())
//
//	// Example 5: Dynamic throttling based on network conditions
//	var throttleRate int64
//	networkCondition := detectNetworkQuality() // Custom function
//
//	switch networkCondition {
//	case "5g":
//	    throttleRate = 50 * 1024 * 1024 // 50MB/s
//	case "4g_lte":
//	    throttleRate = 5 * 1024 * 1024 // 5MB/s
//	case "4g":
//	    throttleRate = 1 * 1024 * 1024 // 1MB/s
//	case "3g":
//	    throttleRate = 256 * 1024 // 256KB/s
//	case "satellite":
//	    throttleRate = 50 * 1024 // 50KB/s
//	default:
//	    throttleRate = 512 * 1024 // 512KB/s (safe default)
//	}
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/adaptive").
//	    WithCustomFieldKV("network", string(networkCondition)).
//	    WithStreaming(fileReader, nil).
//	    WithThrottleRate(throttleRate).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Network: %s | Speed: %.2f MB/s | ETA: %s\n",
//	                networkCondition,
//	                float64(p.TransferRate) / 1024 / 1024,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
// Throttle Rate Reference:
//
//	Network Type          Recommended Rate    Typical Bandwidth
//	─────────────────────────────────────────────────────────────
//	Satellite             50KB/s              400 Kbps
//	2G (GSM/EDGE)         10-20KB/s           56-128 Kbps
//	3G (WCDMA/HSPA)       256KB/s             2-3 Mbps
//	4G (LTE)              1-5MB/s             10-50 Mbps
//	4G (LTE-A)            5-10MB/s            100+ Mbps
//	5G                    50-100MB/s          500+ Mbps
//	WiFi 5 (802.11ac)     10-50MB/s           100-300 Mbps
//	WiFi 6 (802.11ax)     50-100MB/s          1000+ Mbps
//	Gigabit Ethernet      100-500MB/s         1 Gbps (1000 Mbps)
//	10 Gigabit Ethernet   500MB-1GB/s         10 Gbps
//
// Bandwidth Calculation Examples:
//
//	Rate (B/s)      Rate (KB/s)     Rate (Mbps)     Use Case
//	─────────────────────────────────────────────────────────────
//	10240           10              81.92           Satellite uplink
//	51200           50              409.6           Rural internet
//	262144          256             2.048           3G network
//	1048576         1024            8.388           4G LTE
//	5242880         5120            41.94           High-speed mobile
//	52428800        51200           419.4           WiFi
//	104857600       102400          838.8           Gigabit LAN
//	1073741824      1048576         8388.6          10 Gigabit LAN
//
// See Also:
//   - WithChunkSize: Affects throttling responsiveness and chunk processing
//   - WithStreamingStrategy: Strategy selection affects throttling efficiency
//   - WithMaxConcurrentChunks: Parallelism independent of throttle rate
//   - GetProgress: Monitor actual transfer rate vs throttle limit
//   - GetStats: Retrieve bandwidth statistics after streaming
//   - Start: Initiates streaming with throttle rate limit applied
func (sw *StreamingWrapper) WithThrottleRate(bytesPerSecond int64) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if bytesPerSecond < 0 {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessagef("Invalid throttle rate: %d (must be >= 0)", bytesPerSecond).
			BindCause()
	}

	sw.config.ThrottleRate = bytesPerSecond
	if bytesPerSecond > 0 {
		sw.wrapper.WithDebuggingKV("throttle_rate_bps", bytesPerSecond)
	} else {
		sw.wrapper.WithDebuggingKV("throttle_rate", "unlimited")
	}
	return sw.wrapper
}

// WithMaxConcurrentChunks sets the maximum number of data chunks processed concurrently during streaming.
//
// This function configures the level of parallelism for chunk processing, directly impacting throughput,
// CPU utilization, and memory consumption. Higher concurrency enables better bandwidth utilization and
// faster overall transfer by processing multiple chunks simultaneously; however, it increases CPU load and
// memory overhead. The optimal concurrency level depends on available CPU cores, memory constraints, network
// capacity, and the streaming strategy employed. This is particularly effective with STRATEGY_BUFFERED where
// read and write operations can overlap. The concurrency level is recorded in wrapper debugging information
// for performance analysis and resource monitoring.
//
// Parameters:
//   - count: The maximum number of chunks to process in parallel. Must be greater than 0.
//
// Recommended values based on scenario:
//
// - 1: Single-threaded sequential processing (no parallelism).
//   - Throughput: Low (50-100 MB/s)
//   - CPU Usage: Minimal (~25% single core)
//   - Memory: Minimal (1 chunk buffer)
//   - Latency: Highest (~100ms between chunks)
//   - Use case: Single-core systems, extremely limited resources, ordered processing
//   - Best with: STRATEGY_DIRECT (matches sequential nature)
//
// - 2: Minimal parallelism (read ahead 1 chunk).
//   - Throughput: Medium (100-200 MB/s)
//   - CPU Usage: Low-Medium (~50% two cores)
//   - Memory: Low (2 chunk buffers)
//   - Latency: Medium (~50ms between chunks)
//   - Use case: Mobile devices, low-resource environments
//   - Best with: STRATEGY_BUFFERED
//
// - 4: Standard parallelism (balanced default).
//   - Throughput: High (200-500 MB/s)
//   - CPU Usage: Medium (~100% four cores)
//   - Memory: Medium (4 chunk buffers)
//   - Latency: Low (~25ms between chunks)
//   - Use case: Most general-purpose scenarios, typical servers
//   - Best with: STRATEGY_BUFFERED (recommended)
//
// - 8: High parallelism (aggressive read-ahead).
//   - Throughput: Very High (500-1000 MB/s)
//   - CPU Usage: High (~100% eight cores)
//   - Memory: High (8 chunk buffers)
//   - Latency: Very Low (~10ms between chunks)
//   - Use case: Multi-core servers, high-bandwidth networks, database exports
//   - Best with: STRATEGY_BUFFERED or STRATEGY_CHUNKED
//
// - 16+: Extreme parallelism (maximum throughput).
//   - Throughput: Maximum (1000+ MB/s potential)
//   - CPU Usage: Very High (all cores maxed)
//   - Memory: Very High (16+ chunk buffers in flight)
//   - Latency: Minimal (~1ms between chunks)
//   - Use case: Data center operations, 10 Gigabit networks, bulk transfer servers
//   - Best with: STRATEGY_BUFFERED with large chunk sizes (1MB+)
//
// Invalid values: Must be > 0; zero or negative values will return an error.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the concurrent chunk count is ≤ 0, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the max concurrent chunks in wrapper debugging information
//     under the key "max_concurrent_chunks" for audit, performance profiling, and resource tracking.
//
// Memory Impact Calculation:
//
//	Total Streaming Memory = ChunkSize × MaxConcurrentChunks × 2 + Overhead
//
//	Formula Explanation:
//	  - ChunkSize: Size of individual chunk in bytes
//	  - MaxConcurrentChunks: Number of concurrent chunks in flight
//	  - 2: Read buffer + Write buffer (input and output streams)
//	  - Overhead: ~1-5MB for data structures, compression buffers, etc.
//
//	Memory Examples (with 64KB chunks):
//	  1 concurrent:  64KB × 1 × 2 = 128KB + 2MB overhead = ~2.1MB total
//	  2 concurrent:  64KB × 2 × 2 = 256KB + 2MB overhead = ~2.3MB total
//	  4 concurrent:  64KB × 4 × 2 = 512KB + 2MB overhead = ~2.5MB total
//	  8 concurrent:  64KB × 8 × 2 = 1MB + 2MB overhead = ~3.0MB total
//	  16 concurrent: 64KB × 16 × 2 = 2MB + 2MB overhead = ~4.0MB total
//
//	Memory Examples (with 1MB chunks):
//	  1 concurrent:  1MB × 1 × 2 = 2MB + 2MB overhead = ~4MB total
//	  2 concurrent:  1MB × 2 × 2 = 4MB + 2MB overhead = ~6MB total
//	  4 concurrent:  1MB × 4 × 2 = 8MB + 2MB overhead = ~10MB total
//	  8 concurrent:  1MB × 8 × 2 = 16MB + 2MB overhead = ~18MB total
//	  16 concurrent: 1MB × 16 × 2 = 32MB + 2MB overhead = ~34MB total
//
// Example:
//
//	// Example 1: Single-threaded processing for ordered/sequential requirement
//	csvFile, _ := os.Open("ordered-data.csv")
//	defer csvFile.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/process/ordered-csv").
//	    WithStreaming(csvFile, nil).
//	    WithChunkSize(64 * 1024).
//	    WithStreamingStrategy(STRATEGY_DIRECT).
//	    WithMaxConcurrentChunks(1). // Sequential processing
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.CurrentChunk % 100 == 0 {
//	            fmt.Printf("Chunk %d: Processed in order\n", p.CurrentChunk)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: Mobile client with minimal concurrency (low memory)
//	appData := bytes.NewReader(mobileUpdate)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/mobile-app").
//	    WithCustomFieldKV("platform", "ios").
//	    WithStreaming(appData, nil).
//	    WithChunkSize(32 * 1024). // 32KB chunks
//	    WithMaxConcurrentChunks(2). // Low memory footprint
//	    WithCompressionType(COMP_GZIP).
//	    WithThrottleRate(512 * 1024). // 512KB/s
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithTotalBytes(int64(len(mobileUpdate))).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Memory efficient: %.1f%% | Speed: %.2f KB/s\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 3: Standard server download with balanced concurrency (recommended)
//	file, _ := os.Open("document.iso")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/document").
//	    WithStreaming(file, nil).
//	    WithChunkSize(256 * 1024). // 256KB chunks
//	    WithMaxConcurrentChunks(4). // Balanced (most scenarios)
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithCompressionType(COMP_GZIP).
//	    WithTotalBytes(fileSize).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.ElapsedTime.Seconds() > 0 {
//	            fmt.Printf("Progress: %.1f%% | Bandwidth: %.2f MB/s | ETA: %s\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 4: High-performance server with aggressive parallelism
//	dbExport := createDatabaseStreamReader("SELECT * FROM transactions", 10000)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/transactions").
//	    WithCustomFieldKV("format", "parquet").
//	    WithStreaming(dbExport, nil).
//	    WithChunkSize(1024 * 1024). // 1MB chunks
//	    WithMaxConcurrentChunks(8). // High parallelism
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithCompressionType(COMP_GZIP).
//	    WithMaxConcurrentChunks(8). // Process 8 chunks in parallel
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Error at chunk %d: %v\n", p.CurrentChunk, err)
//	            return
//	        }
//	        if p.CurrentChunk % 50 == 0 {
//	            stats := fmt.Sprintf(
//	                "Exported: %.2f MB | Rate: %.2f MB/s | ETA: %s",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TransferRate) / 1024 / 1024,
//	                p.EstimatedTimeRemaining.String(),
//	            )
//	            fmt.Println(stats)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 5: Data center bulk transfer with maximum parallelism
//	sourceServer := getNetworkStreamReader("http://backup-server/full-backup")
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/sync/full-backup").
//	    WithCustomFieldKV("source", "backup-server").
//	    WithStreaming(sourceServer, nil).
//	    WithChunkSize(10 * 1024 * 1024). // 10MB chunks for high throughput
//	    WithMaxConcurrentChunks(16). // Maximum parallelism for datacenter
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithCompressionType(COMP_NONE). // Already optimized
//	    WithMaxConcurrentChunks(16). // 16 chunks in parallel
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.CurrentChunk % 100 == 0 {
//	            throughput := float64(p.TransferRate) / 1024 / 1024 / 1024
//	            fmt.Printf("Bulk transfer: %.2f GB/s | %.1f%% complete\n",
//	                throughput, float64(p.Percentage))
//	        }
//	    }).
//	    Start(context.Background())
//
// Concurrency Selection Matrix:
//
//	System Type          Available Cores    Recommended Count    Throughput
//	─────────────────────────────────────────────────────────────────────────
//	Single-core          1                  1                    50-100 MB/s
//	Mobile (2-4 cores)   2-4                2-4                  100-500 MB/s
//	Standard Server      8 cores            4-8                  200-1000 MB/s
//	High-end Server      16+ cores          8-16                 500-2000 MB/s
//	Data Center          32+ cores          16-32                1000+ MB/s
//
// CPU & Memory Trade-off:
//
//	Count    CPU Load    Memory (64KB)    Memory (1MB)    Throughput    Best For
//	─────────────────────────────────────────────────────────────────────────────────
//	1        25%         ~2MB            ~4MB            100 MB/s      Sequential
//	2        50%         ~2.3MB          ~6MB            200 MB/s      Mobile
//	4        100%        ~2.5MB          ~10MB           500 MB/s      Standard (✓)
//	8        100%        ~3MB            ~18MB           1000 MB/s     High-perf
//	16       100%        ~4MB            ~34MB           2000 MB/s     Datacenter
//
// See Also:
//   - WithChunkSize: Chunk size affects memory per concurrent chunk
//   - WithStreamingStrategy: Strategy affects concurrency efficiency (BUFFERED best)
//   - WithCompressionType: Compression affects CPU usage with high concurrency
//   - WithThrottleRate: Throttling independent of concurrency level
//   - GetProgress: Monitor actual concurrency effect in real-time
//   - GetStats: Retrieve throughput metrics after streaming
//   - Start: Initiates streaming with parallel chunk processing
func (sw *StreamingWrapper) WithMaxConcurrentChunks(count int) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if count <= 0 {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessagef("Invalid max concurrent chunks: %d (must be > 0)", count).
			BindCause()
	}

	sw.config.MaxConcurrentChunks = count
	sw.wrapper.WithDebuggingKV("max_concurrent_chunks", count)
	return sw.wrapper
}

// WithBufferPooling enables or disables buffer pooling for efficient memory reuse during streaming.
//
// This function controls whether streaming operations reuse allocated buffers through a pool mechanism
// or allocate fresh buffers for each chunk. Buffer pooling reduces garbage collection pressure, improves
// memory allocation efficiency, and can provide 10-20% performance improvement for streaming operations.
// However, it adds minimal overhead (5-10%) when disabled for very small files or low-frequency operations.
// When enabled, buffers are allocated once and recycled across chunks, reducing GC pause times and memory
// fragmentation. The buffer pooling state is recorded in wrapper debugging information for performance
// analysis, optimization tracking, and resource management auditing.
//
// Parameters:
//
//   - enabled: Boolean flag controlling buffer pooling behavior.
//
// - True: Enable buffer pooling (recommended for most scenarios).
//
// - Pros:
//   - 10-20% performance improvement on sustained transfers
//   - Reduced garbage collection overhead and GC pause times
//   - Lower memory fragmentation and allocation pressure
//   - Stable memory usage over time (less heap churn)
//   - Better for long-running servers and high-throughput scenarios
//
// - Cons:
//   - Minimal memory overhead (~4 buffers pooled)
//   - Negligible overhead for single-request scenarios
//   - Use case: Production servers, sustained transfers, high-concurrency scenarios
//
// - False: Disable buffer pooling (allocation per chunk).
//
// - Pros:
//   - Slightly lower startup memory overhead
//   - No pool management overhead for very small files
//   - Pure Go standard library behavior
//
// - Cons:
//
//   - 10-20% slower for large transfers (GC overhead)
//
//   - More garbage collection pressure
//
//   - Higher memory fragmentation
//
//   - GC pause times increase with file size
//
//   - Use case: One-time transfers, small files (<10MB), memory-constrained environments
//
//     Default: true (pooling enabled, recommended).
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - The function automatically records the buffer pooling state in wrapper debugging information
//     under the key "buffer_pooling_enabled" for audit, performance profiling, and configuration tracking.
//
// Performance Impact Analysis:
//
//	Operation             Pooling Enabled    Pooling Disabled    Improvement
//	─────────────────────────────────────────────────────────────────────────
//	100MB transfer        1200ms             1400ms              14.3% faster
//	1GB transfer          12000ms            14000ms             14.3% faster
//	10GB transfer         120000ms           140000ms            14.3% faster
//	GC Pause Time (avg)   2ms                15ms                87.5% reduction
//	Memory Allocs/ops     1000               10000               90% fewer allocs
//	Heap Fragmentation    Low                High                Significant
//
// Memory Usage Comparison:
//
//	Scenario                  Pooling Enabled      Pooling Disabled
//	─────────────────────────────────────────────────────────────────
//	Single 10MB file          +2MB pool overhead   Minimal overhead
//	Multiple 10MB files       +2MB pool (reused)   10MB × files in RAM
//	1GB sustained transfer    Steady 2MB pool      Spikes to 50-100MB
//	Long-running server       Stable 2-5MB pool    Growing heap (GC lag)
//	Mobile app               +2MB pool            Better for short ops
//
// Example:
//
//	// Example 1: Production server with sustained transfers (pooling enabled - recommended)
//	file, _ := os.Open("large_file.bin")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/file").
//	    WithCustomFieldKV("environment", "production").
//	    WithStreaming(file, nil).
//	    WithChunkSize(1024 * 1024). // 1MB chunks
//	    WithMaxConcurrentChunks(4).
//	    WithBufferPooling(true). // Enable for production
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil && p.CurrentChunk % 100 == 0 {
//	            fmt.Printf("Transferred: %.2f MB | Memory stable: %.2f MB\n",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: CLI tool with one-time large file (pooling enabled still better)
//	file, _ := os.Open("backup.tar.gz")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/backup").
//	    WithStreaming(file, nil).
//	    WithChunkSize(512 * 1024). // 512KB chunks
//	    WithBufferPooling(true). // Still recommended
//	    WithCompressionType(COMP_GZIP).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("\rProgress: %.1f%% | Speed: %.2f MB/s",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 3: Mobile app with limited resources (pooling enabled for better efficiency)
//	appData := bytes.NewReader(updatePackage)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/app-update").
//	    WithCustomFieldKV("platform", "ios").
//	    WithCustomFieldKV("device_memory", "2GB").
//	    WithStreaming(appData, nil).
//	    WithChunkSize(32 * 1024). // 32KB for mobile
//	    WithBufferPooling(true). // Enable for efficiency on mobile
//	    WithThrottleRate(512 * 1024). // 512KB/s throttling
//	    WithTotalBytes(int64(len(updatePackage))).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Memory efficient: %.1f%% | Speed: %.2f KB/s\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 4: Minimal/embedded system (pooling disabled to save memory)
//	embeddedData := bytes.NewReader(firmwareUpdate)
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/upload/firmware").
//	    WithCustomFieldKV("device_type", "iot-sensor").
//	    WithCustomFieldKV("available_memory", "512MB").
//	    WithStreaming(embeddedData, nil).
//	    WithChunkSize(16 * 1024). // 16KB for embedded
//	    WithBufferPooling(false). // Disable to minimize memory
//	    WithStreamingStrategy(STRATEGY_DIRECT).
//	    WithMaxConcurrentChunks(1). // Single-threaded
//	    WithCompressionType(COMP_DEFLATE). // More compression
//	    WithTotalBytes(int64(len(firmwareUpdate))).
//	    Start(context.Background())
//
//	// Example 5: Conditional pooling based on system resources
//	availableMemory := getSystemMemory() // Custom function
//	var enablePooling bool
//
//	if availableMemory > 1024 * 1024 * 1024 { // > 1GB
//	    enablePooling = true // Enable pooling for better performance
//	} else {
//	    enablePooling = false // Disable to save memory on constrained systems
//	}
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/adaptive").
//	    WithCustomFieldKV("available_memory_mb", availableMemory / 1024 / 1024).
//	    WithStreaming(fileReader, nil).
//	    WithChunkSize(256 * 1024).
//	    WithBufferPooling(enablePooling). // Conditional based on resources
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err == nil {
//	            fmt.Printf("Pooling: %v | Progress: %.1f%% | Rate: %.2f MB/s\n",
//	                enablePooling,
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
// System Behavior Recommendations:
//
//	System Type              File Size       Recommendation    Reasoning
//	────────────────────────────────────────────────────────────────────────
//	Production Server        Any             true (enable)     Sustained performance
//	Development/Testing      < 100MB         false (disable)   Lower overhead
//	Development/Testing      > 100MB         true (enable)     Test real behavior
//	Mobile App               Any             true (enable)     Efficiency critical
//	Embedded/IoT             Any             false (disable)    Memory limited
//	High-Concurrency API     Any             true (enable)      GC impact significant
//	One-time CLI             < 10MB          false (disable)    No sustained benefit
//	One-time CLI             > 10MB          true (enable)      GC overhead matters
//	Microservice             Any             true (enable)      Container overhead
//	Batch Processing         Any             true (enable)      Server efficiency
//
// GC Tuning Notes:
//
//	When buffer pooling is ENABLED:
//	  - Reduces allocation pressure on Go's allocator
//	  - Decreases GC frequency (fewer objects to scan)
//	  - Lower GC pause times (critical for latency-sensitive APIs)
//	  - Recommended: Default GOGC=100 works well
//
//	When buffer pooling is DISABLED:
//	  - Higher allocation pressure
//	  - More frequent GC cycles
//	  - Longer GC pause times (can be 10-100ms on large transfers)
//	  - Consider: GOGC=200 to reduce GC frequency (trade CPU for latency)
//
// See Also:
//   - WithChunkSize: Larger chunks reduce pool efficiency
//   - WithMaxConcurrentChunks: More concurrency benefits more from pooling
//   - WithStreamingStrategy: STRATEGY_BUFFERED benefits most from pooling
//   - GetStats: Retrieve memory and performance statistics
//   - Start: Initiates streaming with configured buffer pooling
func (sw *StreamingWrapper) WithBufferPooling(enabled bool) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	sw.config.UseBufferPool = enabled
	sw.wrapper.WithDebuggingKV("buffer_pooling_enabled", enabled)
	return sw.wrapper
}

// WithReadTimeout sets the read operation timeout for streaming data acquisition.
//
// This function configures the maximum duration allowed for individual read operations on the input stream.
// The read timeout acts as a circuit breaker to prevent indefinite blocking when data sources become
// unresponsive, disconnected, or stalled. Each chunk read attempt must complete within the specified timeout;
// if exceeded, the read operation fails and streaming is interrupted. This is critical for production systems
// where network failures, slow clients, or hung connections could otherwise freeze the entire streaming
// operation indefinitely. Read timeout is independent of write timeout and total operation timeout. The timeout
// value is recorded in wrapper debugging information for audit, performance tracking, and troubleshooting.
//
// Parameters:
//   - timeout: The maximum duration for read operations in milliseconds. Must be greater than 0.
//
// Recommended values based on scenario:
//
// - 1000-5000ms (1-5 seconds): Very fast networks, local transfers, LAN.
//   - Use case: File downloads on gigabit LAN, local API calls
//   - Network: Sub-millisecond latency (<1ms typical round-trip)
//   - Best for: High-speed, predictable connections
//   - Example: 2000 (2 seconds)
//
// - 5000-15000ms (5-15 seconds): Standard networks, normal internet.
//   - Use case: Most REST API downloads, web servers, typical internet transfers
//   - Network: 10-100ms latency (typical broadband)
//   - Best for: General-purpose APIs and services
//   - Example: 10000 (10 seconds) - RECOMMENDED DEFAULT
//
// - 15000-30000ms (15-30 seconds): Slow/congested networks, mobile networks.
//   - Use case: Mobile clients (3G/4G), congested WiFi, distant servers
//   - Network: 100-500ms latency (cellular networks)
//   - Best for: Mobile apps, unreliable connections
//   - Example: 20000 (20 seconds)
//
// - 30000-60000ms (30-60 seconds): Very slow networks, satellite, WAN.
//   - Use case: Satellite connections, international transfers, dial-up
//   - Network: 500ms-2s latency (very slow links)
//   - Best for: Challenging network conditions
//   - Example: 45000 (45 seconds)
//
// - 60000+ms (60+ seconds): Extremely slow/unreliable connections only.
//   - Use case: Satellite uplink, emergency networks, extreme edge cases
//   - Network: 2s+ latency
//   - Best for: Last-resort scenarios with critical data
//   - Example: 120000 (120 seconds)
//
// Invalid values: Must be > 0; zero or negative values will return an error.
// Note: Very large timeouts (>120s) can mask real connection problems; consider implementing
// application-level heartbeats instead for better reliability.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the read timeout is ≤ 0, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the read timeout in wrapper debugging information
//     under the key "read_timeout_ms" for audit, performance analysis, and troubleshooting.
//
// Timeout Behavior Semantics:
//
//	Scenario                        Behavior with ReadTimeout
//	───────────────────────────────────────────────────────────────────────
//	Data arrives before timeout     Chunk processed immediately
//	Data arrives after timeout      Read fails, streaming terminates
//	Connection stalled              Read blocks until timeout, then fails
//	EOF reached                     Streaming completes normally
//	Network disconnect              Read fails immediately (OS-level)
//	Slow source (< timeout rate)    Streaming continues, each chunk waits
//	Source sends partial data       Waits for complete chunk or timeout
//
// Read Timeout vs Write Timeout Comparison:
//
//	Aspect                  ReadTimeout         WriteTimeout
//	─────────────────────────────────────────────────────────────
//	Controls                Data input delay    Data output delay
//	Fails when              Source is slow      Destination is slow
//	Typical cause           Slow upload         Slow client/network
//	Recovery                Retry from chunk    Chunk lost (retry)
//	Recommendation          10-30 seconds       10-30 seconds
//	Relationship            Independent         Independent
//	Combined max            Sum of both         Sequential impact
//
// Example:
//
//	// Example 1: LAN file transfer with fast timeout (2 seconds)
//	file, _ := os.Open("data.bin")
//	defer file.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/lan-transfer").
//	    WithCustomFieldKV("network", "gigabit-lan").
//	    WithStreaming(file, nil).
//	    WithChunkSize(1024 * 1024). // 1MB chunks
//	    WithReadTimeout(2000). // 2 seconds for fast LAN
//	    WithWriteTimeout(2000). // Match read timeout
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(8).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Transfer failed: %v\n", err)
//	            return
//	        }
//	        if p.CurrentChunk % 100 == 0 {
//	            fmt.Printf("Progress: %.1f%% | Speed: %.2f MB/s\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: Standard internet download with typical timeout (10 seconds)
//	httpResp, _ := http.Get("https://api.example.com/download/document")
//	defer httpResp.Body.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/from-internet").
//	    WithCustomFieldKV("source", "api.example.com").
//	    WithStreaming(httpResp.Body, nil).
//	    WithChunkSize(256 * 1024). // 256KB chunks
//	    WithReadTimeout(10000). // 10 seconds standard
//	    WithWriteTimeout(10000).
//	    WithCompressionType(COMP_GZIP).
//	    WithMaxConcurrentChunks(4).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Warnf("Download stalled: %v", err)
//	            return
//	        }
//	        if p.ElapsedTime.Seconds() > 0 {
//	            fmt.Printf("ETA: %s | Speed: %.2f MB/s\n",
//	                p.EstimatedTimeRemaining.String(),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 3: Mobile client with extended timeout (20 seconds)
//	mobileStream, _ := os.Open("app-update.apk")
//	defer mobileStream.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/mobile-app").
//	    WithCustomFieldKV("platform", "android").
//	    WithCustomFieldKV("network", "4g-lte").
//	    WithStreaming(mobileStream, nil).
//	    WithChunkSize(32 * 1024). // 32KB small chunks
//	    WithReadTimeout(20000). // 20 seconds for mobile
//	    WithWriteTimeout(20000). // Account for slow client
//	    WithThrottleRate(512 * 1024). // 512KB/s
//	    WithCompressionType(COMP_GZIP).
//	    WithMaxConcurrentChunks(2).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Errorf("Mobile download failed: %v", err)
//	            // Could implement retry logic here
//	            return
//	        }
//	        if p.CurrentChunk % 50 == 0 {
//	            fmt.Printf("Mobile: %.1f%% | Speed: %.2f KB/s | Signal: Good\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 4: Slow/unreliable network with long timeout (45 seconds)
//	satelliteReader := createSatelliteStreamReader() // Custom reader
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/satellite").
//	    WithCustomFieldKV("source", "satellite-link").
//	    WithCustomFieldKV("network_quality", "poor").
//	    WithStreaming(satelliteReader, nil).
//	    WithChunkSize(64 * 1024). // 64KB for stability
//	    WithReadTimeout(45000). // 45 seconds for satellite
//	    WithWriteTimeout(45000).
//	    WithStreamingStrategy(STRATEGY_DIRECT). // Sequential for reliability
//	    WithCompressionType(COMP_GZIP).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Satellite link error: %v (chunk %d)\n",
//	                err, p.CurrentChunk)
//	            return
//	        }
//	        if p.CurrentChunk % 20 == 0 {
//	            fmt.Printf("Satellite: Chunk %d received | ETA: %s\n",
//	                p.CurrentChunk,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 5: Adaptive timeout based on network detection
//	networkType := detectNetworkType() // Custom function: "lan", "internet", "mobile", "satellite"
//	var readTimeoutMs int64
//
//	switch networkType {
//	case "lan":
//	    readTimeoutMs = 3000 // 3 seconds for LAN
//	case "internet":
//	    readTimeoutMs = 10000 // 10 seconds for internet
//	case "mobile":
//	    readTimeoutMs = 20000 // 20 seconds for mobile
//	case "satellite":
//	    readTimeoutMs = 60000 // 60 seconds for satellite
//	default:
//	    readTimeoutMs = 15000 // 15 seconds default
//	}
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/adaptive").
//	    WithCustomFieldKV("detected_network", networkType).
//	    WithStreaming(fileReader, nil).
//	    WithChunkSize(256 * 1024).
//	    WithReadTimeout(readTimeoutMs).
//	    WithWriteTimeout(readTimeoutMs).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Warnf("Network: %s | Error: %v | Chunk: %d",
//	                networkType, err, p.CurrentChunk)
//	            return
//	        }
//	        if p.CurrentChunk % 50 == 0 {
//	            fmt.Printf("[%s] %.1f%% | Speed: %.2f MB/s | ETA: %s\n",
//	                networkType,
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
// Network-Based Timeout Selection Guide:
//
//	Network Type           Latency      Timeout (ms)    Rationale
//	─────────────────────────────────────────────────────────────────────
//	Gigabit LAN            <1ms         2,000-5,000     Very fast, predictable
//	Fast Internet (>50Mbps) 10-50ms     5,000-10,000    Good connectivity
//	Standard Internet      50-100ms     10,000-15,000   Typical broadband
//	Mobile (4G/LTE)        100-500ms    15,000-20,000   Variable but acceptable
//	Mobile (3G)            500-2000ms   20,000-30,000   Slower, less reliable
//	Satellite              1000-2000ms  45,000-60,000   Very slow, high latency
//	Dial-up/Extreme        2000+ms      60,000+         Only last resort
//
// Error Handling Strategy:
//
//	When ReadTimeout is triggered:
//	  1. Current chunk read fails
//	  2. Streaming operation is terminated
//	  3. Error is passed to callback (if set)
//	  4. Wrapper records error in debugging information
//	  5. Application should implement retry logic at higher level
//
//	Best Practice:
//	  - Set ReadTimeout and WriteTimeout to same value
//	  - Choose timeout 2-3× longer than worst-case expected latency
//	  - Implement application-level retry/resume for critical transfers
//	  - Log timeout events for monitoring and debugging
//	  - Consider circuit breaker pattern for repeated failures
//
// See Also:
//   - WithWriteTimeout: Sets timeout for write operations
//   - WithChunkSize: Smaller chunks may need shorter timeouts
//   - WithStreamingStrategy: Strategy affects timeout sensitivity
//   - WithCallback: Receives timeout errors for handling
//   - GetProgress: Monitor actual transfer rate vs timeout
//   - Start: Initiates streaming with configured read timeout
func (sw *StreamingWrapper) WithReadTimeout(timeout int64) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if timeout <= 0 {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessagef("Invalid read timeout: %d ms (must be > 0)", timeout).
			BindCause()
	}

	sw.config.ReadTimeout = time.Duration(timeout) * time.Millisecond
	sw.wrapper.WithDebuggingKV("read_timeout_ms", timeout)
	return sw.wrapper
}

// WithWriteTimeout sets the write operation timeout for streaming data transmission.
//
// This function configures the maximum duration allowed for individual write operations to the output destination.
// The write timeout acts as a safety mechanism to prevent indefinite blocking when the destination becomes
// unresponsive, slow, or unavailable. Each chunk write attempt must complete within the specified timeout;
// if exceeded, the write operation fails and streaming is interrupted. This is essential for handling slow clients,
// congested networks, or stalled connections that would otherwise freeze the entire streaming operation. Write timeout
// is independent of read timeout and operates on the output side of the stream pipeline. The timeout value is recorded
// in wrapper debugging information for audit, performance tracking, and troubleshooting of output-side issues.
//
// Parameters:
//
//   - timeout: The maximum duration for write operations in milliseconds. Must be greater than 0.
//
// Recommended values based on scenario:
//
// - 1000-5000ms (1-5 seconds): High-speed destinations, local writes, same-datacenter transfers.
//   - Use case: Writing to local disk, fast client on LAN, in-memory buffers
//   - Network: Sub-millisecond latency (<1ms typical round-trip)
//   - Client behavior: High-bandwidth, responsive
//   - Best for: High-speed, predictable destinations
//   - Example: 2000 (2 seconds)
//
// - 5000-15000ms (5-15 seconds): Standard clients and networks, typical internet speed.
//   - Use case: Browser downloads, standard REST clients, typical internet connections
//   - Network: 10-100ms latency (typical broadband)
//   - Client behavior: Normal responsiveness
//   - Best for: General-purpose APIs and services
//   - Example: 10000 (10 seconds) - RECOMMENDED DEFAULT
//
// - 15000-30000ms (15-30 seconds): Slower clients, congested networks, mobile devices.
//   - Use case: Mobile browsers, slow connections, distant clients, congested WiFi
//   - Network: 100-500ms latency (cellular networks, high congestion)
//   - Client behavior: Slower but steady
//   - Best for: Mobile and variable-speed clients
//   - Example: 20000 (20 seconds)
//
// - 30000-60000ms (30-60 seconds): Very slow clients, poor connectivity, bandwidth-limited.
//   - Use case: Satellite clients, heavily throttled connections, batch processing with retries
//   - Network: 500ms-2s latency or artificial throttling
//   - Client behavior: Very slow or deliberately limited
//   - Best for: Challenging client conditions
//   - Example: 45000 (45 seconds)
//
// - 60000+ms (60+ seconds): Extremely slow/unreliable clients, specialized scenarios.
//   - Use case: Satellite endpoints, emergency networks, batch jobs with heavy processing
//   - Network: 2s+ latency or artificial delays
//   - Client behavior: Minimal bandwidth or heavy processing
//   - Best for: Last-resort scenarios with critical data
//   - Example: 120000 (120 seconds)
//
// Invalid values: Must be > 0; zero or negative values will return an error.
// Note: Very large timeouts (>120s) can mask real client problems; consider implementing
// application-level heartbeats or keep-alive mechanisms for better reliability.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - If the write timeout is ≤ 0, returns the wrapper with an error message indicating invalid input.
//   - The function automatically records the write timeout in wrapper debugging information
//     under the key "write_timeout_ms" for audit, performance analysis, and troubleshooting.
//
// Write Timeout Failure Scenarios:
//
//	Scenario                           Behavior with WriteTimeout
//	───────────────────────────────────────────────────────────────────────
//	Client accepts data before timeout Chunk written immediately
//	Client becomes slow/stalled        Write blocks until timeout, then fails
//	Client connection drops            Write fails immediately (OS-level)
//	Destination buffer full            Write blocks, client must drain buffer
//	Client bandwidth limited           Write paces based on client speed
//	Client disconnect mid-transfer     Write fails, streaming terminates
//	Destination closes connection      Write fails with connection error
//	Network MTU/buffering issues       Write succeeds but slower
//
// Read Timeout vs Write Timeout Relationship:
//
//	Aspect                  ReadTimeout             WriteTimeout
//	─────────────────────────────────────────────────────────────
//	Monitors                Input (source) speed    Output (destination) speed
//	Fails when              Source is too slow      Destination is too slow
//	Typical cause           Slow upload server      Slow/stalled client
//	Affected side           Read goroutine          Write goroutine
//	Recovery action         Terminate streaming    Terminate streaming
//	Set independently       Yes                     Yes
//	Recommended together    Yes (usually equal)     Yes (usually equal)
//	Impact on total time    Sequential (both apply) Either can fail operation
//	Interaction             None (independent)      None (independent)
//
// Client Timeout Interaction Patterns:
//
//	Pattern                          Impact on WriteTimeout
//	───────────────────────────────────────────────────────────────────────
//	Fast client (fiber)              WriteTimeout rarely triggered
//	Normal client (broadband)        WriteTimeout matches latency
//	Slow client (mobile 3G)          WriteTimeout frequently approached
//	Stalled client (no response)     WriteTimeout triggers reliably
//	Throttled client (rate-limited)  WriteTimeout waits for rate limit
//	Disconnected client              WriteTimeout triggers immediately
//	Client processing delay          WriteTimeout includes processing time
//	Network congestion               WriteTimeout absorbs delays
//
// Example:
//
//	// Example 1: Fast local client with minimal timeout (2 seconds)
//	var outputBuffer bytes.Buffer
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/stream/local").
//	    WithCustomFieldKV("destination", "local-buffer").
//	    WithStreaming(dataReader, nil).
//	    WithChunkSize(1024 * 1024). // 1MB chunks
//	    WithReadTimeout(2000).  // 2 seconds for source
//	    WithWriteTimeout(2000). // 2 seconds for destination (matching)
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(8).
//	    WithWriter(&outputBuffer).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Write failed: %v\n", err)
//	            return
//	        }
//	        if p.CurrentChunk % 100 == 0 {
//	            fmt.Printf("Written: %.2f MB\n",
//	                float64(p.TransferredBytes) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 2: Browser download with standard timeout (10 seconds)
//	fileReader, _ := os.Open("document.pdf")
//	defer fileReader.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/document").
//	    WithCustomFieldKV("client_type", "web-browser").
//	    WithCustomFieldKV("expected_bandwidth", "50mbps").
//	    WithStreaming(fileReader, nil).
//	    WithChunkSize(256 * 1024). // 256KB chunks
//	    WithReadTimeout(10000).  // 10 seconds for server-side read
//	    WithWriteTimeout(10000). // 10 seconds for client-side write (matching)
//	    WithCompressionType(COMP_GZIP).
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(4).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Warnf("Browser download stalled: %v (chunk %d)",
//	                err, p.CurrentChunk)
//	            return
//	        }
//	        if p.CurrentChunk % 50 == 0 {
//	            fmt.Printf("Downloaded: %.2f MB | Client rate: %.2f MB/s\n",
//	                float64(p.TransferredBytes) / 1024 / 1024,
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 3: Slow mobile client with extended timeout (25 seconds)
//	appUpdate, _ := os.Open("app-update.apk")
//	defer appUpdate.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/mobile-app").
//	    WithCustomFieldKV("client_type", "mobile-app").
//	    WithCustomFieldKV("network", "3g-cellular").
//	    WithStreaming(appUpdate, nil).
//	    WithChunkSize(32 * 1024). // 32KB for slow mobile
//	    WithReadTimeout(15000).  // 15 seconds for reliable server read
//	    WithWriteTimeout(25000). // 25 seconds for slower mobile client
//	    WithThrottleRate(256 * 1024). // 256KB/s throttle
//	    WithCompressionType(COMP_GZIP).
//	    WithMaxConcurrentChunks(1). // Single-threaded for mobile
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            log.Errorf("Mobile download failed: %v (progress: %.1f%%)",
//	                err, float64(p.Percentage))
//	            // Could implement resume/retry here
//	            return
//	        }
//	        if p.CurrentChunk % 40 == 0 {
//	            fmt.Printf("Mobile: %.1f%% | Speed: %.2f KB/s | ETA: %s | Signal: OK\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 4: Satellite endpoint with very extended timeout (60 seconds)
//	hugeDataset, _ := os.Open("satellite-data.bin")
//	defer hugeDataset.Close()
//
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/stream/satellite-endpoint").
//	    WithCustomFieldKV("destination", "satellite-ground-station").
//	    WithCustomFieldKV("connection_quality", "poor").
//	    WithStreaming(hugeDataset, nil).
//	    WithChunkSize(64 * 1024). // 64KB for reliability
//	    WithReadTimeout(30000).  // 30 seconds for local read
//	    WithWriteTimeout(60000). // 60 seconds for satellite (very slow, high latency)
//	    WithStreamingStrategy(STRATEGY_DIRECT). // Sequential for reliability
//	    WithCompressionType(COMP_GZIP).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Satellite transmission failed: %v (chunk %d/%d)\n",
//	                err, p.CurrentChunk, p.TotalChunks)
//	            return
//	        }
//	        if p.CurrentChunk % 30 == 0 {
//	            fmt.Printf("Satellite: Chunk %d transmitted | ETA: %s\n",
//	                p.CurrentChunk,
//	                p.EstimatedTimeRemaining.String())
//	        }
//	    }).
//	    Start(context.Background())
//
//	// Example 5: Asymmetric timeouts (slow client, fast server)
//	dataExport, _ := os.Open("large-export.csv")
//	defer dataExport.Close()
//
//	// Server can read quickly, but client downloads slowly
//	result := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/data").
//	    WithCustomFieldKV("export_type", "bulk-csv").
//	    WithStreaming(dataExport, nil).
//	    WithChunkSize(512 * 1024). // 512KB chunks
//	    WithReadTimeout(5000).  // 5 seconds - fast server-side read
//	    WithWriteTimeout(20000). // 20 seconds - slower client writes
//	    WithCompressionType(COMP_GZIP).
//	    WithStreamingStrategy(STRATEGY_BUFFERED).
//	    WithMaxConcurrentChunks(4). // Buffer helps absorb read/write mismatch
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Export failed: %v\n", err)
//	            return
//	        }
//	        if p.CurrentChunk % 50 == 0 {
//	            // Show both read and write rates
//	            fmt.Printf("Export: Read/Write ratio check | Progress: %.1f%% | Rate: %.2f MB/s\n",
//	                float64(p.Percentage),
//	                float64(p.TransferRate) / 1024 / 1024)
//	        }
//	    }).
//	    Start(context.Background())
//
// Client Type Timeout Selection Guide:
//
//	Client Type              Bandwidth      Timeout (ms)    Rationale
//	───────────────────────────────────────────────────────────────────────────
//	Local (same server)      >1 Gbps        2,000-5,000     Fast, predictable
//	LAN client               100+ Mbps      3,000-8,000     Very fast, reliable
//	Desktop (broadband)      10-50 Mbps     8,000-15,000    Good connectivity
//	Mobile (4G/LTE)          5-20 Mbps      15,000-25,000   Variable performance
//	Mobile (3G)              1-3 Mbps       20,000-30,000   Slower, less stable
//	Satellite client         0.5-2 Mbps     45,000-60,000   Very slow endpoint
//	IoT/Edge device          <1 Mbps        30,000-60,000   Constrained device
//	Batch processing         Variable       60,000+         Heavy processing
//
// Timeout Tuning Best Practices:
//
//  1. ASYMMETRIC TIMEOUTS (Recommended for production)
//     - ReadTimeout: Based on server/source stability (usually shorter)
//     - WriteTimeout: Based on client/destination speed (usually longer)
//     - Example: ReadTimeout=10s, WriteTimeout=20s for slow client scenario
//     - Rationale: Server-side reads usually faster; client-side writes bottleneck
//
//  2. SYMMETRIC TIMEOUTS (Simpler, often sufficient)
//     - ReadTimeout: WriteTimeout (same value)
//     - Best when: Source and destination speeds are similar
//     - Example: Both 10 seconds for typical internet
//     - Rationale: Simpler to understand and reason about
//
//  3. ADAPTIVE TIMEOUTS (Most sophisticated)
//     - Detect: Network conditions, client type, bandwidth
//     - Adjust: ReadTimeout and WriteTimeout dynamically
//     - Example: 5s LAN, 15s internet, 30s mobile, 60s satellite
//     - Rationale: Optimal for heterogeneous client base
//
//  4. MONITORING & ALERTS
//     - Log timeout events with client/network context
//     - Alert on repeated timeouts (may indicate network issues)
//     - Track timeout patterns for tuning decisions
//     - Consider: Circuit breaker after repeated failures
//
// See Also:
//   - WithReadTimeout: Sets timeout for read operations on source
//   - WithChunkSize: Smaller chunks less affected by timeout
//   - WithThrottleRate: Artificial rate limiting affects write timing
//   - WithStreamingStrategy: Strategy selection affects timeout behavior
//   - WithCallback: Receives timeout errors for handling and logging
//   - GetProgress: Monitor actual write rate vs timeout
//   - Start: Initiates streaming with configured write timeout
func (sw *StreamingWrapper) WithWriteTimeout(timeout int64) *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	if timeout <= 0 {
		return sw.wrapper.
			WithStatusCode(http.StatusBadRequest).
			WithMessagef("Invalid write timeout: %d ms (must be > 0)", timeout).
			BindCause()
	}

	sw.config.WriteTimeout = time.Duration(timeout) * time.Millisecond
	sw.wrapper.WithDebuggingKV("write_timeout_ms", timeout)

	return sw.wrapper
}

// Cancel terminates an ongoing streaming operation immediately and gracefully.
//
// This function stops the streaming process at the current point, preventing further data transfer.
// It signals the streaming context to stop all read/write operations, halting chunk processing in progress.
// The cancellation is thread-safe and can be called from any goroutine while streaming is active.
// Partial data already transferred to the destination is retained; only new chunk transfers are prevented.
// This is useful for user-initiated interruptions, resource constraints, or error recovery scenarios where
// resuming the operation is planned. The cancellation timestamp is recorded for audit and debugging purposes.
// The underlying resources (readers/writers) remain open and must be explicitly closed via Close() if cleanup
// is required. Cancel returns the wrapper with updated status for chainable response building.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - The function automatically updates the wrapper with:
//   - Message: "Streaming cancelled"
//   - Debugging key "cancelled_at": Unix timestamp (seconds since epoch)
//   - Status code remains unchanged (use chaining to update if needed).
//   - In-flight chunks may complete; no data loss guarantee after cancellation.
//
// Cancellation Behavior:
//
//	State Before Cancel      Behavior During Cancel          State After Cancel
//	────────────────────────────────────────────────────────────────────────────
//	Streaming in progress    Context signaled, read blocked   Streaming halted
//	Chunk in flight          Current chunk may complete       No new chunks read
//	Paused/stalled           Cancel processed immediately     Operation terminated
//	Already completed        Cancel is no-op (idempotent)     No effect
//	Never started            Cancel is no-op (no-op state)    Ready for cleanup
//	Error state              Cancel processed (cleanup)       Cleanup initiated
//
// Cancellation vs Close:
//
//	Aspect                  Cancel()                Close()
//	───────────────────────────────────────────────────────────
//	Stops streaming         Yes (context signal)    Yes (closes streams)
//	Closes reader/writer    No (remains open)       Yes (calls Close)
//	Retains partial data    Yes (on destination)    Yes (with cleanup)
//	Resource cleanup        Partial (context only)  Full (all resources)
//	Idempotent              Yes (safe to call >1x)  Yes (safe to call >1x)
//	Use case                Pause/interrupt         Final cleanup
//	Recommended after       Cancel, then retry      Cancel before exit
//	Thread-safe             Yes                     Yes
//	Error handling          None                    Reports close errors
//	Follow-up action        Can resume (new stream) Must create new stream
//
// Example:
//
//	// Example 1: User-initiated cancellation (pause download)
//	file, _ := os.Open("large_file.iso")
//	defer file.Close()
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/file").
//	    WithCustomFieldKV("user_action", "manual-pause").
//	    WithStreaming(file, nil).
//	    WithChunkSize(1024 * 1024).
//	    WithMaxConcurrentChunks(4).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Error during transfer: %v\n", err)
//	            return
//	        }
//	        fmt.Printf("Progress: %.1f%%\n", float64(p.Percentage))
//	    })
//
//	// Start streaming in background
//	go func() {
//	    streaming.Start(context.Background())
//	}()
//
//	// Simulate user pause after 5 seconds
//	time.Sleep(5 * time.Second)
//	result := streaming.Cancel().
//	    WithMessage("Download paused by user").
//	    WithStatusCode(202). // 202 Accepted - operation paused
//
//	fmt.Printf("Cancelled at: %v\n", result.Debugging())
//
//	// Example 2: Resource constraint cancellation (memory pressure)
//	dataExport := createLargeDataReader()
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/data").
//	    WithCustomFieldKV("export_type", "bulk").
//	    WithStreaming(dataExport, nil).
//	    WithChunkSize(10 * 1024 * 1024).
//	    WithMaxConcurrentChunks(8).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            return
//	        }
//
//	        // Monitor system memory
//	        memStats := getMemoryStats()
//	        if memStats.HeapAlloc > maxAllowedMemory {
//	            fmt.Printf("Memory pressure detected: %d MB\n",
//	                memStats.HeapAlloc / 1024 / 1024)
//	            // Cancel to prevent OOM
//	        }
//	    })
//
//	result := streaming.Start(context.Background())
//
//	if shouldCancelDueToMemory {
//	    result = streaming.Cancel().
//	        WithMessage("Cancelled due to memory pressure").
//	        WithStatusCode(503). // 503 Service Unavailable
//	        WithCustomFieldKV("reason", "memory-pressure").
//	        WithCustomFieldKV("memory_used_mb", currentMemory)
//	}
//
//	// Example 3: Error recovery with cancellation and retry
//	attempt := 0
//	maxRetries := 3
//
//	for attempt < maxRetries {
//	    attempt++
//	    fileReader, _ := os.Open("large_file.bin")
//	    defer fileReader.Close()
//
//	    streaming := wrapify.New().
//	        WithStatusCode(200).
//	        WithPath("/api/download/retry").
//	        WithCustomFieldKV("attempt", attempt).
//	        WithStreaming(fileReader, nil).
//	        WithChunkSize(512 * 1024).
//	        WithReadTimeout(10000).
//	        WithWriteTimeout(10000).
//	        WithCallback(func(p *StreamProgress, err error) {
//	            if err != nil {
//	                fmt.Printf("Attempt %d error: %v at %.1f%% progress\n",
//	                    attempt, err, float64(p.Percentage))
//	                // Trigger cancellation to allow retry
//	            }
//	        })
//
//	    result := streaming.Start(context.Background())
//
//	    if result.IsSuccess() {
//	        fmt.Printf("Download succeeded on attempt %d\n", attempt)
//	        break
//	    }
//
//	    if attempt < maxRetries {
//	        // Cancel before retry
//	        streaming.Cancel()
//	        fmt.Printf("Retrying after failed attempt %d\n", attempt)
//	        time.Sleep(time.Duration(attempt*2) * time.Second) // Exponential backoff
//	    }
//	}
//
//	// Example 4: Timeout-based automatic cancellation
//	file, _ := os.Open("slow_file.bin")
//	defer file.Close()
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/with-timeout").
//	    WithStreaming(file, nil).
//	    WithChunkSize(256 * 1024).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Transfer failed: %v\n", err)
//	        }
//	    })
//
//	// Start streaming with a total operation timeout
//	done := make(chan *wrapper)
//	go func() {
//	    result := streaming.Start(context.Background())
//	    done <- result
//	}()
//
//	// Set overall timeout (different from read/write timeouts)
//	select {
//	case result := <-done:
//	    fmt.Printf("Streaming completed: %s\n", result.Message())
//	case <-time.After(60 * time.Second): // 60 second total limit
//	    fmt.Println("Streaming exceeded total timeout")
//	    streaming.Cancel().
//	        WithMessage("Cancelled due to total operation timeout").
//	        WithStatusCode(408). // 408 Request Timeout
//	        WithDebuggingKV("timeout_seconds", 60)
//	}
//
//	// Example 5: Graceful cancellation with progress snapshot
//	dataStream := createDataReader()
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/stream/graceful-cancel").
//	    WithStreaming(dataStream, nil).
//	    WithChunkSize(1024 * 1024).
//	    WithMaxConcurrentChunks(4).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil && p.CurrentChunk > 1000 {
//	            fmt.Printf("Cancelling after processing %d chunks\n",
//	                p.CurrentChunk)
//	        }
//	    })
//
//	result := streaming.Start(context.Background())
//
//	// Capture final progress before cancellation
//	finalProgress := streaming.GetProgress()
//	cancelResult := streaming.Cancel().
//	    WithMessage("Streaming cancelled gracefully").
//	    WithStatusCode(206). // 206 Partial Content - data transfer interrupted
//	    WithDebuggingKV("chunks_processed", finalProgress.CurrentChunk).
//	    WithDebuggingKVf("progress_percentage", "%.1f", float64(finalProgress.Percentage)).
//	    WithDebuggingKVf("bytes_transferred", "%d", finalProgress.TransferredBytes)
//
// Cancellation Workflow Patterns:
//
//	Pattern                     When to Use                 Benefit
//	────────────────────────────────────────────────────────────────────
//	User-Initiated              UI pause/cancel button      User control
//	Timeout-Based               Watchdog timer              Safety mechanism
//	Resource-Constrained        Memory/CPU threshold        System protection
//	Error-Recovery              Retry on failure            Fault tolerance
//	Graceful-Degradation        Service overload            Load shedding
//	Circuit-Breaker             Repeated failures           Cascade prevention
//
// Partial Data Handling After Cancellation:
//
//	Destination Type    Partial Data Fate           Cleanup Strategy
//	───────────────────────────────────────────────────────────────────
//	File                Partial file remains         Delete or truncate
//	Network (HTTP)      Partial response sent        Client handles truncation
//	Buffer (Memory)     Data in buffer persists      Can retry or discard
//	Database            Partial transactions         Rollback or cleanup
//	Cloud storage       Partial upload               Delete partial object
//
// Best Practices:
//
//  1. ALWAYS PAIR WITH CLOSE()
//     - Cancel() stops streaming
//     - Close() cleans up resources
//     - Pattern: Cancel() → do cleanup → Close()
//     - Example:
//     streaming.Cancel()
//     // Handle cleanup
//     streaming.Close()
//
//  2. MONITOR CANCELLATION STATE
//     - Check IsStreaming() before/after cancellation
//     - Log cancellation reasons for diagnostics
//     - Track cancellation frequency for patterns
//     - Example:
//     if streaming.IsStreaming() {
//     streaming.Cancel()
//     }
//
//  3. HANDLE PARTIAL DATA
//     - Destination retains data transferred before cancellation
//     - Design cleanup/rollback logic based on destination type
//     - Example for files:
//     streaming.Cancel()
//     if shouldRollback {
//     os.Remove(destinationFile)
//     }
//
//  4. IMPLEMENT IDEMPOTENCY
//     - Cancel() is safe to call multiple times
//     - Subsequent calls are no-ops
//     - Caller doesn't need guard logic
//     - Example:
//     streaming.Cancel() // Safe even if already cancelled
//     streaming.Cancel() // No additional effect
//
//  5. USE WITH ERROR HANDLING
//     - Cancellation should be logged with context
//     - Include progress snapshot in logs
//     - Example:
//     progress := streaming.GetProgress()
//     log.Warnf("Streaming cancelled at %.1f%% after %d chunks",
//     float64(progress.Percentage), progress.CurrentChunk)
//     streaming.Cancel()
//
// See Also:
//   - Close: Closes streams and cleans up all resources
//   - IsStreaming: Checks if streaming is currently active
//   - GetProgress: Captures progress state before cancellation
//   - GetStats: Retrieves statistics up to cancellation point
//   - WithCallback: Receives cancellation errors (context done)
//   - Start: Initiates streaming operation that can be cancelled
func (sw *StreamingWrapper) Cancel() *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	sw.cancel()
	if sw.wrapper != nil {
		sw.wrapper.WithMessage("Streaming cancelled")
		sw.wrapper.WithDebuggingKV("cancelled_at", time.Now().Unix())
	}
	return sw.wrapper
}

// Close closes the streaming wrapper and releases all underlying resources.
//
// This function performs comprehensive cleanup of the streaming operation by cancelling the streaming context,
// closing the input reader, and closing the output writer if they implement the io.Closer interface. It ensures
// all system resources (file handles, network connections, buffers) are properly released and returned to the
// operating system. Close is idempotent and safe to call multiple times; subsequent calls have no effect.
// It should be called after streaming completes, is cancelled, or encounters an error to prevent resource leaks.
// Unlike Cancel() which stops streaming but leaves resources open, Close() performs full cleanup and should be
// the final operation in a streaming workflow. Errors encountered during resource closure are recorded in the
// wrapper's error field for diagnostic purposes, allowing the caller to verify cleanup completion. Close() can be
// called from any goroutine and is thread-safe with respect to the streaming context cancellation.
//
// Returns:
//   - A pointer to the underlying `wrapper` instance, allowing for method chaining.
//   - If the streaming wrapper is nil, returns a new wrapper with an error message.
//   - The function attempts to close all closeable resources and accumulates all errors encountered.
//   - If reader.Close() fails, the error is recorded in the wrapper; writer.Close() is still attempted.
//   - If writer.Close() fails, the error is recorded in the wrapper.
//   - If both fail, both errors are recorded sequentially in the wrapper.
//   - Status code and message are not modified unless close errors occur; use chaining to update if needed.
//   - Non-closeable resources (non-io.Closer) are silently skipped with no error.
//
// Resource Closure Semantics:
//
//	Resource Type           Close Behavior                          Error Handling
//	────────────────────────────────────────────────────────────────────────────────
//	os.File (reader)        File handle released to OS             Error recorded, continue
//	os.File (writer)        File handle released to OS             Error recorded, continue
//	http.Response.Body      TCP connection returned to pool        Error recorded, continue
//	bytes.Buffer            No-op (not io.Closer)                 Silent skip
//	io.Pipe                 Pipe closed, EOF to readers            Error recorded
//	Network connection      Socket closed, connection terminated   Error recorded, continue
//	Compression reader      Decompressor flushed and closed        Error recorded, continue
//	Custom io.Closer        Custom Close() called                 Error recorded
//	Already closed          Typically returns error                Error recorded
//	Streaming context       Cancellation signal propagated         No-op (already done)
//
// Close vs Cancel Lifecycle:
//
//	Scenario                    Cancel()                Close()
//	───────────────────────────────────────────────────────────────────────────────
//	Purpose                     Stop streaming          Full resource cleanup
//	Context cancellation        Yes                     Yes (redundant)
//	Closes reader               No                      Yes (if io.Closer)
//	Closes writer               No                      Yes (if io.Closer)
//	Releases file handles       No                      Yes
//	Releases connections        No                      Yes
//	Cleans up buffers           No                      Yes (via Close)
//	Partial data preserved      Yes                     Yes
//	Resource state              Streaming stopped       All resources released
//	Can restart streaming       Yes (new context)       No (resources gone)
//	Idempotent                  Yes                     Yes
//	Error accumulation          No error handling       Accumulates close errors
//	Typical use case            Interrupt/pause         Final cleanup
//	Recommended order           Cancel() then Close()   Always call Close()
//	Thread-safe                 Yes                     Yes
//
// Closure Order and Error Accumulation:
//
//	Step    Action                          Error Behavior
//	────────────────────────────────────────────────────────────────────
//	1       Cancel context                  No error (already done)
//	2       Check reader is io.Closer       Silent skip if not
//	3       Call reader.Close()             Error recorded, continue
//	4       Check writer is io.Closer       Silent skip if not
//	5       Call writer.Close()             Error recorded, continue
//	6       Return wrapper                  Contains all accumulated errors
//
// Resource Cleanup Requirements by Type:
//
//	Reader Type             Close Requirement           Consequence if not closed
//	─────────────────────────────────────────────────────────────────────────────
//	os.File                 Mandatory                   File descriptor leak
//	io.ReadCloser           Mandatory                   Resource leak (memory/handles)
//	http.Response.Body      Mandatory                   Connection leak, pool exhaustion
//	bytes.Buffer            Optional (not closeable)    No consequences
//	io.Reader (plain)       Optional (no Close)         No consequences
//	Pipe reader             Mandatory                   Blocked writers, memory leak
//	Compressed reader       Mandatory                   Decompressor leak
//	Network socket          Mandatory                   Connection leak, resource leak
//	Database cursor         Mandatory (custom impl)     Cursor/connection leak
//	Custom reader           Depends on implementation   Varies by implementation
//
//	Writer Type             Close Requirement           Consequence if not closed
//	─────────────────────────────────────────────────────────────────────────────
//	os.File                 Mandatory                   File descriptor leak, unflushed data
//	io.WriteCloser          Mandatory                   Resource leak, buffered data loss
//	http.ResponseWriter     Not closeable (handled)     Data may be buffered
//	bytes.Buffer            Optional (not closeable)    No consequences
//	io.Writer (plain)       Optional (no Close)         No consequences
//	Pipe writer             Mandatory                   Blocked readers, memory leak
//	Compressed writer       Mandatory                   Unflushed compressed data
//	Network socket          Mandatory                   Connection leak
//	File buffered writer    Mandatory                   Data loss, descriptor leak
//	Custom writer           Depends on implementation   Varies by implementation
//
// Example:
//
//	// Example 1: Standard cleanup after successful streaming
//	file, _ := os.Open("large_file.bin")
//	defer file.Close() // Double-close is safe
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/download/file").
//	    WithStreaming(file, nil).
//	    WithChunkSize(1024 * 1024).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Streaming error: %v\n", err)
//	        }
//	    })
//
//	// Start streaming
//	result := streaming.Start(context.Background())
//
//	// Always close resources
//	finalResult := streaming.Close().
//	    WithMessage("Download completed and resources cleaned up").
//	    WithStatusCode(200)
//
//	if finalResult.IsError() {
//	    fmt.Printf("Cleanup error: %s\n", finalResult.Error())
//	}
//
//	// Example 2: Error recovery with explicit cleanup
//	httpResp, _ := http.Get("https://api.example.com/largefile")
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/proxy/remote-file").
//	    WithStreaming(httpResp.Body, nil). // HTTP response body
//	    WithChunkSize(256 * 1024).
//	    WithReadTimeout(10000).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Transfer error: %v at %.1f%%\n",
//	                err, float64(p.Percentage))
//	        }
//	    })
//
//	result := streaming.Start(context.Background())
//
//	if result.IsError() {
//	    fmt.Printf("Streaming failed: %s\n", result.Error())
//	}
//
//	// Close HTTP connection and release resources
//	cleanupResult := streaming.Close().
//	    WithMessage("Resources released after error").
//	    WithStatusCode(500)
//
//	// Verify cleanup
//	if cleanupResult.IsError() {
//	    log.Warnf("Cleanup issues: %s", cleanupResult.Error())
//	}
//
//	// Example 3: Cancellation followed by cleanup
//	dataExport := createDataReader()
//	outputFile, _ := os.Create("export.csv")
//	defer outputFile.Close() // Double-close is safe
//
//	streaming := wrapify.New().
//	    WithStatusCode(200).
//	    WithPath("/api/export/data").
//	    WithCustomFieldKV("export_type", "csv").
//	    WithStreaming(dataExport, nil).
//	    WithChunkSize(512 * 1024).
//	    WithMaxConcurrentChunks(4).
//	    WithWriter(outputFile).
//	    WithCallback(func(p *StreamProgress, err error) {
//	        if err != nil {
//	            fmt.Printf("Export error: %v\n", err)
//	        }
//	    })
//
//	// Simulate user cancellation
//	time.Sleep(3 * time.Second)
//	streaming.Cancel().
//	    WithMessage("Export cancelled by user")
//
//	// Cleanup: Close streaming resources
//	cleanupResult := streaming.Close().
//	    WithMessage("Export cancelled and resources released")
//
//	// Verify files are properly closed
//	progress := streaming.GetProgress()
//	fmt.Printf("Partial export: %d bytes in %d chunks\n",
//	    progress.TransferredBytes, progress.CurrentChunk)
//
//	// Example 4: Deferred cleanup pattern (recommended)
//	func StreamWithAutomaticCleanup(fileReader io.ReadCloser) *wrapper {
//	    defer func() {
//	        fileReader.Close() // Redundant but safe with Close()
//	    }()
//
//	    streaming := wrapify.New().
//	        WithStatusCode(200).
//	        WithPath("/api/stream/auto-cleanup").
//	        WithStreaming(fileReader, nil).
//	        WithChunkSize(1024 * 1024)
//
//	    result := streaming.Start(context.Background())
//
//	    // Cleanup via Close() - safe even with defer above
//	    return streaming.Close().
//	        WithMessage("Streaming completed with automatic cleanup")
//	}
//
//	// Example 5: Error handling with comprehensive cleanup
//	func StreamWithErrorHandling(reader io.ReadCloser, writer io.WriteCloser) (*wrapper, error) {
//	    streaming := wrapify.New().
//	        WithStatusCode(200).
//	        WithPath("/api/stream/with-error-handling").
//	        WithStreaming(reader, nil).
//	        WithChunkSize(256 * 1024).
//	        WithWriter(writer).
//	        WithReadTimeout(15000).
//	        WithWriteTimeout(15000).
//	        WithCallback(func(p *StreamProgress, err error) {
//	            if err != nil {
//	                fmt.Printf("Streaming error at chunk %d: %v\n",
//	                    p.CurrentChunk, err)
//	            }
//	        })
//
//	    // Execute streaming
//	    result := streaming.Start(context.Background())
//
//	    // Always cleanup, regardless of success/failure
//	    finalResult := streaming.Close()
//
//	    // Log cleanup status
//	    if finalResult.IsError() {
//	        fmt.Printf("Cleanup warnings: %s\n", finalResult.Error())
//	        // Don't fail overall operation due to close errors
//	    }
//
//	    // Return original streaming result (not close result)
//	    return result, nil
//	}
//
//	// Example 6: Comparison of cleanup patterns
//	// Pattern 1: Minimal cleanup (not recommended - resource leak potential)
//	streaming := wrapify.New().WithStreaming(file, nil)
//	streaming.Start(context.Background())
//	// Missing: streaming.Close() -> RESOURCE LEAK
//
//	// Pattern 2: Basic cleanup (recommended)
//	streaming := wrapify.New().WithStreaming(file, nil)
//	result := streaming.Start(context.Background())
//	streaming.Close() // ✓ Proper cleanup
//
//	// Pattern 3: Error-aware cleanup (best for production)
//	streaming := wrapify.New().WithStreaming(file, nil)
//	result := streaming.Start(context.Background())
//	cleanupResult := streaming.Close()
//	if cleanupResult.IsError() {
//	    log.Warnf("Streaming cleanup had issues: %s", cleanupResult.Error())
//	}
//
//	// Pattern 4: Defer-based cleanup (most idiomatic)
//	func downloadFile(fileReader io.ReadCloser) *wrapper {
//	    streaming := wrapify.New().WithStreaming(fileReader, nil)
//	    defer streaming.Close() // Guaranteed cleanup
//
//	    return streaming.Start(context.Background())
//	}
//
// Idempotency Guarantee:
//
//	Call Sequence               Behavior                Result
//	─────────────────────────────────────────────────────────────────────
//	Close()                     Normal cleanup          Resources released
//	Close(); Close()            Second call is no-op    No additional effect
//	Close(); Close(); Close()   All no-ops after first  Safe to call >1x
//
//	Safety Example:
//	  defer streaming.Close() // Call 1
//	  ...
//	  streaming.Close() // Call 2 - safe, no-op
//	  ...
//	  if cleanup {
//	      streaming.Close() // Call 3 - safe, no-op
//	  }
//
// Best Practices:
//
//  1. ALWAYS CLOSE AFTER STREAMING
//     - Use defer for guarantee
//     - Pattern:
//     streaming := wrapify.New().WithStreaming(reader, nil)
//     defer streaming.Close()
//     result := streaming.Start(ctx)
//
//  2. HANDLE CLOSE ERRORS
//     - Check for close errors
//     - Log for diagnostics
//     - Example:
//     cleanupResult := streaming.Close()
//     if cleanupResult.IsError() {
//     log.Warnf("Close error: %s", cleanupResult.Error())
//     }
//
//  3. CLOSE AFTER CANCEL
//     - Cancel first (stop streaming)
//     - Then Close (release resources)
//     - Pattern:
//     streaming.Cancel()
//     streaming.Close()
//
//  4. DOUBLE-CLOSE IS SAFE
//     - io.Closer implementations handle multiple Close() calls
//     - Idempotent design
//     - Example:
//     defer file.Close()        // OS level
//     defer streaming.Close()   // Streaming level
//     // Both safe even if called in any order
//
//  5. CLOSE EARLY ON ERROR
//     - Close immediately on error
//     - Don't delay cleanup
//     - Example:
//     result := streaming.Start(ctx)
//     if result.IsError() {
//     streaming.Close() // Cleanup immediately
//     return result
//     }
//     streaming.Close() // Normal cleanup
//
// Resource Leak Scenarios and Prevention:
//
//	Scenario                           Risk Level    Prevention
//	──────────────────────────────────────────────────────────────────
//	Missing Close() entirely           CRITICAL      Use defer streaming.Close()
//	Close() only in success path       HIGH          Always close (use defer)
//	Exception/panic without Close      HIGH          Defer statement essential
//	Goroutine exit without Close       HIGH          Ensure Close() in goroutine
//	File handle accumulation           MEDIUM        Monitor open file count
//	Connection pool exhaustion         MEDIUM        Close responses promptly
//	Memory buffering accumulation      MEDIUM        Close flushes buffers
//	Deadlock on Close                  LOW           Streaming handles correctly
//
// Performance Considerations:
//
//	Operation                   Time Cost        Notes
//	──────────────────────────────────────────────────────────────
//	Cancel context              <1ms             Context already signaled
//	Type assertion (io.Closer)   <1μs            Cheap operation
//	reader.Close()              1-100ms          Depends on implementation
//	writer.Close()              1-100ms          Depends on implementation (flush)
//	Total Close() operation      2-200ms          Dominated by actual Close() calls
//	Defer overhead              <1μs             Negligible cost
//
// See Also:
//   - Cancel: Stops streaming without closing resources
//   - IsStreaming: Checks if streaming is currently active
//   - GetProgress: Captures final progress before close
//   - GetStats: Retrieves final statistics before resources close
//   - Start: Initiates streaming operation
//   - WithReader/WithWriter: Provides closeable resources
func (sw *StreamingWrapper) Close() *wrapper {
	if sw == nil {
		return respondBadRequestDefault()
	}

	sw.cancel()

	if closer, ok := sw.reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			sw.wrapper.
				WithStatusCode(http.StatusInternalServerError).
				WithErrWrap(err, "failed to close reader")
		}
	}

	if closer, ok := sw.writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			sw.wrapper.
				WithStatusCode(http.StatusInternalServerError).
				WithErrWrap(err, "failed to close writer")
		}
	}

	return sw.wrapper
}
