/*
This example shows how to use the limiter middleware with the Iris web framework.
Iris web framework docs: https://docs.iris-go.com/iris
By: Saif Aljanahi
https://github.com/1saifj
*/

package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/hello", IPRateLimit(), hello)

	app.Listen(":8080")
}

func hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}

func IPRateLimit() fiber.Handler {
	// 1. Configure
	rate := limiter.Rate{
		Period: 2 * time.Second,
		Limit:  1,
	}
	store := memory.NewStore()
	ipRateLimiter := limiter.New(ipRateLimiterConfig{
		Store:        store,
		Rate:         rate,
		KeyGenerator: func(c *fiber.Ctx) string { return utils.GetIP(c.IP()) },
	})

	// 2. Return middleware handler
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		limiterCtx, err := ipRateLimiter.Get(ctx, utils.GetIP(c.IP()))
		if err != nil {
			log.Printf("IPRateLimit - ipRateLimiter.Get - err: %v, %s on %s", err, c.IP(), c.Path())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err,
			})
		}

		c.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		c.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		if limiterCtx.Reached {
			log.Printf("Too Many Requests from %s on %s", c.IP(), c.Path())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too Many Requests on " + c.Path(),
			})
		}
		return c.Next()
	}
}
