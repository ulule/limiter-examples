package main

import (
	"fmt"
	basePhi "github.com/cj1128/phi"
	"github.com/go-redis/redis"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/phi"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/valyala/fasthttp"
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

	// Create a phi server
	middleware := phi.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	router := basePhi.NewRouter()
	router.Use(middleware.Handler)
	router.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.SetBodyString(`{"message": "ok"}`)
		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	fmt.Println("Server is running on port 7777")
	log.Fatal(fasthttp.ListenAndServe(":7777", router.ServeFastHTTP))
}
