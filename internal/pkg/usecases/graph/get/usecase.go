package get

import (
  "context"
  "fmt"

  "github.com/ushakovn-org/service-registry/internal/pkg/models/graph"
  "github.com/ushakovn-org/service-registry/internal/pkg/models/proto_dep"
)

type UseCase func(ctx context.Context) (*Result, error)

type Result struct {
  RenderedGraph []byte
}

type useCase struct {
  depsRepo      DepsRepository
  graphRenderer GraphRenderer
}

type DepsRepository interface {
  Scan(ctx context.Context) ([]proto_dep.ProtoDep, error)
}

type GraphRenderer interface {
  Render(ctx context.Context, root graph.Root) ([]byte, error)
}

type Config struct {
  DepsRepo      DepsRepository
  GraphRenderer GraphRenderer
}

func NewUseCase(config Config) UseCase {
  uc := &useCase{
    depsRepo:      config.DepsRepo,
    graphRenderer: config.GraphRenderer,
  }
  return uc.handle
}

func (uc *useCase) handle(ctx context.Context) (*Result, error) {
  deps, err := uc.depsRepo.Scan(ctx)
  if err != nil {
    return nil, fmt.Errorf("depsRepo.Scan: %w", err)
  }
  root := toGraphRoot(deps)

  render, err := uc.graphRenderer.Render(ctx, root)
  if err != nil {
    return nil, fmt.Errorf("graphRenderer.Render: %w", err)
  }
  return &Result{RenderedGraph: render}, nil
}

func toGraphRoot(deps []proto_dep.ProtoDep) graph.Root {
  nodes := make([]*graph.Node, 0, len(deps))

  for _, dep := range deps {
    linked := make([]*graph.Node, 0, len(dep.Targets))

    for _, target := range dep.Targets {
      linked = append(linked, &graph.Node{
        Attrs: []*graph.Attr{
          {Key: "import", Value: target.Import},
          {Key: "commit", Value: target.Commit},
        },
        Key: target.Source,
      })
    }
    nodes = append(nodes, &graph.Node{
      Key:    dep.Source,
      Linked: linked,
    })
  }
  return graph.Root{Nodes: nodes}
}
