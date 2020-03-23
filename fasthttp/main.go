package main

import (
	"fmt"
	"log"

	libredis "github.com/go-redis/redis/v7"
	libfasthttp "github.com/valyala/fasthttp"

	limiter "github.com/ulule/limiter/v3"
	mfasthttp "github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func main() {
	// Define a limit rate to 4 requests per minute.
	rate, err := limiter.NewRateFromFormatted("4-H")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create redis client.
	option, err := libredis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := libredis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_phi_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a fasthttp middleware.
	middleware := mfasthttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	// Create a fasthttp request handler.
	handler := func(ctx *libfasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			ctx.Response.Header.SetContentType("application/json")
			ctx.Response.SetBodyString(`{"message": "ok"}`)
			ctx.SetStatusCode(libfasthttp.StatusOK)
		default:
			ctx.Response.Header.SetContentType("application/json")
			ctx.Response.SetBodyString(`{"message": "object-not-found"}`)
			ctx.SetStatusCode(libfasthttp.StatusNotFound)
		}
	}

	// Launch fasthttp server.
	fmt.Println("Server is running on port 7777")
	log.Fatal(libfasthttp.ListenAndServe(":7777", middleware.Handle(handler)))
}
