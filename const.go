package wrapify

// Standard HTTP headers related to content negotiation and encoding.
const (
	HeaderAccept                        = "Accept"                           // Specifies the media types that are acceptable for the response.
	HeaderAcceptCharset                 = "Accept-Charset"                   // Specifies the character sets that are acceptable.
	HeaderAcceptEncoding                = "Accept-Encoding"                  // Specifies the content encodings that are acceptable.
	HeaderAcceptLanguage                = "Accept-Language"                  // Specifies the acceptable languages for the response.
	HeaderAuthorization                 = "Authorization"                    // Contains the credentials for authenticating the client with the server.
	HeaderCacheControl                  = "Cache-Control"                    // Specifies directives for caching mechanisms in both requests and responses.
	HeaderContentDisposition            = "Content-Disposition"              // Specifies if the content should be displayed inline or treated as an attachment.
	HeaderContentEncoding               = "Content-Encoding"                 // Specifies the encoding transformations that have been applied to the body of the response.
	HeaderContentLength                 = "Content-Length"                   // Specifies the size of the response body in octets.
	HeaderContentType                   = "Content-Type"                     // Specifies the media type of the resource.
	HeaderCookie                        = "Cookie"                           // Contains stored HTTP cookies sent to the server by the client.
	HeaderHost                          = "Host"                             // Specifies the domain name of the server (for virtual hosting) and the TCP port number.
	HeaderOrigin                        = "Origin"                           // Specifies the origin of the cross-origin request or preflight request.
	HeaderReferer                       = "Referer"                          // Contains the address of the previous web page from which a link to the currently requested page was followed.
	HeaderUserAgent                     = "User-Agent"                       // Contains information about the user agent (browser or client) making the request.
	HeaderIfMatch                       = "If-Match"                         // Makes the request conditional on the target resource having the same entity tag as the one provided.
	HeaderIfNoneMatch                   = "If-None-Match"                    // Makes the request conditional on the target resource not having the same entity tag as the one provided.
	HeaderETag                          = "ETag"                             // Provides the entity tag for the resource.
	HeaderLastModified                  = "Last-Modified"                    // Specifies the last modified date of the resource.
	HeaderLocation                      = "Location"                         // Specifies the URL to redirect a client to.
	HeaderPragma                        = "Pragma"                           // Specifies implementation-specific directives that might affect caching.
	HeaderRetryAfter                    = "Retry-After"                      // Specifies the time after which the client should retry the request after receiving a 503 Service Unavailable status code.
	HeaderServer                        = "Server"                           // Contains information about the software used by the origin server to handle the request.
	HeaderWWWAuthenticate               = "WWW-Authenticate"                 // Used in HTTP response headers to indicate that the client must authenticate to access the requested resource.
	HeaderDate                          = "Date"                             // Specifies the date and time at which the message was sent.
	HeaderExpires                       = "Expires"                          // Specifies the date/time after which the response is considered stale.
	HeaderAge                           = "Age"                              // Specifies the age of the response in seconds.
	HeaderConnection                    = "Connection"                       // Specifies control options for the current connection (e.g., keep-alive or close).
	HeaderContentLanguage               = "Content-Language"                 // Specifies the language of the content.
	HeaderForwarded                     = "Forwarded"                        // Contains information about intermediate proxies or gateways that have forwarded the request.
	HeaderIfModifiedSince               = "If-Modified-Since"                // Makes the request conditional on the target resource being modified since the specified date.
	HeaderUpgrade                       = "Upgrade"                          // Requests the server to switch to a different protocol.
	HeaderVia                           = "Via"                              // Provides information about intermediate protocols and recipients between the user agent and the server.
	HeaderWarning                       = "Warning"                          // Carries additional information about the status or transformation of a message.
	HeaderXForwardedFor                 = "X-Forwarded-For"                  // Contains the originating IP address of a client connecting to a web server through an HTTP proxy or load balancer.
	HeaderXForwardedHost                = "X-Forwarded-Host"                 // Contains the original host requested by the client in the Host HTTP request header.
	HeaderXForwardedProto               = "X-Forwarded-Proto"                // Specifies the protocol (HTTP or HTTPS) used by the client.
	HeaderXRequestedWith                = "X-Requested-With"                 // Identifies the type of request being made (e.g., Ajax requests).
	HeaderXFrameOptions                 = "X-Frame-Options"                  // Specifies whether the browser should be allowed to render the page in a <frame>, <iframe>, <object>, <embed>, or <applet>.
	HeaderXXSSProtection                = "X-XSS-Protection"                 // Controls browser's built-in XSS (Cross-Site Scripting) filter.
	HeaderXContentTypeOpts              = "X-Content-Type-Options"           // Prevents browsers from interpreting files as a different MIME type than what is specified.
	HeaderContentSecurity               = "Content-Security-Policy"          // Specifies security policy for web applications, helping to prevent certain types of attacks.
	HeaderStrictTransport               = "Strict-Transport-Security"        // Enforces the use of HTTPS for the website to reduce security risks.
	HeaderPublicKeyPins                 = "Public-Key-Pins"                  // Specifies public key pins to prevent man-in-the-middle attacks.
	HeaderExpectCT                      = "Expect-CT"                        // Allows websites to specify a Certificate Transparency policy.
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"      // Specifies which domains are allowed to access the resources.
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"     // Specifies which HTTP methods are allowed when accessing the resource.
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"     // Specifies which HTTP headers can be used during the actual request.
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"           // Specifies how long the results of a preflight request can be cached.
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"    // Specifies which headers can be exposed as part of the response.
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"    // Used by the browser to indicate which HTTP method will be used during the actual request.
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"   // Specifies which headers can be sent with the actual request.
	HeaderAcceptPatch                   = "Accept-Patch"                     // Specifies which patch document formats are acceptable in the response.
	HeaderDeltaBase                     = "Delta-Base"                       // Specifies the URI of the delta information.
	HeaderIfUnmodifiedSince             = "If-Unmodified-Since"              // Makes the request conditional on the resource not being modified since the specified date.
	HeaderAcceptRanges                  = "Accept-Ranges"                    // Specifies the range of the resource that the client is requesting.
	HeaderContentRange                  = "Content-Range"                    // Specifies the range of the resource being sent in the response.
	HeaderAllow                         = "Allow"                            // Specifies the allowed methods for a resource.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials" // Indicates whether the response to the request can expose credentials.
	HeaderXCSRFToken                    = "X-CSRF-Token"                     // Used to prevent Cross-Site Request Forgery (CSRF) attacks.
	HeaderXRealIP                       = "X-Real-IP"                        // Contains the real IP address of the client, often used in proxies or load balancers.
	HeaderContentSecurityPolicy         = "Content-Security-Policy"          // Specifies content security policies to prevent certain attacks.
	HeaderReferrerPolicy                = "Referrer-Policy"                  // Controls how much information about the referring page is sent.
	HeaderExpectCt                      = "Expect-CT"                        // Specifies a Certificate Transparency policy for the web server.
	HeaderStrictTransportSecurity       = "Strict-Transport-Security"        // Enforces HTTPS to reduce the chance of security breaches.
	HeaderUpgradeInsecureRequests       = "Upgrade-Insecure-Requests"        // Requests the browser to upgrade any insecure requests to secure HTTPS requests.
)

