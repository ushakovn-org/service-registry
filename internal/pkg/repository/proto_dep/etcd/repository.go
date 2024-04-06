package etcd

import (
  "context"
  "encoding/json"
  "fmt"
  "strings"
  "sync/atomic"

  "github.com/iancoleman/strcase"
  "github.com/ushakovn-org/service-registry/internal/pkg/memory/cache"
  v3 "go.etcd.io/etcd/client/v3"

  "github.com/ushakovn-org/service-registry/internal/pkg/models/proto_dep"
)

type Repository struct {
  etcd    *v3.Client
  scanned atomic.Bool
  cache   *cache.Cache[string, proto_dep.ProtoDep]
}

type Config struct {
  EtcdClient *v3.Client
}

func NewRepository(config Config) *Repository {
  return &Repository{
    etcd:  config.EtcdClient,
    cache: cache.NewCache[string, proto_dep.ProtoDep](),
  }
}

func (r *Repository) Save(ctx context.Context, dep proto_dep.ProtoDep) error {
  value, err := buildValue(dep)
  if err != nil {
    return fmt.Errorf("buildValue: %w", err)
  }
  key := buildKey(dep)

  if _, err = r.etcd.Put(ctx, key, value); err != nil {
    return fmt.Errorf("etcd.Save: %w", err)
  }
  r.cache.Put(key, dep)

  return nil
}

func (r *Repository) Scan(ctx context.Context) ([]proto_dep.ProtoDep, error) {
  if r.scanned.Load() {
    return r.cache.Values(), nil
  }
  resp, err := r.etcd.Get(ctx, keyPrefix, v3.WithPrefix())
  if err != nil {
    return nil, fmt.Errorf("etcd.Get: %w", err)
  }
  deps := make([]proto_dep.ProtoDep, 0, len(resp.Kvs))

  for _, kv := range resp.Kvs {
    key := sanitizeKey(string(kv.Key))

    dep := proto_dep.ProtoDep{}

    if err = json.Unmarshal(kv.Value, &dep); err != nil {
      return nil, fmt.Errorf("kv.Value: json.Unmarshal: %w", err)
    }
    deps = append(deps, dep)

    r.cache.Put(key, dep)
  }
  r.scanned.Store(true)

  return deps, nil
}

func buildKey(dep proto_dep.ProtoDep) string {
  key := fmt.Sprintf("%s_%s", keyPrefix, dep.Source)
  key = strcase.ToSnake(key)
  return key
}

func buildValue(dep proto_dep.ProtoDep) (string, error) {
  value, err := json.Marshal(dep)
  if err != nil {
    return "", fmt.Errorf("json.Marshal: %w", err)
  }
  return string(value), nil
}

func sanitizeKey(key string) string {
  key = strings.TrimPrefix(key, keyPrefix)
  return key
}

const keyPrefix = "proto_deps"
