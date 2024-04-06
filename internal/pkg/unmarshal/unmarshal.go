package unmarshal

import (
  "encoding/json"
  "fmt"

  "gopkg.in/yaml.v3"
)

func JSON[T any](buf []byte) (*T, error) {
  return unmarshal[T](buf, json.Unmarshal)
}

func YAML[T any](buf []byte) (*T, error) {
  return unmarshal[T](buf, yaml.Unmarshal)
}

func unmarshal[T any](buf []byte, call func([]byte, any) error) (*T, error) {
  out := new(T)

  if err := call(buf, out); err != nil {
    return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
  }
  return out, nil
}