// Media Type constants define commonly used MIME types for different content types in HTTP requests and responses.
const (
	MediaTypeApplicationJSON         = "application/json"                                                          // Specifies that the content is JSON-formatted data.
	MediaTypeApplicationXML          = "application/xml"                                                           // Specifies that the content is XML-formatted data.
	MediaTypeApplicationForm         = "application/x-www-form-urlencoded"                                         // Specifies that the content is URL-encoded form data.
	MediaTypeApplicationOctetStream  = "application/octet-stream"                                                  // Specifies that the content is binary data (not interpreted by the browser).
	MediaTypeTextPlain               = "text/plain"                                                                // Specifies that the content is plain text.
	MediaTypeTextHTML                = "text/html"                                                                 // Specifies that the content is HTML-formatted data.
	MediaTypeImageJPEG               = "image/jpeg"                                                                // Specifies that the content is a JPEG image.
	MediaTypeImagePNG                = "image/png"                                                                 // Specifies that the content is a PNG image.
	MediaTypeImageGIF                = "image/gif"                                                                 // Specifies that the content is a GIF image.
	MediaTypeAudioMP3                = "audio/mpeg"                                                                // Specifies that the content is an MP3 audio file.
	MediaTypeAudioWAV                = "audio/wav"                                                                 // Specifies that the content is a WAV audio file.
	MediaTypeVideoMP4                = "video/mp4"                                                                 // Specifies that the content is an MP4 video file.
	MediaTypeVideoAVI                = "video/x-msvideo"                                                           // Specifies that the content is an AVI video file.
	MediaTypeApplicationPDF          = "application/pdf"                                                           // Specifies that the content is a PDF file.
	MediaTypeApplicationMSWord       = "application/msword"                                                        // Specifies that the content is a Microsoft Word document.
	MediaTypeApplicationMSPowerPoint = "application/vnd.ms-powerpoint"                                             // Specifies that the content is a Microsoft PowerPoint presentation.
	MediaTypeApplicationExcel        = "application/vnd.ms-excel"                                                  // Specifies that the content is a Microsoft Excel spreadsheet.
	MediaTypeApplicationZip          = "application/zip"                                                           // Specifies that the content is a ZIP archive.
	MediaTypeApplicationGzip         = "application/gzip"                                                          // Specifies that the content is a GZIP-compressed file.
	MediaTypeMultipartFormData       = "multipart/form-data"                                                       // Specifies that the content is a multipart form, typically used for file uploads.
	MediaTypeImageBMP                = "image/bmp"                                                                 // Specifies that the content is a BMP image.
	MediaTypeImageTIFF               = "image/tiff"                                                                // Specifies that the content is a TIFF image.
	MediaTypeTextCSS                 = "text/css"                                                                  // Specifies that the content is CSS (Cascading Style Sheets).
	MediaTypeTextJavaScript          = "text/javascript"                                                           // Specifies that the content is JavaScript code.
	MediaTypeApplicationJSONLD       = "application/ld+json"                                                       // Specifies that the content is a JSON-LD (JSON for Linked Data) document.
	MediaTypeApplicationRDFXML       = "application/rdf+xml"                                                       // Specifies that the content is in RDF (Resource Description Framework) XML format.
	MediaTypeApplicationGeoJSON      = "application/geo+json"                                                      // Specifies that the content is a GeoJSON (geospatial data) document.
	MediaTypeApplicationMsgpack      = "application/msgpack"                                                       // Specifies that the content is in MessagePack format (binary JSON).
	MediaTypeApplicationOgg          = "application/ogg"                                                           // Specifies that the content is an Ogg multimedia container format.
	MediaTypeApplicationGraphQL      = "application/graphql"                                                       // Specifies that the content is in GraphQL format.
	MediaTypeApplicationProtobuf     = "application/protobuf"                                                      // Specifies that the content is in Protocol Buffers format (binary serialization).
	MediaTypeImageWebP               = "image/webp"                                                                // Specifies that the content is a WebP image.
	MediaTypeFontWOFF                = "font/woff"                                                                 // Specifies that the content is a WOFF (Web Open Font Format) font.
	MediaTypeFontWOFF2               = "font/woff2"                                                                // Specifies that the content is a WOFF2 (Web Open Font Format 2) font.
	MediaTypeAudioFLAC               = "audio/flac"                                                                // Specifies that the content is a FLAC audio file (Free Lossless Audio Codec).
	MediaTypeVideoWebM               = "video/webm"                                                                // Specifies that the content is a WebM video file.
	MediaTypeApplicationDart         = "application/dart"                                                          // Specifies that the content is a Dart programming language file.
	MediaTypeApplicationXLSX         = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"         // Specifies that the content is an Excel file in XLSX format.
	MediaTypeApplicationPPTX         = "application/vnd.openxmlformats-officedocument.presentationml.presentation" // Specifies that the content is a PowerPoint file in PPTX format.
	MediaTypeApplicationGRPC         = "application/grpc"                                                          // Specifies that the content is in gRPC format (a high-performance RPC framework).
)

const (
	UnknownXC string = "unknown"
)
