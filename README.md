# Golang Middleware - Rate Limiter

This project implements a rate limiter using the **Token Bucket algorithm** in Go, designed to mitigate **Denial of Service (DoS)** attacks. It helps manage burst traffic by setting a maximum token capacity and refill rate, ensuring consistent request flow.

# Features

- **Rate Limiting Support**: Tracks remaining requests with `X-RateLimit-Remaining` header.
- **Simple Integration**: Easily integrates into existing Go applications.
- **Client-based Rate Limiting**: Rate limits API calls based on client IP address.

# Installion

You can integrate this middleware into your Go project by following these steps:

1.Run the following command to fetch the package:

```bash
go get github.com/Charan010/golang-middleware@v0.0.2
```

2.Import the package in Go:

In your Go code, import the middleware package:

```go
import "github.com/Charan010/golang-middleware"
```

3. Use RateLimitMiddleware function as a wrapper.

For example:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/Charan010/golang-middleware/ratelimiter"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
}

func main() {
    // Create a new HTTP ServeMux (router)
    mux := http.NewServeMux()

    // Register your routes
    mux.HandleFunc("/", HomePage)

    // Apply the rate-limiting middleware to the mux
    rateLimitedHandler := ratelimiter.RateLimitMiddleware(10, 1, mux)

    // Start the server with rate-limiting enabled
    log.Fatal(http.ListenAndServe(":8080", rateLimitedHandler))
}

```

In this example, the rate limiter allows up to 10 requests per second and refills 1 token per second. You can customize the rate and the maximum number of tokens according to your needs.


# Contributing:

Contributions to this project are welcome! Feel free to fork the repository, create a new branch, and submit a pull request with your improvements.

Steps for Contribution:
--Fork this repository.

--Create a new branch: 
```
bash git checkout -b feature-branch 
```

--Commit your changes:
```
 bash git commit -am 'Add new feature'
 ```

--Push to your forked repository:
```
bash git push origin featurebranch
```

--Open a pull request to the main branch.

