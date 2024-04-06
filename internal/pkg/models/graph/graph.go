package graph

type Root struct {
  Nodes []*Node
}

type Node struct {
  Key    string
  Attrs  []*Attr
  Linked []*Node
}

type Attr struct {
  Key   string
  Value string
}
