package contracts

import "net/http"

// RequestInterceptor defines the contract for intercepting and modifying requests.
type RequestInterceptor func(*http.Request) (*http.Request, error)

// ResponseInterceptor defines the contract for intercepting and inspecting responses.
type ResponseInterceptor func(*http.Response) (*http.Response, error)

// DataTransformer defines the contract for transforming response data before unmarshaling.
type DataTransformer func([]byte) ([]byte, error)

// ProgressCallback defines the contract for tracking upload/download progress.
type ProgressCallback func(bytesTransferred, totalBytes int64)
