package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	libFastHttp "github.com/valyala/fasthttp"
	"log"
)

func main() {
	// Define a limit rate to 4 requests per second
	rate, err := limiter.NewRateFromFormatted("4-S")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create redis client
	option, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := redis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_phi_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a fasthttp server
	middleware := fasthttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	requestHandler := func(ctx *libFastHttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			ctx.Response.Header.SetContentType("application/json")
			ctx.Response.SetBodyString(`{"message": "ok"}`)
			ctx.SetStatusCode(libFastHttp.StatusOK)
			break
		}
	}

	fmt.Println("Server is running on port 7777")
	log.Fatal(libFastHttp.ListenAndServe(":7777", middleware.Handle(requestHandler)))
}
