# gin-rate

Controll the rate of calls of your (sensitive) endpoints.

## Description

Gin middleware that limit access rate to your endpoints using a custom key (e.g. client's IP address). 
This code is just an extension and update of the original version github.com/yangxikun/gin-limit-by-key.

It depends on two libraries:

* [golang.org/x/time/rate](https://godoc.org/golang.org/x/time/rate): rate limit
* [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache): expire limiter related key

The rate limiter uses the token method - one request eats up one token and tokens respawn after defined time.
The client can burst N requests in a second but then she will need to wait the whole time period for the tokens
to respawn. Or she can behave - make a request now and then and never feel any restriction.

It's crucial tu use the limiter for sensitive stuff like account creation and so on. If the account
creation takes at most 6 calls then it is sensible to set a `rate.ByIP(6, 10*time.Minute)` so there
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