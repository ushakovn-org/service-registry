package graphviz

import (
  "bytes"
  "context"
  "fmt"

  "github.com/goccy/go-graphviz"
  "github.com/goccy/go-graphviz/cgraph"
  log "github.com/sirupsen/logrus"
  "github.com/ushakovn-org/service-registry/internal/pkg/models/graph"
)

type Renderer struct{}

func NewRenderer() *Renderer {
  return &Renderer{}
}

func (r *Renderer) Render(_ context.Context, root graph.Root) ([]byte, error) {
  gvz := graphviz.New()

  cg, err := gvz.Graph()
  if err != nil {
    return nil, fmt.Errorf("gvz.Graph: %w", err)
  }

  defer func() {
    if rec := recover(); err != nil {
      log.Errorf("graphviz: Renderer.Render: panic recovered: %v", rec)
    }
    if err = gvz.Close(); err != nil {
      log.Errorf("graphviz: Renderer.Render: gvz.Close: %v", err)
    }
    if err = cg.Close(); err != nil {
      log.Errorf("graphviz: Renderer.Render: cg.Close: %v", err)
    }
  }()

  if err = r.createGraph(cg, root); err != nil {
    return nil, fmt.Errorf("createGraph: %w", err)
  }
  render := new(bytes.Buffer)

  if err = gvz.Render(cg, graphviz.SVG, render); err != nil {
    return nil, fmt.Errorf("gvz.Render: %w", err)
  }
  return render.Bytes(), nil
}

func (r *Renderer) createGraph(cg *cgraph.Graph, root graph.Root) error {
  nodeByKey := make(map[string]*cgraph.Node)

  for _, node := range root.Nodes {
    if err := r.createEdge(cg, node, nodeByKey); err != nil {
      return fmt.Errorf("createEdge: %w", err)
    }
  }
  return nil
}

func (r *Renderer) createEdge(cg *cgraph.Graph, node *graph.Node, nodeByKey map[string]*cgraph.Node) error {
  var (
    start *cgraph.Node
    end   *cgraph.Node

    err error
  )
  start, err = r.createNode(cg, node, nodeByKey)
  if err != nil {
    return fmt.Errorf("createNode: %w", err)
  }
  for _, node = range node.Linked {
    if err = r.createEdge(cg, node, nodeByKey); err != nil {
      return fmt.Errorf("createEdge: %w", err)
    }
    end, err = r.createNode(cg, node, nodeByKey)
    if err != nil {
      return fmt.Errorf("createNode: %w", err)
    }
    if _, err = cg.CreateEdge("", start, end); err != nil {
      return fmt.Errorf("graph.CreateEdge: %w", err)
    }
  }
  return nil
}

func (r *Renderer) createNode(cg *cgraph.Graph, node *graph.Node, nodeByKey map[string]*cgraph.Node) (*cgraph.Node, error) {
  cn, ok := nodeByKey[node.Key]
  if !ok {
    var err error

    cn, err = cg.CreateNode(node.Key)
    if err != nil {
      return nil, fmt.Errorf("g.CreateNode: %w", err)
    }
    nodeByKey[node.Key] = cn
  }
  return cn, nil
}
