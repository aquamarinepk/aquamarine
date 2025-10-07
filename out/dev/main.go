package main

import (
  "context"
  "embed"
  "log"
  "net/http"
  "os/signal"
  "syscall"

  "github.com/go-chi/chi/v5"
  "github.com/aquamarinepk/aquamarine/pkg/lib/am"
)

//go:embed assets
var assetsFS embed.FS

func main() {
  apiPort := ":8081"
  webPort := ":8080"

  logger := am.NewLogger("info")

  ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
  defer cancel()

  apiRouter := chi.NewRouter()
  webRouter := chi.NewRouter()
  
  apiRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("api ok"))
  })
  
  webRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("web ok"))
  })

  var deps []any

  starts, stops := am.Setup(ctx, apiRouter, webRouter, deps...)

  if err := am.Start(ctx, starts, stops); err != nil {
    log.Fatal(err)
  }

  servers := []am.Server{
    {Name: "API", Addr: apiPort, Handler: apiRouter},
    {Name: "Web", Addr: webPort, Handler: webRouter},
  }

  am.StartServers(servers, logger)

  <-ctx.Done()
  
  am.Shutdown(servers, stops, logger)
}
