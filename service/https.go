package service

import (
  "github.com/goccy/go-json"

  "github.com/gofiber/fiber/v3"
  "github.com/gofiber/fiber/v3/middleware/cache"
  "github.com/gofiber/fiber/v3/middleware/compress"
  "github.com/gofiber/fiber/v3/middleware/etag"
  "github.com/gofiber/fiber/v3/middleware/idempotency"
  "github.com/gofiber/fiber/v3/middleware/logger"
  "github.com/gofiber/fiber/v3/middleware/recover"
  "github.com/gofiber/fiber/v3/middleware/requestid"
  "github.com/gofiber/fiber/v3/middleware/cors"
  "github.com/gofiber/fiber/v3/client"
)

/*
HTTPS wraps the fiber app and the mongo service.
*/
type HTTPS struct {
  app      *fiber.App
}

/*
NewHTTPS sets up the fiber app and the mongo service.
*/
func NewHTTPS() *HTTPS {
  return &HTTPS{
    app: fiber.New(fiber.Config{
      CaseSensitive:            true,
      StrictRouting:            true,
      EnableSplittingOnParsers: true,
      ServerHeader:             "Fiber",
      AppName:                  "Integration",
      JSONEncoder:              json.Marshal,
      JSONDecoder:              json.Unmarshal,
    }),
  }
}

/*
Up starts the fiber app, adding the middleware and routes.
*/
func (https *HTTPS) Up() error {
  https.app.Use(
    logger.New(),
    recover.New(),
    cache.New(),
    etag.New(),
    compress.New(),
    idempotency.New(),
    requestid.New(),
    cors.New(),
  )

  https.app.Get("/", func(ctx fiber.Ctx) error {
    ctx.Status(fiber.StatusOK)
    return ctx.SendString("OK")
  })

  https.app.Post("/ingress", func(ctx fiber.Ctx) (err error) {
    ctx.Response().Header.Set("Content-Type", "application/json")

    var response *client.Response

    profile := &Profile{}

    if err = json.Unmarshal(ctx.Body(), profile); err != nil {
      return err
    }

    if response, err = client.Post("http://nuner:5051/ingress", client.Config{
      Header: map[string]string{
        "Method":       "POST",
        "Accept":       "application/json",
        "Content-Type": "application/json",
      },
      Body: profile,
    }); err != nil {
      return err
    }

    ctx.Status(response.StatusCode())
    return ctx.SendString("INGRESSED")
  })

  // Preforking allows us to bind multiple Go processes to a single port.
  // This enables significant performance gains, next to the already added
  // benefits of fasthttp.
  return https.app.Listen(":5050", fiber.ListenConfig{
    EnablePrefork: true,
  })
}
