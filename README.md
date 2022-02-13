# gin-rate

Controll the rate of calls of your (sensitive) endpoints.

## Description

An in-memory Gin middleware to limit access rate by custom key and rate. This code is just an
extension and update of the original version _github.com/yangxikun/gin-limit-by-key_.

It depends on two library:

* [golang.org/x/time/rate](https://godoc.org/golang.org/x/time/rate): rate limit
* [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache): expire limiter related key

The rate limiter uses the token method - for every request, there is a token with default renewal
period. The limiter starts with all tokens (let's say 10) and once a token is used it becomes
unavailable for the defined period (e.g. 5 minutes). Thus the client can burst 10 requests in a
second but then they need to wait full 5 minutes. Or they can behave - make a request every
minute and they will never run out of tokens.

It's crucial tu use the limiter for sensitive stuff like account creation and so on. If the account
creation takes at most 6 calls then it is sensible to set a `rate.OnIP(6, 10*time.Minute)` so there
can be only one account created every 10 minutes.

## Usage

```go
package main

import (
    "github.com/katomaso/gin-rate"
    "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.Use(rate.ByIP(10, 5*time.Minute)) // will allow 10 calls from one IP within 5 minutes (e.g. one call every 30 seconds)

	r.GET("/", func(c *gin.Context) {})

	r.Run(":8888")
}
```