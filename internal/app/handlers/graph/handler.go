package graph

import (
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/create"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/get"
)

type Handler struct {
  createUseCase create.UseCase
  getUseCase    get.UseCase
}

type Config struct {
  CreateUseCase create.UseCase
  GetUseCase    get.UseCase
}

func NewHandler(config Config) Handler {
  return Handler{
    createUseCase: config.CreateUseCase,
    getUseCase:    config.GetUseCase,
  }
}
