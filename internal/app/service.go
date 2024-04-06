package app

import (
  "net/http"

  "github.com/go-chi/chi/v5"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn-org/service-registry/internal/app/handlers/graph"
  "github.com/ushakovn/boiler/pkg/app"
)

type Config struct {
  Graph graph.Handler
}

type Service struct {
  graph graph.Handler
}

func NewService(config Config) *Service {
  return &Service{graph: config.Graph}
}

// RegisterService stub for boiler app
func (s *Service) RegisterService(*app.RegisterParams) error {
  mux := chi.NewRouter()

  mux.Get("/graph", s.graph.Get)
  mux.Post("/graph", s.graph.Create)

  go func() {
    if err := http.ListenAndServe("localhost:8080", mux); err != nil {
      log.Errorf("http server run failed: %v", err)
    }
  }()
  return nil
}
