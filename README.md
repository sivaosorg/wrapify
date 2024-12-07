# wrapify

**wrapify** is a Go library designed to simplify and standardize API response wrapping for RESTful services. It leverages the Decorator Pattern to dynamically add error handling, metadata, pagination, and other response features in a clean and human-readable format. With Wrapify, you can ensure consistent and extensible API responses with minimal boilerplate. Perfect for building robust, maintainable REST APIs in Go!

### Requirements

- Go version 1.23 or higher

### Installation

To install, you can use the following commands based on your preference:

- For a specific version:

  ```bash
  go get github.com/sivaosorg/wrapify@v0.0.1
  ```

- For the latest version:
  ```bash
  go get -u github.com/sivaosorg/wrapify@latest
  ```

### Getting started

#### Getting wrapify

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```go
import "github.com/sivaosorg/wrapify"
```

### Contributing

To contribute to project, follow these steps:

1. Clone the repository:

   ```bash
   git clone --depth 1 https://github.com/sivaosorg/wrapify.git
   ```

2. Navigate to the project directory:

   ```bash
   cd wrapify
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```
