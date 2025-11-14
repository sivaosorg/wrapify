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
