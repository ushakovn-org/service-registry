package create

import (
  "context"
  "fmt"

  "github.com/ushakovn-org/service-registry/internal/pkg/models/proto_dep"
  "github.com/ushakovn-org/service-registry/internal/pkg/usecases/graph/get"
)

type UseCase func(ctx context.Context, params Params) (*Result, error)

type Params struct {
  Owner string
  Repo  string
  Path  string
}

type Result struct {
  RenderedGraph []byte
}

type useCase struct {
  depsLoader DepsLoader
  depsSaver  DepsSaver
  getUseCase get.UseCase
}

type DepsLoader interface {
  Load(ctx context.Context, params LoadParams) (*proto_dep.ProtoDep, error)
}

type LoadParams struct {
  Owner string
  Repo  string
  Path  string
}

type DepsSaver interface {
  Save(ctx context.Context, dep proto_dep.ProtoDep) error
}

type Config struct {
  DepsLoader DepsLoader
  DepsSaver  DepsSaver
  GetUseCase get.UseCase
}

func NewUseCase(config Config) UseCase {
  uc := &useCase{
    depsLoader: config.DepsLoader,
    depsSaver:  config.DepsSaver,
    getUseCase: config.GetUseCase,
  }
  return uc.handle
}

func (uc *useCase) handle(ctx context.Context, params Params) (*Result, error) {
  protoDep, err := uc.depsLoader.Load(ctx, LoadParams(params))
  if err != nil {
    return nil, fmt.Errorf("depsLoader.Load: %w", err)
  }
  if err = uc.depsSaver.Save(ctx, *protoDep); err != nil {
    return nil, fmt.Errorf("depsSaver.Save: %w", err)
  }
  result, err := uc.getUseCase(ctx)
  if err != nil {
    return nil, fmt.Errorf("graphRenderer.Render: %w", err)
  }
  return &Result{RenderedGraph: result.RenderedGraph}, nil
}
